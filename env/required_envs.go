package env

import (
	"errors"
	"github.com/getsentry/sentry-go"
	"os"
)

// RequiredEnvs means that each env in the array is required or a error is thrown.
func RequiredEnvs(Envs ...string) map[string]string {
	m := make(map[string]string, len(Envs))
	for _, v := range Envs {
		e := os.Getenv(v)
		if e == "" {
			err := errors.New(v + " does not exist")
			sentry.CaptureException(err)
			panic(err)
		}
		m[v] = e
	}
	return m
}
