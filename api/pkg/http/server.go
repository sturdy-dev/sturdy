package http

import (
	"fmt"
	"net/http"
)

type Server struct {
	handler http.Handler
}

func ProvideServer(handler http.Handler) *Server {
	return &Server{
		handler: handler,
	}
}
func (s *Server) ListenAndServe(addr string) error {
	if err := http.ListenAndServe(addr, s.handler); err != http.ErrServerClosed {
		return fmt.Errorf("failed to start http server: %w", err)
	}
	return nil
}
