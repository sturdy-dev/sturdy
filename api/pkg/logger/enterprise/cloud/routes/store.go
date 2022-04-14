package routes

import (
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

func parseTimestamp(msg json.RawMessage) (time.Time, error) {
	var t string
	if err := json.Unmarshal(msg, &t); err != nil {
		return time.Time{}, err
	}
	if t, err := time.Parse(time.RFC3339Nano, t); err == nil {
		return t, nil
	}
	return time.Parse("2006-01-02T15:04:05.00", t)
}

func parseException(msg json.RawMessage) ([]sentry.Exception, error) {
	var exceptions []sentry.Exception
	if err := json.Unmarshal(msg, &exceptions); err == nil {
		return exceptions, err
	}
	var exception sentry.Exception
	if err := json.Unmarshal(msg, &exception); err != nil {
		return nil, err
	}
	return []sentry.Exception{exception}, nil
}

// this route handles raven-go's sentry event
func Store(logger *zap.Logger, sentryClient *sentry.Client) http.HandlerFunc {
	type e struct {
		*sentry.Event

		// Raven library sends only a single exception, while in Sentry it's a slice.
		// Override the exception field to handle that difference.
		Exception json.RawMessage `json:"exception,omitempty"`

		// Raven library sends timestmap in "2006-01-02T15:04:05.00" format, while in Sentry it's RFC3339Nano.
		// Override the timestamp field to handle that difference.
		Timestamp json.RawMessage `json:"timestamp,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		event := &e{}
		switch contentType {
		case "application/json":
			if err := json.NewDecoder(r.Body).Decode(event); err != nil {
				logger.Error("failed to decode sentry packet", zap.Error(err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r.Body.Close()
		case "application/octet-stream":
			b64 := base64.NewDecoder(base64.StdEncoding, r.Body)
			deflate, _ := zlib.NewReader(b64)
			if err := json.NewDecoder(deflate).Decode(event); err != nil {
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

		// now, handle format differences
		ts, err := parseTimestamp(event.Timestamp)
		if err != nil {
			logger.Error("failed to parse timestamp", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		exceptions, err := parseException(event.Exception)
		if err != nil {
			logger.Error("failed to parse exception", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sentryEvent := event.Event
		sentryEvent.Timestamp = ts
		sentryEvent.Exception = exceptions

		// send event!
		sentryClient.CaptureEvent(event.Event, nil, nil)
	}
}
