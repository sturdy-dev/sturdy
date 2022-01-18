package routes

import (
	"context"
	"database/sql"
	"errors"
	"mash/pkg/auth"
	"mash/pkg/codebase"
	"mash/pkg/ctxlog"
	"net/http"

	service_auth "mash/pkg/auth/service"
	"mash/pkg/codebase/acl"
	"mash/pkg/user"
	"mash/pkg/view"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type aclProvider interface {
	GetByCodebaseID(context.Context, string) (acl.ACL, error)
}

type viewRepository interface {
	Get(string) (*view.View, error)
}

type userRpository interface {
	Get(id string) (*user.User, error)
}

func ListAllows(
	logger *zap.Logger,
	viewRepo viewRepository,
	authService *service_auth.Service,
) func(*gin.Context) {
	type listAllowsResponse struct {
		Allows []string `json:"allows"`
	}

	return func(c *gin.Context) {
		viewID := c.Param("id")

		viewObj, err := viewRepo.Get(viewID)
		switch {
		case err == nil:
		case errors.Is(err, sql.ErrNoRows):
			c.AbortWithStatus(http.StatusNotFound)
			return
		default:
			ctxlog.ErrorOrWarn(logger, "failed to get view", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// TODO: Authenticate internal requests
		ctx := auth.NewContext(c.Request.Context(), &auth.Subject{ID: viewObj.UserID, Type: auth.SubjectMutagen})

		allower, err := authService.GetAllower(ctx, &codebase.Codebase{ID: viewObj.CodebaseID})
		if err != nil {
			ctxlog.ErrorOrWarn(logger, "failed to list allowed pattern", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, &listAllowsResponse{
			Allows: allower.Patterns,
		})
	}
}
