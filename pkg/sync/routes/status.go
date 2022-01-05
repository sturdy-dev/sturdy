package routes

import (
	"errors"
	"fmt"
	"net/http"

	"mash/pkg/sync"
	"mash/pkg/sync/service"
	db_view "mash/pkg/view/db"
	"mash/vcs"
	"mash/vcs/executor"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Status does not perform any authentication. That's expected to be done by ValidateViewAccessMiddleware or similar.
func Status(repo db_view.Repository, executorProvider executor.Provider, logger *zap.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("viewID")
		view, err := repo.Get(id)
		if err != nil {
			logger.Error("failed to get view", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		var status *sync.RebaseStatusResponse
		if err := executorProvider.New().AllowRebasingState().Write(func(repo vcs.RepoWriter) error {
			rebasing, err := repo.OpenRebase()
			if err != nil {
				if errors.Is(err, vcs.NoRebaseInProgress) {
					status = &sync.RebaseStatusResponse{}
					return nil
				}
				return fmt.Errorf("failed to open rebase: %w", err)
			}

			rebaseStatus, err := service.Status(logger, rebasing)
			if err != nil {
				return fmt.Errorf("failed to get status: %w", err)
			}
			status = rebaseStatus
			return nil
		}).ExecView(view.CodebaseID, view.ID, "rebaseStatus"); err != nil {
			logger.Error("failed to get rebase status", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, status)
	}
}
