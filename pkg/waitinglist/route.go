package waitinglist

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/posthog/posthog-go"
	"go.uber.org/zap"
)

type WaitingListRequest struct {
	Email string `json:"email" binding:"required"`
}

func Insert(logger *zap.Logger, postHogClient posthog.Client, repo WaitingListRepo) func(c *gin.Context) {
	return func(c *gin.Context) {
		logger := logger

		var req WaitingListRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to parse or validate input"})
			return
		}

		req.Email = strings.TrimSpace(req.Email)

		if !strings.Contains(req.Email, "@") {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
			return
		}

		logger = logger.With(zap.String("email", req.Email))

		err := repo.Insert(req.Email)
		if err != nil {
			logger.Error("failed to add to waitinglist", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		err = postHogClient.Enqueue(&posthog.Capture{
			DistinctId: req.Email,
			Event:      "signed up for waiting list",
		})
		if err != nil {
			logger.Error("posthog failed", zap.Error(err))
		}

		logger.Info("added to waitinglist")
		c.Status(http.StatusOK)
	}
}
