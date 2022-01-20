package gitserver

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	service_codebase "getsturdy.com/api/pkg/codebase/service"
	"getsturdy.com/api/pkg/gitserver/pack"
	"getsturdy.com/api/pkg/jwt"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	"getsturdy.com/api/pkg/servicetokens"
	service_servicetokens "getsturdy.com/api/pkg/servicetokens/service"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger

	serviceTokensService *service_servicetokens.Service
	jwtTokensService     *service_jwt.Service
	codebaseService      *service_codebase.Service
	executorProvider     executor.Provider

	router *gin.Engine
}

func New(
	logger *zap.Logger,
	serviceTokensService *service_servicetokens.Service,
	jwtTokensService *service_jwt.Service,
	codebaeService *service_codebase.Service,
	executorProvider executor.Provider,
) *Server {
	return &Server{
		logger: logger,

		serviceTokensService: serviceTokensService,
		jwtTokensService:     jwtTokensService,
		codebaseService:      codebaeService,
		executorProvider:     executorProvider,

		router: gin.New(),
	}
}

func (h *Server) Start(ctx context.Context, addr string) error {
	h.router.Use(ginzap.Ginzap(h.logger, time.RFC3339, true))
	h.router.Use(ginzap.RecoveryWithZap(h.logger, true))

	ciIntegrationGroup := h.router.Group("/").Use(h.serviceTokenAuth)
	ciIntegrationGroup.GET("/info/refs", h.handleInfoRefs)
	ciIntegrationGroup.POST("/git-upload-pack", h.handleGitUploadPack)

	importGroup := h.router.Group("/:codebaseId").Use(h.jwtTokenAuth)
	importGroup.GET("/info/refs", h.handleInfoRefs)
	importGroup.POST("/git-receive-pack", h.handleGitReceivePack)

	h.logger.Info("start gitserver", zap.String("addr", addr))

	if err := h.router.Run(addr); err != http.ErrServerClosed {
		return fmt.Errorf("failed to run the server: %w", err)
	}

	return nil
}

const (
	tokenKey  = "token"
	userIDKey = "user_id"
	ciRepo    = "ci"
)

func (h *Server) jwtTokenAuth(c *gin.Context) {
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		c.Header("WWW-Authenticate", "Basic realm=Authorization Required")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if username != "import" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userToken, err := h.jwtTokensService.Verify(c.Request.Context(), password, jwt.TokenTypeAuth)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	accessAllowed, err := h.codebaseService.CanAccess(c.Request.Context(), userToken.Subject, c.Param("codebaseId"))
	if err != nil {
		h.logger.Error("failed to check access", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if !accessAllowed {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	c.Set(userIDKey, userToken.Subject)
}

func (h *Server) serviceTokenAuth(c *gin.Context) {
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		c.Header("WWW-Authenticate", "Basic realm=Authorization Required")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := h.serviceTokensService.Get(c.Request.Context(), username)
	if errors.Is(err, sql.ErrNoRows) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	} else if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := token.Verify(password); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set(tokenKey, token)
}

func getToken(c *gin.Context) *servicetokens.Token {
	token, ok := c.Get(tokenKey)
	if !ok {
		return nil
	}

	return token.(*servicetokens.Token)
}

func getServiceName(r *http.Request) string {
	if service, fromQuery := r.URL.Query()["service"]; fromQuery {
		return strings.Replace(service[0], "git-", "", 1)
	}

	if len(r.Form["service"]) > 0 {
		return strings.Replace(r.Form["service"][0], "git-", "", 1)
	}

	return ""
}

func (h *Server) handleGitReceivePack(c *gin.Context) {
	codebaseID := c.Param("codebaseId")

	c.Writer.Header().Set("Content-Type", "application/x-git-receive-pack-result")

	// TODO: This buffers the whole request (including _all_ object data etc), in memory.
	// Rewrite so that the header can be parsed of a stream (r.Body) as it's read by ProcInput?
	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error("receive-pack failed to read request", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var firstRow []byte
	if idx := bytes.Index(requestBody, []byte("\n")); idx > 0 {
		firstRow = requestBody[0:idx]
	} else {
		firstRow = requestBody
	}

	// TODO: Learn more about git-receive-pack and the protocol to figure out why the headers are not always sent
	//       and how to keep doing this in a safe way.

	if string(firstRow) != "0000" {
		// If we have a header, parse it!
		header, err := pack.ParseHeader(firstRow)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			h.logger.Error("receive-pack failed to parse header", zap.Error(err))
			return
		}

		if header.Branch != "sturdytrunk" {
			c.AbortWithStatus(http.StatusBadRequest)
			h.logger.Error("receive-pack request to non sturdytrunk branch",
				zap.String("branch", header.Branch),
				zap.String("codebase_id", codebaseID),
			)
			return
		}
	}

	if err := h.executorProvider.New().Write(func(repo vcs.RepoWriter) error {
		args := []string{"receive-pack", "--stateless-rpc", repo.Path()}
		cmd := exec.Command("git", args...)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("failed to get stdout: %w", err)
		}
		cmd.Stdin = bytes.NewReader(requestBody)

		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start command: %w", err)
		}

		if _, err := io.Copy(c.Writer, stdout); err != nil {
			return fmt.Errorf("failed to copy stdout: %w", err)
		}

		return nil
	}).ExecTrunk(codebaseID, "gitserverGitReceivePack"); err != nil {
		h.logger.Error("failed to handle git receive pack", zap.Error(err))
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (h *Server) handleGitUploadPack(c *gin.Context) {
	token := getToken(c)

	c.Writer.Header().Set("Content-Type", "application/x-git-upload-pack-result")

	if err := h.executorProvider.New().Read(func(repo vcs.RepoReader) error {
		args := []string{"upload-pack", "--stateless-rpc", repo.Path()}
		cmd := exec.Command("git", args...)
		cmd.Stdin = c.Request.Body
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("failed to get stdout: %w", err)
		}

		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start command: %w", err)
		}

		if _, err := io.Copy(c.Writer, stdout); err != nil {
			return fmt.Errorf("failed to copy stdout: %w", err)
		}

		return nil
	}).ExecView(token.CodebaseID, ciRepo, "gitserverGitUploadPack"); err != nil {
		h.logger.Error("failed to handle git upload pack", zap.Error(err))
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (h *Server) handleInfoRefs(c *gin.Context) {
	// todo: what is this?
	serviceName := getServiceName(c.Request)

	c.Writer.WriteHeader(200)
	c.Header("Content-Type", fmt.Sprintf("application/x-git-%s-advertisement", serviceName))

	str := fmt.Sprintf("# service=git-%s", serviceName)
	fmt.Fprintf(c.Writer, "%.4x%s\n", len(str)+5, str)
	fmt.Fprintf(c.Writer, "0000")

	executor := h.executorProvider.New().Read(func(repo vcs.RepoReader) error {
		args := []string{serviceName, "--stateless-rpc", "--advertise-refs", repo.Path()}
		cmd := exec.Command("git", args...)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("failed to get stdout: %w", err)
		}

		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start command: %w", err)
		}

		if _, err := io.Copy(c.Writer, stdout); err != nil {
			return fmt.Errorf("failed to copy stdout: %w", err)
		}

		return nil
	})

	if token := getToken(c); token != nil { // this is ci flow
		if err := executor.ExecView(token.CodebaseID, ciRepo, "gitserverInfoRefs"); err != nil {
			h.logger.Error("failed to handle info refs", zap.Error(err))
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	} else { // this is import flow
		if err := executor.ExecTrunk(c.Param("codebaseId"), "gitserverInfoRefs"); err != nil {
			h.logger.Error("failed to handle info refs", zap.Error(err))
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
}
