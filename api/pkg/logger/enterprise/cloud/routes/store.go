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
	type request struct {
		raven.Packet

		Exception *struct {
			raven.Exception
			Values []*raven.Exception `json:"values"`
		} `json:"exception"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		packet := &request{}
		switch contentType {
		case "application/json":
			if err := json.NewDecoder(r.Body).Decode(packet); err != nil {
				logger.Error("failed to decode sentry packet", zap.Error(err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r.Body.Close()
		case "application/octet-stream":
			b64 := base64.NewDecoder(base64.StdEncoding, r.Body)
			deflate, _ := zlib.NewReader(b64)
			if err := json.NewDecoder(deflate).Decode(packet); err != nil {
				logger.Error("failed to decode sentry packet", zap.Error(err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			deflate.Close()
			r.Body.Close()
		default:
			logger.Error("failed to decode sentry packet: unexpected content-type", zap.String("content-type", contentType))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if packet.Exception != nil {
			packet.Interfaces = append(packet.Interfaces, &packet.Exception.Exception)
			if len(packet.Exception.Values) > 0 {
				packet.Interfaces = append(packet.Interfaces, &raven.Exceptions{
					Values: packet.Exception.Values,
				})
			}
		}

		sentryClient.Capture(&packet.Packet, nil)
	}
}
