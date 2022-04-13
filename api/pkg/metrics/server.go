package metrics

import (
	"errors"
	"fmt"
	"net/http"

	"getsturdy.com/api/pkg/metrics/configuration"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	srv http.Server
}

func New(cfg *configuration.Configuration) *Server {
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
