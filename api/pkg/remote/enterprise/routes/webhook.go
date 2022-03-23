package routes

import (
	"context"
	"fmt"
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

		if c.Request.Method != "POST" {
			c.Status(http.StatusBadRequest)
			_, _ = c.Writer.WriteString(fmt.Sprintf("Hey! Send a POST request to this endpoint to activate the magic. (got a %s-request)", c.Request.Method))
			return
		}

		codebaseID := codebases.ID(c.Param("id"))

		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			_, _ = c.Writer.WriteString("Unexpected ID")
			return
		}
		defer c.Request.Body.Close()

		logger.Info("received hook",
			zap.Stringer("codebase_id", codebaseID),
			zap.String("body", string(body)))

		ctx := context.Background()
		if err := svc.Pull(ctx, codebaseID); err != nil {
			c.Status(http.StatusInternalServerError)
			_, _ = c.Writer.WriteString("InternalServerError, please try again later...")
			return
		}

		c.Status(http.StatusOK)
		_, _ = c.Writer.WriteString("OK!")
		return
	}
}
