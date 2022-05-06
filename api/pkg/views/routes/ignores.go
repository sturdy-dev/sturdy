package routes

import (
	"net/http"
	"os"

	db_view "getsturdy.com/api/pkg/views/db"
	"getsturdy.com/api/pkg/views/ignore"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Ignores(logger *zap.Logger, executorProvider executor.Provider, viewRepo db_view.Repository) func(*gin.Context) {
	return func(c *gin.Context) {
		view, err := viewRepo.Get(c.Param("viewID"))
		if err != nil {
			logger.Error("failed to get view", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		var res []string
		err = executorProvider.New().
			AllowRebasingState(). // allowed to get .gitignore even if conflicting
			Read(func(repo vcs.RepoReader) error {
				var err error
				res, err = ignore.FindIgnore(os.DirFS(repo.Path()))
				if err != nil {
					return err
				}
				return nil
			}).ExecView(view.CodebaseID, view.ID, "findIgnores")
		if err != nil {
			logger.Error("failed to find ignores", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"paths": res,
		})
	}
}
