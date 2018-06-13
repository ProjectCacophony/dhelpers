package net

import (
	"net/http"
	"time"
)

// Get does a GET request and returns the result, returns an error if the StatusCode was not 2xx
func Get(url string) ([]byte, error) {
	return GetTimeout(url, time.Second*15)
}

// GetResilient does a GET request with up to three retries and concurrency, and returns the result, returns an error if the StatusCode was not 2xx
func GetResilient(url string) ([]byte, error) {
	return GetTimeoutAndRetries(url, time.Second*15, 3)
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

	// read body or return error
	return readResponse(response)
}

// GetTimeoutAndRetries does a GET request, with up to n retries, and returns the result, with a specified timeout, returns an error if the StatusCode was not 2xx
func GetTimeoutAndRetries(url string, timeout time.Duration, retries int) ([]byte, error) {
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
	response, err := GetPesterClient(timeout, 2, retries).Do(request)
	if err != nil {
		return nil, err
	}

	// read body or return error
	return readResponse(response)
}
