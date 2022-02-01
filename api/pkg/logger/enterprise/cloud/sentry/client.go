package sentry

import (
	"github.com/getsentry/raven-go"
)

func NewClient() (*raven.Client, error) {
	return raven.New("https://7ca135fcbbbc4fa3bf816695f743c98f@o952367.ingest.sentry.io/6177866")
}
