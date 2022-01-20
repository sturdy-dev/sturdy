package routes

import (
	"io"
	"net/http"

	service_auth "getsturdy.com/api/pkg/auth/service"
	service_suggestions "getsturdy.com/api/pkg/suggestions/service"
	view2 "getsturdy.com/api/pkg/view"
	"getsturdy.com/api/pkg/view/db"
	"getsturdy.com/api/pkg/view/events"
	"getsturdy.com/api/pkg/view/stream"
	db_workspace "getsturdy.com/api/pkg/workspace/db"
	service_workspace "getsturdy.com/api/pkg/workspace/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Stream(
	logger *zap.Logger,
	viewRepo db.Repository,
	viewEventsReader events.EventReader,
	workspaceReader db_workspace.WorkspaceReader,
	authService *service_auth.Service,
	workspaceService service_workspace.Service,
	suggestionService *service_suggestions.Service,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		logger := logger

		workspaceID := c.Query("workspace_id")
		if len(workspaceID) == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "workspace_id not set"})
			return
		}
		ws, err := workspaceReader.Get(workspaceID)
		if err != nil {
			logger.Error("workspace not found", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "workspace not found"})
			return
		}

		if err := authService.CanRead(c.Request.Context(), ws); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Optional
		viewID := c.Query("view_id")

		logger = logger.With(
			zap.String("workspace_id", workspaceID),
			zap.String("view_id", viewID),
		)

		var view *view2.View
		if viewID != "" {
			view, err = viewRepo.Get(viewID)
			if err != nil {
				logger.Error("view not found", zap.Error(err))
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "view not found"})
				return
			}
		}

		done := make(chan bool)

		chanStream, err := stream.Stream(c.Request.Context(), logger, ws, view, done, viewEventsReader, authService, workspaceService, suggestionService)
		if err != nil {
			logger.Error("stream failed", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to stream changes"})
			return
		}

		// Blocks as long as the request is open
		c.Stream(func(w io.Writer) bool {
			if ev, ok := <-chanStream; ok {
				c.SSEvent(string(ev.Name), ev.Message)
				return true
			}
			return false
		})

		// Indicate to worker that it should stop
		done <- true
	}
}
