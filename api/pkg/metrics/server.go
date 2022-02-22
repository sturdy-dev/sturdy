package metrics

import (
	"errors"
	"fmt"
	"net/http"

	"getsturdy.com/api/pkg/configuration/flags"
	"getsturdy.com/api/pkg/di"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Configuration struct {
	Addr flags.Addr `long:"addr" description:"Address to listen on" default:"127.0.0.1:2112"`
}

type Server struct {
	srv http.Server
}

func Module(c *di.Container) {
	c.Register(New)
}

func New(cfg *Configuration) *Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return &Server{
		srv: http.Server{
			Addr:    cfg.Addr.String(),
			Handler: mux,
		},
	}
}

func (s *Server) Start() error {
	if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start metrics server: %w", err)
	}
	return nil
}
