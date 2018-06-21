package middleware

import (
	"net/http"

	"gitlab.com/Cacophony/dhelpers"
)

// Recoverer middleware handles errors gracefully, sents them to sentry, logs them, and responds with an internal server error
func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err, ok := recover().(error)
			if err == nil || !ok {
				return
			}

			// log error
			dhelpers.HandleHTTPErrorWith(GetService(r.Context()), r, err, dhelpers.SentryErrorHandler)

			// send internal server error response
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
