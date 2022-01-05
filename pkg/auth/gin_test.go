package auth_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"mash/pkg/auth"
	"mash/pkg/jwt"
	db_jwt_keys "mash/pkg/jwt/keys/db"
	service_jwt "mash/pkg/jwt/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	oneDay   = 24 * time.Hour
	oneMonth = 30 * oneDay
)

func TestGinMiddleware__shouldAllowCIAuthInHeader(t *testing.T) {
	jwtTokenService := service_jwt.NewService(zap.NewNop(), db_jwt_keys.NewInMemory())

	token, err := jwtTokenService.IssueToken(context.Background(), "id", oneMonth, jwt.TokenTypeCI)
	assert.NoError(t, err)

	router := gin.New()
	router.Use(auth.GinMiddleware(zap.NewNop(), jwtTokenService))
	router.GET("/ping", func(c *gin.Context) {
		subject, found := auth.SubjectFromGinContext(c)
		if assert.True(t, found) {
			assert.Equal(t, subject.ID, "id")
			assert.Equal(t, subject.Type, auth.SubjectCI)
		}

		subject, found = auth.FromContext(c.Request.Context())
		if assert.True(t, found) {
			assert.Equal(t, subject.ID, "id")
			assert.Equal(t, subject.Type, auth.SubjectCI)
		}
	})

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/ping", nil)
	assert.NoError(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", token.Token))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Len(t, w.Result().Cookies(), 0)
}

func TestGinMiddleware__shouldAllowUserAuthInHeader(t *testing.T) {
	jwtTokenService := service_jwt.NewService(zap.NewNop(), db_jwt_keys.NewInMemory())

	token, err := jwtTokenService.IssueToken(context.Background(), "id", oneMonth, jwt.TokenTypeAuth)
	assert.NoError(t, err)

	router := gin.New()
	router.Use(auth.GinMiddleware(zap.NewNop(), jwtTokenService))
	router.GET("/ping", func(c *gin.Context) {
		subject, found := auth.SubjectFromGinContext(c)
		if assert.True(t, found) {
			assert.Equal(t, subject.ID, "id")
			assert.Equal(t, subject.Type, auth.SubjectUser)
		}

		subject, found = auth.FromContext(c.Request.Context())
		if assert.True(t, found) {
			assert.Equal(t, subject.ID, "id")
			assert.Equal(t, subject.Type, auth.SubjectUser)
		}
	})

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/ping", nil)
	assert.NoError(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", token.Token))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Len(t, w.Result().Cookies(), 0)
}

func TestGinMiddleware__shouldNotRefreshExpiringHeaderToken(t *testing.T) {
	jwtTokenService := service_jwt.NewService(zap.NewNop(), db_jwt_keys.NewInMemory())

	token, err := jwtTokenService.IssueToken(context.Background(), "id", time.Hour, jwt.TokenTypeAuth)
	assert.NoError(t, err)

	router := gin.New()
	router.Use(auth.GinMiddleware(zap.NewNop(), jwtTokenService))
	router.GET("/ping", func(c *gin.Context) {
		subject, found := auth.SubjectFromGinContext(c)
		if assert.True(t, found) {
			assert.Equal(t, subject.ID, "id")
			assert.Equal(t, subject.Type, auth.SubjectUser)
		}

		subject, found = auth.FromContext(c.Request.Context())
		if assert.True(t, found) {
			assert.Equal(t, subject.ID, "id")
			assert.Equal(t, subject.Type, auth.SubjectUser)
		}
	})

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/ping", nil)
	assert.NoError(t, err)
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", token.Token))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Len(t, w.Result().Cookies(), 0)
}

func TestGinMiddleware__shouldRefreshExpiringCookie(t *testing.T) {
	jwtTokenService := service_jwt.NewService(zap.NewNop(), db_jwt_keys.NewInMemory())

	token, err := jwtTokenService.IssueToken(context.Background(), "id", time.Hour, jwt.TokenTypeAuth)
	assert.NoError(t, err)

	router := gin.New()
	router.Use(auth.GinMiddleware(zap.NewNop(), jwtTokenService))
	router.GET("/ping", func(c *gin.Context) {
		subject, found := auth.SubjectFromGinContext(c)
		if assert.True(t, found) {
			assert.Equal(t, subject.ID, "id")
			assert.Equal(t, subject.Type, auth.SubjectUser)
		}

		subject, found = auth.FromContext(c.Request.Context())
		if assert.True(t, found) {
			assert.Equal(t, subject.ID, "id")
			assert.Equal(t, subject.Type, auth.SubjectUser)
		}
	})

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/ping", nil)
	assert.NoError(t, err)
	req.AddCookie(&http.Cookie{
		Name:  "auth",
		Value: token.Token,
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	if cookies := w.Result().Cookies(); assert.Len(t, cookies, 1) {
		assert.Equal(t, "auth", cookies[0].Name)
		assert.NotEqual(t, token.Token, cookies[0].Value, "cookie was not refreshed")
	}
}

func TestGinMiddleware__shouldAllowUserAuthInCookie(t *testing.T) {
	jwtTokenService := service_jwt.NewService(zap.NewNop(), db_jwt_keys.NewInMemory())

	token, err := jwtTokenService.IssueToken(context.Background(), "id", oneMonth, jwt.TokenTypeAuth)
	assert.NoError(t, err)

	router := gin.New()
	router.Use(auth.GinMiddleware(zap.NewNop(), jwtTokenService))
	router.GET("/ping", func(c *gin.Context) {
		subject, found := auth.SubjectFromGinContext(c)
		if assert.True(t, found) {
			assert.Equal(t, subject.ID, "id")
			assert.Equal(t, subject.Type, auth.SubjectUser)
		}

		subject, found = auth.FromContext(c.Request.Context())
		if assert.True(t, found) {
			assert.Equal(t, subject.ID, "id")
			assert.Equal(t, subject.Type, auth.SubjectUser)
		}
	})

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/ping", nil)
	assert.NoError(t, err)
	req.AddCookie(&http.Cookie{
		Name:  "auth",
		Value: token.Token,
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Len(t, w.Result().Cookies(), 0)
}

func TestGinMiddleware__shouldAllowNoAuth(t *testing.T) {
	jwtTokenService := service_jwt.NewService(zap.NewNop(), db_jwt_keys.NewInMemory())

	router := gin.New()
	router.Use(auth.GinMiddleware(zap.NewNop(), jwtTokenService))
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/ping", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Len(t, w.Result().Cookies(), 0)
	assert.Equal(t, "pong", w.Body.String())
}
