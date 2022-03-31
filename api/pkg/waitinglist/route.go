package waitinglist

import (
	"log"
	"net/http"
	"strings"

	"getsturdy.com/api/pkg/analytics"
	service_analytics "getsturdy.com/api/pkg/analytics/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WaitingListRequest struct {
	Email string `json:"email" binding:"required"`
}

func Insert(logger *zap.Logger, analyticsService *service_analytics.Service, repo WaitingListRepository) func(c *gin.Context) {
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

		analyticsService.Capture(c.Request.Context(), "signed up for waiting list",
			analytics.DistinctID(req.Email),
		)

		logger.Info("added to waitinglist")
		c.Status(http.StatusOK)
	}
}
