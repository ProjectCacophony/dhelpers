package dhelpers

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
)

// TODO: version
var defaultUA = "ProjectD/0.1 (https://gitlab.com/project-d-collab)"

// NetGet does a GET request and returns the result, returns an error if the StatusCode was not 2xx
func NetGet(url string) ([]byte, error) {
	return NetGetTimeout(url, time.Second*15)
}

// NetGetTimeout does a GET request and returns the result, with a specified timeout, returns an error if the StatusCode was not 2xx
func NetGetTimeout(url string, timeout time.Duration) ([]byte, error) {
	// Allocate client
	client := &http.Client{
		Timeout: timeout,
	}

	// Prepare request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}
	request.Header.Set("User-Agent", defaultUA)

	// Do request
	response, err := client.Do(request)
	if err != nil {
		return []byte{}, err
	}

	// Only continue if code was 200
	if response.StatusCode/100 != 2 {
		if err != nil {
			return []byte{}, errors.New("expected status 200; got " + strconv.Itoa(response.StatusCode))
		}
	} else {
		// Read body
		defer response.Body.Close() // nolint: errcheck

		buf := bytes.NewBuffer(nil)
		_, err := io.Copy(buf, response.Body)
		if err != nil {
			return []byte{}, err
		}

		return buf.Bytes(), nil
	}
	return []byte{}, errors.New("internal error")
}
