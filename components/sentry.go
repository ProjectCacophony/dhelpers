package components

import (
	"os"

	"github.com/getsentry/raven-go"
)

// InitSentry sets up the raven client
// it reads the DSN from the environment variable RAVEN_DSN
// it sets the release to the environment variable VERSION if set
func InitSentry() (err error) {
	if os.Getenv("RAVEN_DSN") != "" {
		err = raven.SetDSN(os.Getenv("RAVEN_DSN"))
		if os.Getenv("VERSION") != "" {
			raven.SetRelease(os.Getenv("VERSION"))
		}
	}
	return err
}
