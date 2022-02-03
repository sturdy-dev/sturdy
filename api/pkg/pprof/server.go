package pprof

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"getsturdy.com/api/pkg/configuration/flags"
	"getsturdy.com/api/pkg/di"
)

type Configuration struct {
	Addr flags.Addr `long:"addr" description:"address to listen on" default:"127.0.0.1:6060"`
}

type Server struct {
	cfg *Configuration
}

func New(cfg *Configuration) *Server {
	return &Server{
		cfg: cfg,
	}
}

func Module(c *di.Container) {
	c.Register(New)
}

func (s *Server) Start() error {
	if err := http.ListenAndServe(s.cfg.Addr.String(), nil); err != http.ErrServerClosed {
		return fmt.Errorf("failed to start http pprof server: %w", err)
	}
	return nil
}
