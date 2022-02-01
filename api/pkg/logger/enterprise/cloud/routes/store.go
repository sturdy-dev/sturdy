package routes

import (
	"encoding/json"
	"net/http"

	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
)

func Store(logger *zap.Logger, sentryClient *raven.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		packet := &raven.Packet{}
		if err := json.NewDecoder(r.Body).Decode(packet); err != nil {
			logger.Error("failed to decode sentry packet", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		sentryClient.Capture(packet, nil)
	}
}
