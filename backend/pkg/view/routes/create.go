package routes

import (
	"net/http"
	"time"

	"mash/pkg/analytics"
	"mash/pkg/auth"
	"mash/pkg/codebase/access"
	db_codebase "mash/pkg/codebase/db"
	db_snapshots "mash/pkg/snapshots/db"
	"mash/pkg/snapshots/snapshotter"
	"mash/pkg/view"
	"mash/pkg/view/db"
	"mash/pkg/view/events"
	"mash/pkg/view/open"
	"mash/pkg/view/vcs"
	db_workspace "mash/pkg/workspace/db"
	"mash/vcs/executor"
	"mash/vcs/provider"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CreateRequest struct {
	CodebaseID    string `json:"codebase_id" binding:"required"`
	WorkspaceID   string `json:"workspace_id" binding:"required"`
	Name          string `json:"name"`
	MountPath     string `json:"mount_path"`
	MountHostname string `json:"mount_hostname"`
}

func Create(
	logger *zap.Logger,
	viewRepo db.Repository,
	codebaseUserRepo db_codebase.CodebaseUserRepository,
	analyticsClient analytics.Client,
	workspaceReader db_workspace.WorkspaceReader,
	snapshotter snapshotter.Snapshotter,
	snapshotRepo db_snapshots.Repository,
	workspaceWriter db_workspace.WorkspaceWriter,
	executorProvider executor.Provider,
	eventSender events.EventSender,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req CreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Warn("failed to parse input", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to parse or validate input"})
			return
		}

		userID, err := auth.UserID(c.Request.Context())
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !access.UserHasAccessToCodebase(codebaseUserRepo, userID, req.CodebaseID) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		t := time.Now()
		e := view.View{
			ID:            uuid.New().String(),
			UserID:        userID,
			CodebaseID:    req.CodebaseID,
			WorkspaceID:   req.WorkspaceID,
			Name:          &req.Name,          // It's optional
			MountPath:     &req.MountPath,     // It's optional
			MountHostname: &req.MountHostname, // It's optional
			CreatedAt:     &t,
		}

		if err := viewRepo.Create(e); err != nil {
			logger.Error("failed to create view", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to create view"})
			return
		}

		if err := executorProvider.New().
			AllowRebasingState(). // allowed because the view does not exist yet
			Schedule(func(repoProvider provider.RepoProvider) error {
				return vcs.Create(repoProvider, req.CodebaseID, req.WorkspaceID, e.ID)
			}).ExecView(req.CodebaseID, e.ID, "createView"); err != nil {
			logger.Error("failed to create view", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to create view"})
			return
		}

		ws, err := workspaceReader.Get(req.WorkspaceID)
		if err != nil {
			logger.Error("failed to get workspace", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Use workspace on view
		if err := open.OpenWorkspaceOnView(logger, &e, ws, viewRepo, workspaceReader, snapshotter, snapshotRepo, workspaceWriter, executorProvider, eventSender); err != nil {
			logger.Error("failed to open workspace on view", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if err := analyticsClient.Enqueue(analytics.Capture{
			DistinctId: userID,
			Event:      "create view",
			Properties: analytics.NewProperties().
				Set("codebase_id", req.CodebaseID).
				Set("workspace_id", req.WorkspaceID).
				Set("view_id", e.ID).
				Set("mount_path", e.MountPath).
				Set("mount_hostname", e.MountHostname),
		}); err != nil {
			logger.Error("analytics failed", zap.Error(err))
			// do not fail
		}

		c.JSON(http.StatusOK, e)
	}
}
