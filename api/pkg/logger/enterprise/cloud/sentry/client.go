package sentry

import (
	"context"

	"getsturdy.com/api/pkg/version"

	"github.com/getsentry/sentry-go"
)

func NewClient() (*sentry.Client, error) {
	return sentry.NewClient(sentry.ClientOptions{
		Dsn:        "https://7ca135fcbbbc4fa3bf816695f743c98f@o952367.ingest.sentry.io/6177866",
		ServerName: version.Type.String(),
		Release:    version.Version,
		IgnoreErrors: []string{
			context.Canceled.Error(),
			"pq: canceling statement due to user request",
		},
	})
}
