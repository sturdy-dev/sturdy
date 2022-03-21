package routes

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/remote/enterprise/service"
)

type TriggerSyncCodebaseWebhookHandler gin.HandlerFunc

func TriggerSyncCodebaseWebhook(svc *service.EnterpriseService, logger *zap.Logger) TriggerSyncCodebaseWebhookHandler {
	logger = logger.Named("TriggerSyncCodebaseWebhookHandler")
	return func(c *gin.Context) {
		logger := logger

		codebaseID := codebases.ID(c.Param("id"))

		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		defer c.Request.Body.Close()

		logger.Info("received hook",
			zap.Stringer("codebase_id", codebaseID),
			zap.String("body", string(body)))

		ctx := context.Background()
		if err := svc.Pull(ctx, codebaseID); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		c.AbortWithStatus(http.StatusOK)
	}
}
