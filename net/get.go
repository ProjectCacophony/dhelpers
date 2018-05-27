package net

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Get does a GET request and returns the result, returns an error if the StatusCode was not 2xx
func Get(url string) ([]byte, error) {
	return GetTimeout(url, time.Second*15)
}

// GetTimeout does a GET request and returns the result, with a specified timeout, returns an error if the StatusCode was not 2xx
func GetTimeout(url string, timeout time.Duration) ([]byte, error) {
	// Prepare request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// set user agent header
	request.Header.Set("User-Agent", defaultUA)
	// allow receiving gzip
	request.Header.Add("Accept-Encoding", "gzip")

	// Do request
	response, err := GetHTTPClientTimeoutWithoutKeepAlive(timeout).Do(request)
	if err != nil {
		return nil, err
	}

	// Only continue if code was 200 - 299
	if response.StatusCode < 200 || response.StatusCode > 299 {
		if err != nil {
			return nil, errors.New("expected status 200; got " + strconv.Itoa(response.StatusCode))
		}
	}

	// Read body
	if response.Body != nil {
		defer response.Body.Close() // nolint: errcheck
	}

	// create bytes buffer
	buf := bytes.NewBuffer(nil)

	// read content-encoding
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		// decompress gzip if required
		var gzipReader *gzip.Reader
		gzipReader, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
		// close reader
		if gzipReader != nil {
			defer gzipReader.Close() // nolint: errcheck
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
