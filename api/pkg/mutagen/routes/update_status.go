package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	eventsv2 "getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/mutagen"
	"getsturdy.com/api/pkg/mutagen/db"
	db_view "getsturdy.com/api/pkg/view/db"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type stateStatus struct {
	Identifier     string         `json:"identifier,omitempty"`
	Name           string         `json:"name,omitempty"`
	Status         string         `json:"status,omitempty"`
	AlphaURL       string         `json:"alphaUrl,omitempty"`
	BetaURL        string         `json:"betaUrl,omitempty"`
	AlphaConnected bool           `json:"alphaConnected"`
	BetaConnected  bool           `json:"betaConnected"`
	LastError      string         `json:"lastError,omitempty"`
	StagingStatus  receiverStatus `json:"stagingStatus,omitempty"`
	AlphaProblems  []problem      `json:"alphaProblems,omitempty"`
	BetaProblems   []problem      `json:"betaProblems,omitempty"`
	Paused         bool           `json:"paused"`
	SturdyVersion  string         `json:"sturdyVersion,omitempty"`
}

type problem struct {
	Path  string `json:"path,omitempty"`
	Error string `json:"error,omitempty"`
}

type receiverStatus struct {
	Path     string `json:"path,omitempty"`
	Received uint64 `json:"received,omitempty"`
	Total    uint64 `json:"total,omitempty"`
}

// This endpoint is unauthenticated
// TODO: Find a way to propagate the auth context to the mutagen client
func UpdateStatus(logger *zap.Logger, viewStatusRepo db.ViewStatusRepository, viewRepo db_view.Repository, eventsSender *eventsv2.Publisher) func(*gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var input stateStatus
		if err := c.BindJSON(&input); err != nil {
			logger.Warn("failed to read status", zap.Error(err))
			c.Status(http.StatusBadRequest)
			return
		}

		logger.Info("mutagen status", zap.Any("state", input))

		prefix := "view-"
		if !strings.HasPrefix(input.Name, prefix) {
			c.Status(http.StatusBadRequest)
			return
		}
		viewID := input.Name[len(prefix):]

		status, err := viewStatusRepo.GetByViewID(viewID)
		if errors.Is(err, sql.ErrNoRows) {
			status := mutagen.ViewStatus{ID: viewID}
			setStatus(&status, input)
			err = viewStatusRepo.Create(status)
			if err != nil {
				logger.Error("failed to create status", zap.Error(err))
				c.Status(http.StatusInternalServerError)
				return
			}
			c.Status(http.StatusOK)
			return
		} else if err != nil {
			logger.Error("failed to get view status", zap.Error(err))
			c.Status(http.StatusInternalServerError)
			return
		}

		setStatus(status, input)

		err = viewStatusRepo.Update(status)
		if err != nil {
			logger.Error("failed to update status", zap.Error(err))
			c.Status(http.StatusInternalServerError)
			return
		}

		// Send event
		vw, err := viewRepo.Get(viewID)
		if err != nil {
			logger.Error("failed to get view", zap.Error(err))
			c.Status(http.StatusInternalServerError)
			return
		}

		if err := eventsSender.Codebase(ctx, vw.CodebaseID).ViewStatusUpdated(vw); err != nil {
			logger.Error("failed to send event", zap.Error(err))
			// do not fail
		}

		c.JSON(http.StatusOK, struct{}{})
	}
}

func setStatus(status *mutagen.ViewStatus, input stateStatus) {
	status.State = mutagen.ViewStatusState(input.Status)
	status.SturdyVersion = input.SturdyVersion

	if input.LastError != "" {
		status.LastError = &input.LastError
	} else {
		status.LastError = nil
	}

	if input.StagingStatus.Path != "" {
		status.StagingStatusPath = &input.StagingStatus.Path
		r := int(input.StagingStatus.Received)
		status.StagingStatusReceived = &r
		t := int(input.StagingStatus.Total)
		status.StagingStatusTotal = &t
	} else {
		status.StagingStatusPath = nil
		status.StagingStatusReceived = nil
		status.StagingStatusTotal = nil
	}
}

/*
	{
		"identifier": "sync_277Pz5lrIUpQRbbzlK8edubAXGpXU79oZlb61rP0SgB",
		"alphaUrl": "path:\"C:\\\\Users\\\\owend\\\\Desktop\\\\ue4-tribeproject\"",
		"betaUrl": "protocol:SSH user:\"9ccbc0e1-937a-4b7c-a823-cff42abae04d\" host:\"sync.getsturdy.com\" path:\"/repos/588f9cf7-8ab5-49b1-8d3e-061a4b581fa9/6de0f1ac-6de7-42bd-b2ec-f3ae2e9bfae1/\"",
		"status": "StagingBeta",
		"alphaConnected": true,
		"betaConnected": true,
		"stagingStatus": {
			"path": "Content/GraduallyRottenFoodVol1/Textures/Pear/pear_08_NRM.uasset",
			"received": 4518,
			"total": 8104
		},
		"sturdyVersion": "v0.5.17"
	}
*/

/*
	{
		"identifier": "sync_eV6PeO79BVUwfKxWkBgVVfIob1ggpOHwY2GiY3URBmp",
		"alphaUrl": "path:\"/Users/gustav/src/sturdy\"",
		"betaUrl": "protocol:SSH  user:\"847dfd0c-49bf-40c5-8870-74a12fca0d60\"  host:\"sync.getsturdy.com\"  path:\"/repos/31596772-e9d6-445e-8144-856a3022744b/757787ca-0928-4c37-bdae-d9b2e0a09f4c/\"",
		"status": "WaitingForRescan",
		"alphaConnected": true,
		"betaConnected": true,
		"lastError": "alpha scan error: hashed size mismatch (cmd/sturdy/sturdy-v0.5.18-darwin-amd64.tar.gz): 3686400 != 3563520",
		"stagingStatus": {},
		"sturdyVersion": "v0.5.17"
	}
*/

/*
	{
	        "identifier": "sync_277Pz5lrIUpQRbbzlK8edubAXGpXU79oZlb61rP0SgB",
	        "alphaUrl": "path:\"C:\\\\Users\\\\owend\\\\Desktop\\\\ue4-tribeproject\"",
	        "betaUrl": "protocol:SSH user:\"9ccbc0e1-937a-4b7c-a823-cff42abae04d\" host:\"sync.getsturdy.com\" path:\"/repos/588f9cf7-8ab5-49b1-8d3e-061a4b581fa9/6de0f1ac-6de7-42bd-b2ec-f3ae2e9bfae1/\"",
	        "status": "WaitingForRescan",
	        "alphaConnected": true,
	        "betaConnected": true,
	        "lastError": "alpha scan error: unable to open file (Content/Level.umap): unable to open file handle: unable to open path: The process cannot access the file because it is being used by another process.",
	        "stagingStatus": {},
	        "sturdyVersion": "v0.5.17"
	    }
*/
