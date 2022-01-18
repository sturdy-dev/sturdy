package routes

import (
	"time"

	"mash/pkg/analytics"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Batch implements posthog /batch/ api endpoint
func Batch(logger *zap.Logger, client analytics.Client) gin.HandlerFunc {
	type event struct {
		Type       string                 `json:"type"`
		Event      string                 `json:"event"`
		Timestamp  time.Time              `json:"timestamp"`
		Properties map[string]interface{} `json:"properties"`
		DistinctID string                 `json:"distinct_id"`
		Set        map[string]interface{} `json:"$set"`
	}
	type request struct {
		Batch []event `json:"batch"`
	}

	logger = logger.Named("analytics.proxy")

	return func(c *gin.Context) {
		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		for _, event := range req.Batch {
			switch event.Type {
			case "identify":
				if err := client.Enqueue(analytics.Identify{
					DistinctId: event.DistinctID,
					Properties: event.Set,
					Timestamp:  event.Timestamp,
				}); err != nil {
					logger.Error("failed to enqueue identify", zap.Error(err))
					return
				}
			case "capture":
				delete(event.Properties, "$lib")
				delete(event.Properties, "$lib_version")
				if err := client.Enqueue(analytics.Capture{
					Event:      event.Event,
					DistinctId: event.DistinctID,
					Properties: event.Properties,
					Timestamp:  event.Timestamp,
				}); err != nil {
					logger.Error("failed to enqueue identify", zap.Error(err))
					return
				}
			default:
				logger.Error("unknown event type", zap.String("type", event.Type))
				return
			}
		}
	}
}
