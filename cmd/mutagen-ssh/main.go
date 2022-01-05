package main

import (
	"context"
	"flag"
	"mash/pkg/mutagen/ssh"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	sshListenAddr := flag.String("ssh-listen-addr", "127.0.0.1:2222", "")
	sturdyApiAddr := flag.String("sturdy-api-addr", "http://host.docker.internal:3000", "")
	mutagenAgentBinaryDir := flag.String("mutagen-agent-binary-dir", "/usr/bin/", "")
	keyHostPath := flag.String("ssh-key-path", "id_ed25519", "")
	httpPprofListenAddr := flag.String("http-pprof-listen-addr", "127.0.0.1:6060", "")
	flag.Parse()

	logger, _ := zap.NewProduction()

	srv := ssh.New(logger, &ssh.Config{
		ListenAddr:            *sshListenAddr,
		KeyHostPath:           *keyHostPath,
		MutagenAgentBinaryDir: *mutagenAgentBinaryDir,
		SturdyApiAddr:         *sturdyApiAddr,
	})

	// Pprof server
	go func() {
		logger.Info("Starting pprof server", zap.String("listen-addr", *httpPprofListenAddr))
		if err := http.ListenAndServe(*httpPprofListenAddr, nil); err != http.ErrServerClosed {
			logger.Fatal("pprof server error", zap.Error(err))
		}
	}()

	// Wait for shutdown in a separate goroutine.
	errCh := make(chan error)
	go func() {
		shutdownCh := make(chan os.Signal, 1)
		signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)
		sig := <-shutdownCh

		logger.Info("shutting down", zap.Stringer("signal", sig))

		shutdownTimeout := 15 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		errCh <- srv.Shutdown(shutdownCtx)
	}()

	mainCtx := context.Background()
	if err := srv.Start(mainCtx); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}

	// Handle shutdown errors.
	if err := <-errCh; err != nil {
		logger.Fatal("failed to shutdown server", zap.Error(err))
	}

	logger.Info("server shutdown")
}
