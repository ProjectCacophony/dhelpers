package net

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"strconv"

	"net/http"

	"github.com/sethgrid/pester"
	"gitlab.com/Cacophony/dhelpers/cache"
)

func pesterLogHook(e pester.ErrEntry) {
	cache.GetLogger().WithField("module", "pester").Warnln(
		"failed", e.Method, e.URL, "(", e.Attempt, " attempt):", e.Err.Error(),
	)
}

func readResponse(response *http.Response) ([]byte, error) {
	// Only continue if code was 200 - 299
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, errors.New("expected status 200; got " + strconv.Itoa(response.StatusCode))
	}

	// Read body
	if response.Body != nil {
		defer func() {
			closeBodyErr := response.Body.Close()
			if closeBodyErr != nil {
				cache.GetLogger().WithError(closeBodyErr).Errorln("error closing body")
			}
		}()
	}

	// create bytes buffer
	buf := bytes.NewBuffer(nil)

	// read content-encoding
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		// decompress gzip if required
		var gzipReader *gzip.Reader
		gzipReader, err := gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
		// close reader
		if gzipReader != nil {
			defer func() {
				closeGzipErr := gzipReader.Close()
				if closeGzipErr != nil {
					cache.GetLogger().WithError(closeGzipErr).Errorln("error closing gzip reader")
				}
			}()
		}
		// copy result to buffer
		_, err = io.Copy(buf, gzipReader)
		if err != nil {
			return nil, err
		}
	default:
		// copy raw data to buffer
		_, err := io.Copy(buf, response.Body)
		if err != nil {
			return nil, err
		}
	}

	// return bytes
	return buf.Bytes(), nil
}
