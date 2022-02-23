package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/github/enterprise/config"
	"getsturdy.com/api/pkg/github/enterprise/db"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	db_user "getsturdy.com/api/pkg/users/db"

	"github.com/gin-gonic/gin"
	gh "github.com/google/go-github/v39/github"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	githubOAuth "golang.org/x/oauth2/github"
)

func Oauth(
	logger *zap.Logger,
	config *config.GitHubAppConfig,
	userRepo db_user.Repository,
	gitHubUserRepo db.GitHubUserRepo,
	gitHubService *service_github.Service,
) func(*gin.Context) {
	type GitHubAuthReq struct {
		Code string `json:"code" binding:"required"`
	}
	oauthCfg := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.Secret,
		Endpoint:     githubOAuth.Endpoint,
	}
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

		logger = logger.With(zap.String("user_id", userID))

		if ghUser, err := gitHubUserRepo.GetByUserID(userID); errors.Is(err, sql.ErrNoRows) {
			// This user doesn't have a github connected yet, connect it.
		} else if err != nil {
			logger.Error("failed to get user", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		} else {
			// If user has a github connected already, refresh permissions without asking for a new token
			if err := gitHubService.AddUserToCodebases(c.Request.Context(), ghUser); err != nil {
				logger.Error("failed to add user to codebases", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			return
		}

		user, err := userRepo.Get(userID)
		if err != nil {
			logger.Error("failed to get user by ID provided in the state query param", zap.Error(err))
			c.Status(http.StatusNotFound)
			return
		}

		token, err := oauthCfg.Exchange(c.Request.Context(), incomingReq.Code)
		if err != nil {
			logger.Error("failed to exchange code for access token", zap.Error(err))
			c.Status(http.StatusInternalServerError)
			return
		}

		githubOAuth2Client := oauthCfg.Client(c.Request.Context(), token)
		githubAPIClient := gh.NewClient(githubOAuth2Client)

		ghApiUser, resp, err := githubAPIClient.Users.Get(c.Request.Context(), "")
		if err != nil {
			if resp.StatusCode == http.StatusUnauthorized {
				c.JSON(http.StatusBadRequest, gin.H{"error": "GitHub authentication failed: invalid code provided"})
				return
			} else {
				logger.Error("failed to get github user from api", zap.Error(err))
				c.Status(http.StatusInternalServerError)
				return
			}
		}

		ghUser, err := gitHubUserRepo.GetByUsername(ghApiUser.GetLogin())
		if errors.Is(err, sql.ErrNoRows) {
			ghUser = &github.GitHubUser{
				ID:          uuid.NewString(),
				UserID:      user.ID,
				Username:    ghApiUser.GetLogin(),
				AccessToken: token.AccessToken,
				CreatedAt:   time.Now(),
			}
			if err := gitHubUserRepo.Create(*ghUser); err != nil {
				logger.Error("failed to create github user in db", zap.Error(err))
				c.Status(http.StatusInternalServerError)
				return
			}
		} else if err != nil {
			logger.Error("failed to lookup github user repo in db", zap.Error(err))
			c.Status(http.StatusInternalServerError)
			return
		} else if ghUser.UserID != user.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("This GitHub account is already used by another Sturdy user (%s)", user.Email)})
			return
		}

		if err := gitHubService.AddUserToCodebases(c.Request.Context(), ghUser); err != nil {
			logger.Error("failed to grant user access", zap.Error(err))
			c.Status(http.StatusInternalServerError)
			return
		}
	}
}
