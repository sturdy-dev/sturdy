package routes

import (
	"database/sql"
	"errors"
	"net/http"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/codebases/access"
	"getsturdy.com/api/pkg/codebases/db"
	service_user "getsturdy.com/api/pkg/users/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Get(
	repo db.CodebaseRepository,
	codebaseUserRepo db.CodebaseUserRepository,
	logger *zap.Logger,
	userService service_user.Service,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		userID, err := auth.UserID(c.Request.Context())
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var cb *codebases.Codebase

		// If ID is a short ID
		if len(id) == 7 {
			cb, err = repo.GetByShortID(codebases.ShortCodebaseID(id))
		} else {
			cb, err = repo.Get(codebases.ID(id))
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

		memberAuthors, err := membersAsAuthors(c.Request.Context(), codebaseUserRepo, userService, cb.ID)
		if err != nil {
			logger.Error("failed to get members", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, codebases.CodebaseWithMetadata{
			Codebase: *cb,
			Members:  memberAuthors,
		})
	}
}
