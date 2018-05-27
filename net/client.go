package net

import (
	"net/http"
	"time"
)

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
