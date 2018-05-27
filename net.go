package dhelpers

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TODO: version
var defaultUA = "Cacophony/0.1 (https://gitlab.com/Cacophony)"

// GetHTTPClient returns a HTTP client with 15 seconds timeout
func GetHTTPClient() *http.Client {
	return GetHTTPClientTimeout(time.Second * 15)
}

// GetHTTPClientTimeout returns a HTTP client with a specified timeout
func GetHTTPClientTimeout(timeout time.Duration) *http.Client {
	// create http client with given timeout
	return &http.Client{
		Timeout: timeout,
	}
}

// GetHTTPClientTimeoutWithoutKeepAlive returns a HTTP client with a specified timeout and without http keep alive
func GetHTTPClientTimeoutWithoutKeepAlive(timeout time.Duration) *http.Client {
	// create http client with given timeout
	return &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
		Timeout: timeout,
	}
}

// CleanURL makes a URL posted in discord ready to use for further usage
func CleanURL(uncleanedURL string) (url string) {
	if strings.HasPrefix(uncleanedURL, "<") {
		uncleanedURL = strings.TrimLeft(uncleanedURL, "<")
	}
	if strings.HasSuffix(uncleanedURL, ">") {
		uncleanedURL = strings.TrimRight(uncleanedURL, ">")
	}
	return uncleanedURL
}

// NetGet does a GET request and returns the result, returns an error if the StatusCode was not 2xx
func NetGet(url string) ([]byte, error) {
	return NetGetTimeout(url, time.Second*15)
}

// NetGetTimeout does a GET request and returns the result, with a specified timeout, returns an error if the StatusCode was not 2xx
func NetGetTimeout(url string, timeout time.Duration) ([]byte, error) {
	// Prepare request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}
	request.Header.Set("User-Agent", defaultUA)

	// Do request
	response, err := GetHTTPClientTimeoutWithoutKeepAlive(timeout).Do(request)
	if err != nil {
		return []byte{}, err
	}

	// Only continue if code was 200 - 299
	if response.StatusCode < 200 || response.StatusCode > 299 {
		if err != nil {
			return []byte{}, errors.New("expected status 200; got " + strconv.Itoa(response.StatusCode))
		}
	} else {
		// Read body
		if response.Body != nil {
			defer response.Body.Close() // nolint: errcheck
		}

		buf := bytes.NewBuffer(nil)
		_, err := io.Copy(buf, response.Body)
		if err != nil {
			return []byte{}, err
		}

		return buf.Bytes(), nil
	}
	return []byte{}, errors.New("internal error")
}
