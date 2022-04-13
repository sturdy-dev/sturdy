package pprof

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"getsturdy.com/api/pkg/pprof/configuration"
)

type Server struct {
	cfg *configuration.Configuration
}

func New(cfg *configuration.Configuration) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) Start() error {
	if err := http.ListenAndServe(s.cfg.Addr.String(), nil); err != http.ErrServerClosed {
		return fmt.Errorf("failed to start http pprof server: %w", err)
	}
	return nil
}
