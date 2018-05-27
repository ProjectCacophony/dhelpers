package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
	"sync"

	"gitlab.com/Cacophony/dhelpers"
)

// modified https://gist.github.com/erikdubbelboer/7df2b2b9f34f9f839a84

var (
	// sync pool for responses
	pool = sync.Pool{
		New: func() interface{} {
			w, _ := gzip.NewWriterLevel(nil, gzip.BestSpeed) // nolint: 	gas
			return &gzipResponseWriter{
				w: w,
			}
		},
	}
)

// gzipResponseWriter to return
type gzipResponseWriter struct {
	http.ResponseWriter

	w             *gzip.Writer
	statusCode    int
	headerWritten bool
}

// WriteHeader sets the content length and content encoding
func (gzr *gzipResponseWriter) WriteHeader(statusCode int) {
	gzr.statusCode = statusCode
	gzr.headerWritten = true

	if gzr.statusCode != http.StatusNotModified && gzr.statusCode != http.StatusNoContent {
		gzr.ResponseWriter.Header().Del("Content-Length")
		gzr.ResponseWriter.Header().Set("Content-Encoding", "gzip")
	}

	gzr.ResponseWriter.WriteHeader(statusCode)
}

// Write sets content type if not set yet, and write into the response
func (gzr *gzipResponseWriter) Write(b []byte) (int, error) {
	if _, ok := gzr.Header()["Content-Type"]; !ok {
		// If no content type, apply sniffing algorithm to un-gzipped body.
		gzr.ResponseWriter.Header().Set("Content-Type", http.DetectContentType(b))
	}

	if !gzr.headerWritten {
		// This is exactly what Go would also do if it hasn't been written yet.
		gzr.WriteHeader(http.StatusOK)
	}

	return gzr.w.Write(b)
}

// Flush flushes any pending compressed data to the writer
func (gzr *gzipResponseWriter) Flush() {
	if gzr.w != nil {
		err := gzr.w.Flush()
		dhelpers.LogError(err)
	}

	if fw, ok := gzr.ResponseWriter.(http.Flusher); ok {
		fw.Flush()
	}
}

// GzipMiddleware encodes a response as gzip if accepted by requester
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if requester doesn't accept gzip, call next handler with raw response
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// get response writer from pool
		gzr := pool.Get().(*gzipResponseWriter)
		gzr.statusCode = 0
		gzr.headerWritten = false
		gzr.ResponseWriter = w
		gzr.w.Reset(w)

		// write response async
		defer func() {
			// gzr.w.Close will write a footer even if no data has been written.
			// StatusNotModified and StatusNoContent expect an empty body so don't close it.
			if gzr.statusCode != http.StatusNotModified && gzr.statusCode != http.StatusNoContent {
				err := gzr.w.Close()
				// log any closing errors
				dhelpers.LogError(err)
			}
			pool.Put(gzr)
		}()

		// call next handler
		next.ServeHTTP(gzr, r)
	})
}
