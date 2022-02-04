package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/codebase"
	"getsturdy.com/api/pkg/codebase/access"
	"getsturdy.com/api/pkg/codebase/db"
	"getsturdy.com/api/pkg/codebase/vcs"
	db_user "getsturdy.com/api/pkg/users/db"
	vcsvcs "getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func Get(
	repo db.CodebaseRepository,
	codebaseUserRepo db.CodebaseUserRepository,
	logger *zap.Logger,
	userRepo db_user.Repository,
	executorProvider executor.Provider,
) gin.HandlerFunc {
	lastUpdatedAt := func(codebaseID string) time.Time {
		var gitTime time.Time
		if err := executorProvider.New().GitRead(func(repo vcsvcs.RepoGitReader) error {
			changes, err := vcs.ListChanges(repo, 1)
			if err != nil || len(changes) == 0 {
				return fmt.Errorf("failed to list changes: %w", err)
			}
			gitTime = changes[0].Time
			return nil
		}).ExecTrunk(codebaseID, "codebase.LastUpdatedAt"); err != nil {
			return time.Unix(0, 0)
		}

		return gitTime
	}
	return func(c *gin.Context) {
		id := c.Param("id")

		userID, err := auth.UserID(c.Request.Context())
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var cb *codebase.Codebase

		// If ID is a short ID
		if len(id) == 7 {
			cb, err = repo.GetByShortID(id)
		} else {
			cb, err = repo.Get(id)
		}

		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("codebase not found", zap.String("codebase_id", id), zap.Error(err))
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		if err != nil {
			logger.Error("failed to get codebase", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unable to get codebase"})
			return
		}

		// Since this API allows lookup by the ShortID, this access check (which uses the UUID) is done after fetching the codebase
		if !access.UserHasAccessToCodebase(codebaseUserRepo, userID, cb.ID) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		memberAuthors, err := membersAsAuthors(codebaseUserRepo, userRepo, cb.ID)
		if err != nil {
			logger.Error("failed to get members", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, codebase.CodebaseWithMetadata{
			Codebase:          *cb,
			LastUpdatedAtUnix: lastUpdatedAt(cb.ID).Unix(),
			Members:           memberAuthors,
		})
	}
}
