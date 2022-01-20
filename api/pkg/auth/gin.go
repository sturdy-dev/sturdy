package auth

import (
	"errors"
	"fmt"
	"net/http"

	"getsturdy.com/api/pkg/ctxlog"
	"getsturdy.com/api/pkg/jwt"
	service_jwt "getsturdy.com/api/pkg/jwt/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	ginContextKey = "auth.subject"
)

func GinMiddleware(logger *zap.Logger, jwtService *service_jwt.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, shouldRefresh, err := jwtFromRequest(c.Request, jwtService)
		if err != nil && !errors.Is(err, ErrUnauthenticated) {
			ctxlog.ErrorOrWarn(logger, "failed to authenticate user", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if shouldRefresh {
			if err := refreshToken(c, token, jwtService); err != nil {
				ctxlog.ErrorOrWarn(logger, "failed to refresh token", err)
			}
		}

		subject := subjectFromToken(token)

		c.Set(ginContextKey, subject)
		c.Request = c.Request.WithContext(NewContext(c.Request.Context(), subject))

		c.Next()
	}
}

func refreshToken(c *gin.Context, token *jwt.Token, jwtService *service_jwt.Service) error {
	token, err := jwtService.IssueToken(c.Request.Context(), token.Subject, oneMonth, token.Type)
	if err != nil {
		return fmt.Errorf("failed to issue new token: %w", err)
	}

	isSecure := c.Request.URL.Scheme == "https"
	setAuthCookie(c.Writer, isSecure, token.Token)
	return nil
}

func SetAuthCookieForUser(c *gin.Context, userID string, jwtService *service_jwt.Service) error {
	token, err := jwtService.IssueToken(c.Request.Context(), userID, oneMonth, jwt.TokenTypeAuth)
	if err != nil {
		return fmt.Errorf("failed to issue new token: %w", err)
	}

	isSecure := c.Request.URL.Scheme == "https"
	setAuthCookie(c.Writer, isSecure, token.Token)
	return nil
}

func SubjectFromGinContext(c *gin.Context) (*Subject, bool) {
	v, found := c.Get(ginContextKey)
	if !found {
		return nil, false
	}

	subject, ok := v.(*Subject)
	return subject, ok
}
