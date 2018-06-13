package net

import (
	"net/http"
	"time"

	"github.com/sethgrid/pester"
)

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

// GetPesterClient returns a Pester client with the specified variables
func GetPesterClient(timeout time.Duration, concurrency int, retries int) *pester.Client {
	client := pester.NewExtendedClient(GetHTTPClientTimeout(timeout))
	client.Concurrency = concurrency
	client.MaxRetries = retries
	client.Backoff = pester.ExponentialBackoff
	client.LogHook = pesterLogHook
	return client
}
