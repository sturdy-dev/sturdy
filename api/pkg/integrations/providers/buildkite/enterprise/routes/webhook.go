package routes

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	svc_ci "getsturdy.com/api/pkg/ci/service"
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/integrations/providers/buildkite"
	service_buildkite "getsturdy.com/api/pkg/integrations/providers/buildkite/enterprise/service"
	service_servicetokens "getsturdy.com/api/pkg/servicetokens/service"
	"getsturdy.com/api/pkg/statuses"
	svc_statuses "getsturdy.com/api/pkg/statuses/service"

	bk "github.com/buildkite/go-buildkite/v3/buildkite"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	acceptEvents = map[string]bool{
		"build.scheduled": true,
		"build.running":   true,
		"build.finished":  true,
		"job.scheduled":   true,
		"job.running":     true,
		"job.finished":    true,
		"job.activated":   true,
	}

	// Valid states: running, scheduled, passed, failed, blocked, canceled, canceling, skipped, not_run, finished
	buildkiteStateToType = map[string]statuses.Type{
		"running":   statuses.TypePending,
		"blocked":   statuses.TypePending,
		"canceling": statuses.TypePending,
		"scheduled": statuses.TypePending,

		"passed":   statuses.TypeHealty,
		"skipped":  statuses.TypeHealty,
		"not_run":  statuses.TypeHealty,
		"finished": statuses.TypeHealty,

		"failed":   statuses.TypeFailing,
		"canceled": statuses.TypeFailing,
	}
)

var errInvalidSignature = errors.New("invalid signature")

var allowedWindow = 5 * time.Minute

// The X-Buildkite-Signature header contains a timestamp and an HMAC signature.
// The timestamp is prefixed by timestamp= and the signature is prefixed by signature=.
// e.g. timestamp=1637075221,signature=dbdabe3596995f7bd1f39f50f135df4c48e4291f5368c0eb5c5a02664ae536e9
//
// Buildkite generates the signature using HMAC-SHA256; a hash-based message authentication code HMAC used with
// the SHA-256 hash function and a secret key. The webhook token value is used as the secret key. The timestamp
// is an integer representation of a UTC timestamp.
func parseSignatureHeader(xBuildkiteSignature string) (timestamp time.Time, signature string, err error) {
	params := strings.Split(xBuildkiteSignature, ",")
	for _, param := range params {
		kv := strings.Split(param, "=")
		switch kv[0] {
		case "timestamp":
			ts, err := strconv.Atoi(kv[1])
			if err != nil {
				return time.Time{}, "", fmt.Errorf("invalid timestamp: %w", err)
			}
			timestamp = time.Unix(int64(ts), 0)
		case "signature":
			signature = kv[1]
		}
	}
	return
}

func WebhookHandler(
	logger *zap.Logger,
	statusesService *svc_statuses.Service,
	ciService *svc_ci.Service,
	serviceTokensService *service_servicetokens.Service,
	buildkiteService *service_buildkite.Service,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("failed to read body", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		payload := &webhookPayload{}
		if err := json.Unmarshal(requestBody, payload); err != nil {
			logger.Error("failed to parse payload", zap.Error(err))
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to parse payload"))
			return
		}

		// Short-circuit events that we're not interested in
		if !acceptEvents[payload.Event] {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		if payload.Pipeline.Repository == nil {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		logger := logger.With(
			zap.String("pipeline_repository", *payload.Pipeline.Repository),
			zap.String("x-buildkite-signature", c.GetHeader("X-Buildkite-Signature")),
			zap.String("x-buildkite-event", c.GetHeader("X-Buildkite-Event")),
		)

		// Parse the pipeline URL, to extract service token
		pipelineUrl, err := url.Parse(*payload.Pipeline.Repository)
		if err != nil {
			logger.Error("invalid URL from buildkite", zap.Error(err))
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("malformed pipeline repository url"))
			return
		}
		serviceTokenID := pipelineUrl.User.Username()

		serviceToken, err := serviceTokensService.Get(c.Request.Context(), serviceTokenID)
		if errors.Is(err, sql.ErrNoRows) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		} else if err != nil {
			logger.Error("failed to get service token", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if err := validateSignature(c.Request.Context(), c.GetHeader("X-Buildkite-Signature"), serviceToken.CodebaseID, requestBody, buildkiteService); err != nil {
			if errors.Is(err, errInvalidSignature) {
				logger.Error("failed to validate signature", zap.Error(err))
				c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to validate signature"))
				return
			} else {
				logger.Error("failed to validate signature", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

		trunkCommitSHA, err := ciService.GetTrunkCommitSHA(c.Request.Context(), serviceToken.CodebaseID, *payload.Build.Commit)
		if errors.Is(err, sql.ErrNoRows) {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("unknown commit: %s", *payload.Build.Commit))
			return
		} else if err != nil {
			logger.Error("could not find trunk commit", zap.Stringp("buildkite_build_commit", payload.Build.Commit))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		state := payload.State()
		statusType, ok := buildkiteStateToType[state]
		if !ok {
			logger.Error("invalid status from buildkite", zap.String("status", state))
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid status: %s", state))
			return
		}
		webURL := payload.WebURL()
		if err := statusesService.Set(c, &statuses.Status{
			ID:         uuid.NewString(),
			CommitSHA:  trunkCommitSHA,
			CodebaseID: serviceToken.CodebaseID,
			Type:       statusType,
			Title:      payload.Title(),
			DetailsURL: &webURL,
			Timestamp:  time.Now(),
		}); err != nil {
			logger.Error("failed to update status", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		logger.Info("got webhook from buildkite", zap.String("path", pipelineUrl.Path), zap.Stringer("codebase_id", serviceToken.CodebaseID))
	}
}

func validateSignature(ctx context.Context, xBuildkiteSignature string, codebaseID codebases.ID, requestBody []byte, buildkiteService *service_buildkite.Service) error {
	buildkiteConfigs, err := buildkiteService.GetConfigurationsByCodebaseID(ctx, codebaseID)
	if err != nil {
		return fmt.Errorf("failed to get buildkite configuration: %w", err)
	}

	var lastError error

	for _, cfg := range buildkiteConfigs {
		if err := validateSingleSignature(cfg, xBuildkiteSignature, requestBody); err != nil {
			lastError = err
		} else if err == nil {
			// Successfully validated
			return nil
		}
	}

	// Unexpected
	if lastError == nil {
		return fmt.Errorf("failed to validate buildkite signature (unexpected no success)")
	}

	return lastError
}

func validateSingleSignature(buildkiteCfg *buildkite.Config, xBuildkiteSignature string, requestBody []byte) error {
	timestamp, signature, err := parseSignatureHeader(xBuildkiteSignature)
	if err != nil {
		return errInvalidSignature
	}

	now := time.Now()
	if timestamp.After(now) {
		return errInvalidSignature
	}
	if timestamp.Before(now.Add(-1 * allowedWindow)) {
		return errInvalidSignature
	}

	hmacSum := hmac.New(sha256.New, []byte(buildkiteCfg.WebhookSecret))
	if _, err := hmacSum.Write([]byte(fmt.Sprint(timestamp.Unix()))); err != nil {
		return fmt.Errorf("failed to write timestamp to sha: %w", err)
	}
	if _, err := hmacSum.Write([]byte(".")); err != nil {
		return fmt.Errorf("failed to write dot to sha: %w", err)
	}
	if _, err := hmacSum.Write(requestBody); err != nil {
		return fmt.Errorf("failed to write request body to sha: %w", err)
	}

	signatureValid := fmt.Sprintf("%x", hmacSum.Sum(nil)) == signature
	if !signatureValid {
		return errInvalidSignature
	}

	return nil
}

type webhookPayload struct {
	Event    string      `json:"event"`
	Build    bk.Build    `json:"build"`
	Job      *bk.Job     `json:"job"`
	Pipeline bk.Pipeline `json:"pipeline"`
}

func (p *webhookPayload) WebURL() string {
	if p.Job != nil {
		return p.Job.WebURL
	}
	return *p.Build.WebURL
}

var buildkiteEmoji = regexp.MustCompile(`:[a-z0-9_]+:`)

func sanitize(s string) string {
	s = buildkiteEmoji.ReplaceAllString(s, "")
	s = strings.TrimSpace(s)
	return s
}

func (p *webhookPayload) Title() string {
	if p.Job != nil {
		return fmt.Sprintf("%s: %s", sanitize(*p.Pipeline.Name), sanitize(*p.Job.Name))
	}
	return sanitize(*p.Pipeline.Name)
}

func (p *webhookPayload) State() string {
	if p.Job != nil {
		return *p.Job.State
	}
	return *p.Build.State
}
