package ssh

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Config struct {
	ListenAddr            string
	SturdyApiAddr         string
	KeyHostPath           string
	MutagenAgentBinaryDir string
}

type Server struct {
	logger *zap.Logger
	cfg    *Config

	sshServer *ssh.Server
}

func New(logger *zap.Logger, cfg *Config) *Server {
	return &Server{
		logger: logger,
		cfg:    cfg,
	}
}

func (srv *Server) Start(ctx context.Context) error {
	srv.logger.Info("Starting SSH server",
		zap.String("addr", srv.cfg.ListenAddr),
	)

	srv.sshServer = &ssh.Server{
		Handler: srv.sshHandler,
		Addr:    srv.cfg.ListenAddr,

		// On the mutagen client side: -oConnectTimeout=5 -oServerAliveInterval=10 -oServerAliveCountMax=1

		IdleTimeout: time.Second * 30,

		// TODO(zegl): If this is _set_ all rsync transfers are forced to complete within 45s, which is not reasonable for
		// large file transfers.
		// This is disabled for now.
		// MaxTimeout:  time.Second * 45,
	}
	if err := srv.sshServer.SetOption(ssh.HostKeyFile(srv.cfg.KeyHostPath)); err != nil {
		return fmt.Errorf("failed to set host key file: %w", err)
	}
	if err := srv.sshServer.SetOption(ssh.PublicKeyAuth(validateKey(srv.logger, srv.cfg.SturdyApiAddr))); err != nil {
		return fmt.Errorf("failed to set public key auth: %w", err)
	}

	if err := srv.sshServer.SetOption(ssh.WrapConn(func(_ ssh.Context, conn net.Conn) net.Conn {
		if err := conn.(*net.TCPConn).SetKeepAlive(true); err != nil {
			log.Fatal(err)
		}
		if err := conn.(*net.TCPConn).SetKeepAlivePeriod(time.Second * 10); err != nil {
			log.Fatal(err)
		}
		return conn
	})); err != nil {
		return fmt.Errorf("failed to set tcp connection option: %w", err)
	}

	if err := srv.sshServer.ListenAndServe(); err != ssh.ErrServerClosed {
		return fmt.Errorf("failed to listen and serve: %w", err)
	}

	srv.logger.Info("SSH server shut down")

	return nil
}

func (srv *Server) Shutdown(ctx context.Context) error {
	srv.logger.Info("Shutting down SSH server")
	if err := srv.sshServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown SSH server: %w", err)
	}
	return nil
}

func (srv *Server) sshHandler(s ssh.Session) {
	logger := srv.logger.With(
		zap.String("user", s.User()),
		zap.Strings("command", s.Command()),
		zap.String("connection_id", uuid.NewString()),
	)

	var binary string
	switch s.Command()[0] {
	case ".sturdy-sync/agents/0.12.0-beta2/mutagen-agent":
		binary = path.Join(srv.cfg.MutagenAgentBinaryDir, "mutagen-agent-v0.12.0-beta2")
	case ".sturdy-sync/agents/0.12.0-beta6/mutagen-agent":
		binary = path.Join(srv.cfg.MutagenAgentBinaryDir, "mutagen-agent-v0.12.0-beta6")
	case ".sturdy-sync/agents/0.12.0-beta7/mutagen-agent":
		binary = path.Join(srv.cfg.MutagenAgentBinaryDir, "mutagen-agent-v0.12.0-beta7")
	case ".sturdy-sync/agents/0.13.0-beta2/mutagen-agent":
		binary = path.Join(srv.cfg.MutagenAgentBinaryDir, "mutagen-agent-v0.13.0-beta2")
	default:
		logger.Error("connection with unknown binary")
		return
	}

	logger = logger.With(zap.String("binary", binary))
	logger.Info("SSH connection")

	t0 := time.Now()

	cmd := exec.Command(binary, "synchronizer")

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("STURDY_AUTHENTICATED_USER_ID=%s", s.User()))
	cmd.Env = append(cmd.Env, fmt.Sprintf("STURDY_API_ADDR=%s", srv.cfg.SturdyApiAddr))

	stdin, err := cmd.StdinPipe()
	if err != nil {
		logger.Error("failed to create pipe", zap.Error(err))
		return
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error("failed to create pipe", zap.Error(err))
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Error("failed to create pipe", zap.Error(err))
		return
	}

	if err := cmd.Start(); err != nil {
		logger.Error("failed to start agent", zap.Error(err))
		return
	}

	// stdout
	go func() {
		_, err := io.Copy(s, stdout)
		if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, os.ErrClosed) {
			logger.Error("stdout copy failed", zap.Error(err))
		}
	}()

	// stderr
	se := s.Stderr()
	pr, pw := io.Pipe()
	go func() {
		teeStdeer := io.TeeReader(stderr, pw)
		if _, err := io.Copy(se, teeStdeer); err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, os.ErrClosed) {
			logger.Error("stderr copy failed", zap.Error(err))
		}
	}()

	go func() {
		b := bufio.NewReader(pr)
		for {
			ln, _, err := b.ReadLine()
			if err != nil {
				logger.Error("stderr", zap.String("output", string(ln)), zap.Error(err))
				break
			}
			logger.Info("stderr", zap.String("output", string(ln)))
		}
	}()

	// stdin
	if _, err := io.Copy(stdin, s); err != nil {
		if errIsConnectionClosed(err) {
			logger.Warn("stdin copy failed", zap.Error(err))
		} else {
			logger.Error("stdin copy failed", zap.Error(err))
		}
	}

	// This point is reached once stdin (from SSH) is closed.
	// Close stdin to the program to cause it to self-terminate
	if err := stdin.Close(); err != nil {
		logger.Error("failed to close stdin", zap.Error(err))
	}

	if err := cmd.Wait(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok && ee.ExitCode() == 1 {
			// this is fine
		} else {
			logger.Error("program wait failed", zap.Error(err))
		}
	}

	logger.Info("Disconnected SSH connection",
		zap.Duration("connection_duration", time.Since(t0)),
	)
}

type VerifyPublicKeyRequest struct {
	// In what format?
	PublicKey []byte `json:"public_key"`
	UserID    string `json:"user_id"`
}

func validateKey(logger *zap.Logger, sturdyApiAddr string) func(ctx ssh.Context, key ssh.PublicKey) bool {
	return func(ctx ssh.Context, key ssh.PublicKey) bool {
		var res struct{}
		err := Request(sturdyApiAddr, "POST", "/v3/pki/verify", "", &VerifyPublicKeyRequest{
			UserID:    ctx.User(),
			PublicKey: key.Marshal(),
		}, &res)

		if err != nil {
			if err.Error() != "unexpected response code 404" {
				logger.Error("key verification failed", zap.Error(err))
			}
			return false
		}
		return true
	}
}

func errIsConnectionClosed(err error) bool {
	if err == nil {
		return false
	}
	if strings.Contains(err.Error(), "connection reset by peer") {
		return true
	}
	if strings.Contains(err.Error(), "broken pipe") {
		return true
	}
	return false
}
