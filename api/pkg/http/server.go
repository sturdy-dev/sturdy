package http

import (
	"fmt"
	"net/http"
)

type Server struct {
	handler http.Handler
	config  *Configuration
}

func ProvideServer(cfg *Configuration, handler http.Handler) *Server {
	return &Server{
		handler: handler,
		config:  cfg,
	}
}

func (s *Server) Start() error {
	if err := http.ListenAndServe(s.config.Addr.String(), s.handler); err != http.ErrServerClosed {
		return fmt.Errorf("failed to start http server: %w", err)
	}
	return nil
}
