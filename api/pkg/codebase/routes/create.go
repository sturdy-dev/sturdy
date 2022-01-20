package routes

import (
	"net/http"
	"strings"

	"go.uber.org/zap"

	service_codebase "getsturdy.com/api/pkg/codebase/service"

	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"` // TODO: remove this field, it's unused
}

func Create(
	logger *zap.Logger,
	svc *service_codebase.Service,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req CreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Warn("failed to bind input", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}
		req.Name = strings.TrimSpace(req.Name)
		if req.Name == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		cb, err := svc.Create(c.Request.Context(), req.Name, nil)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, cb)
	}
}
