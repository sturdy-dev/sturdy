package routes

import (
	"database/sql"
	"errors"
	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/codebase/access"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
	db_view "getsturdy.com/api/pkg/view/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MutagenView struct {
	ID                 string `json:"id"`
	UserID             string `json:"user_id"`
	CodebaseID         string `json:"codebase_id"`
	CodebaseName       string `json:"codebase_name"`
	CodebaseIsArchived bool   `json:"codebase_is_archived"`
}

func GetView(logger *zap.Logger, repo db_view.Repository, codebaseUserRepo db_codebase.CodebaseUserRepository, codebaseRepo db_codebase.CodebaseRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		viewID := c.Param("id")

		viewObj, err := repo.Get(viewID)
		if err != nil {
			logger.Error("failed to get view", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		userID, err := auth.UserID(c.Request.Context())
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Request must be made by the owner of the view
		if viewObj.UserID != userID {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !access.UserHasAccessToCodebase(codebaseUserRepo, userID, viewObj.CodebaseID) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		res := MutagenView{
			ID:     viewObj.ID,
			UserID: viewObj.UserID,
		}

		codebaseObj, err := codebaseRepo.Get(viewObj.CodebaseID)
		if errors.Is(err, sql.ErrNoRows) {
			res.CodebaseIsArchived = true
		} else if err != nil {
			logger.Error("failed to get codebase", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		} else {
			res.CodebaseID = codebaseObj.ID
			res.CodebaseName = codebaseObj.Name
		}
		c.JSON(http.StatusOK, res)
	}
}
