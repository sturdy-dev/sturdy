package routes

import (
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
)

func Store(logger *zap.Logger, sentryClient *raven.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		switch contentType {
		case "application/json":
			packet := &raven.Packet{}
			if err := json.NewDecoder(r.Body).Decode(packet); err != nil {
				logger.Error("failed to decode sentry packet", zap.Error(err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r.Body.Close()
			sentryClient.Capture(packet, nil)
		case "application/octet-stream":
			b64 := base64.NewDecoder(base64.StdEncoding, r.Body)
			deflate, _ := zlib.NewReader(b64)
			packet := &raven.Packet{}
			if err := json.NewDecoder(deflate).Decode(packet); err != nil {
				logger.Error("failed to decode sentry packet", zap.Error(err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			deflate.Close()
			r.Body.Close()
			sentryClient.Capture(packet, nil)
		default:
			logger.Error("failed to decode sentry packet: unexpected content-type", zap.String("content-type", contentType))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
