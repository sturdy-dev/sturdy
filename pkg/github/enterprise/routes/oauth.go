package routes

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"mash/pkg/auth"
	service_github "mash/pkg/github/service"
	"net/http"
	"time"

	"mash/pkg/github"
	"mash/pkg/github/config"
	"mash/pkg/github/db"
	db_user "mash/pkg/user/db"

	"github.com/gin-gonic/gin"
	gh "github.com/google/go-github/v39/github"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type GitHubAuthReq struct {
	Code string `json:"code" binding:"required"`
}

func Oauth(
	logger *zap.Logger,
	config config.GitHubAppConfig,
	userRepo db_user.Repository,
	gitHubUserRepo db.GitHubUserRepo,
	gitHubService *service_github.Service,
) func(*gin.Context) {
	return func(c *gin.Context) {
		var incomingReq GitHubAuthReq
		if err := c.ShouldBindJSON(&incomingReq); err != nil {
			logger.Error("failed to parse request", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to parse or validate input"})
			return
		}

		userID, err := auth.UserID(c.Request.Context())
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// This user already have a github account connected
		// Do a quick refresh of the permissions, and return early
		if _, err := gitHubUserRepo.GetByUserID(userID); err == nil {
			err = gitHubService.AddUserToCodebases(context.Background(), userID)
			if err != nil {
				logger.Error("failed to update codebase access for already authed user", zap.Error(err))
				c.Status(http.StatusInternalServerError)
				return
			}

			c.Status(http.StatusOK)
			return
		} else if !errors.Is(err, sql.ErrNoRows) {
			logger.Error("failed to read github_users db", zap.Error(err))
			c.Status(http.StatusInternalServerError)
			return
		}

		accessToken, err := getAccessToken(logger, config, incomingReq.Code)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		client := gh.NewClient(tc)

		ghApiUser, _, err := client.Users.Get(ctx, "")
		if err != nil {
			logger.Error("failed to get github user from api", zap.Error(err))
			c.Status(http.StatusInternalServerError)
			return
		}

		user, err := userRepo.Get(userID)
		if err != nil {
			logger.Error("failed to get user by ID provided in the state query param", zap.Error(err))
			c.Status(http.StatusInternalServerError)
			return
		}

		ghUser, err := gitHubUserRepo.GetByUsername(ghApiUser.GetLogin())
		if errors.Is(err, sql.ErrNoRows) {
			t0 := time.Now()
			ghUser = &github.GitHubUser{
				ID:          uuid.NewString(),
				UserID:      user.ID,
				Username:    ghApiUser.GetLogin(),
				AccessToken: accessToken,
				CreatedAt:   t0,
			}
			err := gitHubUserRepo.Create(*ghUser)
			if err != nil {
				logger.Error("failed to create github user in db", zap.Error(err))
				c.Status(http.StatusInternalServerError)
				return
			}
		} else if err != nil {
			logger.Error("failed to lookup github user repo in db", zap.Error(err))
			c.Status(http.StatusInternalServerError)
			return
		}

		if err := gitHubService.AddUserToCodebases(ctx, user.ID); err != nil {
			logger.Error("failed to grant user access", zap.Error(err))
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
		return
	}
}

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}

func getAccessToken(logger *zap.Logger, config config.GitHubAppConfig, code string) (token string, err error) {
	httpClient := http.Client{}

	// Next, lets for the HTTP request to call the github oauth endpoint
	// to get our access token
	reqURL := fmt.Sprintf("%s?client_id=%s&client_secret=%s&code=%s",
		"https://github.com/login/oauth/access_token/",
		config.GitHubAppClientID,
		config.GitHubAppSecret,
		code)

	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		logger.Error("failed to build oauth request", zap.Error(err))
		return "", err
	}
	// We set this header since we want the response
	// as JSON
	req.Header.Set("accept", "application/json")

	// Send out the HTTP request
	res, err := httpClient.Do(req)
	if err != nil {
		logger.Error("failed to call endpoint", zap.Error(err))
		return "", err
	}
	defer res.Body.Close()

	// Parse the request body into the `OAuthAccessResponse` struct
	var t OAuthAccessResponse
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		logger.Error("failed to parse response", zap.Error(err))
		return "", err
	}

	return t.AccessToken, nil
}
