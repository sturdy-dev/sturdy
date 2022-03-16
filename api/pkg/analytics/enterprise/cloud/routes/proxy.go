package routes

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/posthog/posthog-go"
	"go.uber.org/zap"
)

// Batch implements posthog /batch/ api endpoint
func Batch(logger *zap.Logger, client posthog.Client) gin.HandlerFunc {
	type event struct {
		Type       string         `json:"type"`
		Event      string         `json:"event"`
		Timestamp  time.Time      `json:"timestamp"`
		Properties map[string]any `json:"properties"`
		Groups     map[string]any `json:"groups"`
		DistinctID string         `json:"distinct_id"`
		Set        map[string]any `json:"$set"`
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
			delete(event.Properties, "$lib")
			delete(event.Properties, "$lib_version")
			if event.Event == "$groupidentify" {
				key := fmt.Sprint(event.Properties["$group_key"])
				typ := fmt.Sprint(event.Properties["$group_type"])
				delete(event.Properties, "$group_key")
				delete(event.Properties, "$group_type")
				if err := client.Enqueue(posthog.GroupIdentify{
					Key:        key,
					Type:       typ,
					DistinctId: event.DistinctID,
					Timestamp:  event.Timestamp,
					Properties: event.Properties,
				}); err != nil {
					logger.Error("failed to enqueue group", zap.Error(err))
				}
			} else {
				switch event.Type {
				case "identify":
					if err := client.Enqueue(posthog.Identify{
						DistinctId: event.DistinctID,
						Properties: event.Set,
						Timestamp:  event.Timestamp,
					}); err != nil {
						logger.Error("failed to enqueue identify", zap.Error(err))
					}
				case "capture":
					if err := client.Enqueue(posthog.Capture{
						Event:      event.Event,
						DistinctId: event.DistinctID,
						Properties: event.Properties,
						Groups:     event.Groups,
						Timestamp:  event.Timestamp,
					}); err != nil {
						logger.Error("failed to enqueue identify", zap.Error(err))
					}
				default:
					logger.Error("unknown event type", zap.String("type", event.Type))
				}
			}
		}
	}
}
