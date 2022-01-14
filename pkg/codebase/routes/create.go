package routes

import (
	"net/http"
	"strings"
	"time"

	"mash/pkg/auth"
	"mash/pkg/shortid"
	"mash/pkg/view/events"
	service_workspace "mash/pkg/workspace/service"
	"mash/vcs/executor"
	"mash/vcs/provider"

	"github.com/google/uuid"

	"github.com/posthog/posthog-go"
	"go.uber.org/zap"

	"mash/pkg/codebase"
	"mash/pkg/codebase/db"
	"mash/pkg/codebase/vcs"

	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func Create(
	logger *zap.Logger,
	repo db.CodebaseRepository,
	codebaseUserRepo db.CodebaseUserRepository,
	executorProvider executor.Provider,
	postHogClient posthog.Client,
	eventsSender events.EventSender,
	workspaceService service_workspace.Service,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req CreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Warn("failed to bind input", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}
		req.Name = strings.TrimSpace(req.Name)
		if req.Name == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		userID, err := auth.UserID(c.Request.Context())
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		cb, err := DoCreateCodebase(
			logger,
			uuid.NewString(),
			req.Name,
			req.Description,
			repo,
			executorProvider,
		)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Add user to codebase
		t := time.Now()
		err = codebaseUserRepo.Create(codebase.CodebaseUser{
			ID:         uuid.New().String(),
			UserID:     userID,
			CodebaseID: cb.ID,
			CreatedAt:  &t,
		})
		if err != nil {
			logger.Error("failed to add user to codebase", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		err = postHogClient.Enqueue(posthog.Capture{
			DistinctId: userID,
			Event:      "create codebase",
			Properties: posthog.NewProperties().
				Set("codebase_id", cb.ID).
				Set("name", cb.Name),
		})
		if err != nil {
			logger.Error("posthog failed", zap.Error(err))
		}

		if err := workspaceService.CreateWelcomeWorkspace(cb.ID, userID, cb.Name); err != nil {
			logger.Error("failed to create welcome workspace", zap.Error(err))
			// not a critical error, continue
		}

		// Send events
		if err := eventsSender.Codebase(cb.ID, events.CodebaseUpdated, cb.ID); err != nil {
			logger.Error("failed send events", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, cb)
	}
}

func DoCreateCodebase(
	logger *zap.Logger,
	codebaseID string,
	name, description string,
	repo db.CodebaseRepository,
	executorProvider executor.Provider,
) (codebase.Codebase, error) {
	t := time.Now()
	cb := codebase.Codebase{
		ID:              codebaseID,
		ShortCodebaseID: codebase.ShortCodebaseID(shortid.New()),
		Name:            name,
		Description:     description,
		Emoji:           "🌟",
		CreatedAt:       &t,
		IsReady:         true, // No additional setup needed
	}

	// Create codebase in database
	if err := repo.Create(cb); err != nil {
		logger.Error("failed to create codebase", zap.Error(err))
		return codebase.Codebase{}, err
	}

	if err := executorProvider.New().
		AllowRebasingState(). // allowed because the repo does not exist yet
		Schedule(func(trunkProvider provider.RepoProvider) error {
			return vcs.Create(trunkProvider, cb.ID)
		}).ExecTrunk(cb.ID, "createCodebase"); err != nil {
		logger.Error("failed to create codebase repo", zap.Error(err))
		return codebase.Codebase{}, err
	}

	return cb, nil
}
