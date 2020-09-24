package infrastructure

import (
	"context"
	"net/http"

	"github.mpi-internal.com/Yapo/pro-carousel/pkg/interfaces/loggers"
)

// Server struct that implements http server to routes incoming requests
// to be proccessed. Server also includes logger to log messages in case of error
type Server struct {
	logger loggers.Logger
	server *http.Server
}

// NewHTTPServer returns a new Server suitable for use http.server and loggerHandler
// methods. NewHttpServer also includes close method to implements io.closer
func NewHTTPServer(addr string,
	routes http.Handler,
	logger loggers.Logger) *Server {
	return &Server{
		logger: logger,
		server: &http.Server{
			Addr:    addr,
			Handler: routes,
		},
	}
}

// ListenAndServe starts an HTTP server with a given address and handler.
// This method encapsulates *http.Server.ListenAndServe method and thus add
// close() method to Server struct
func (s *Server) ListenAndServe() {
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		s.logger.Crit("Error on server: %+v", err)
	}

	s.logger.Info("Closing server...")
}

// Close shuts down http.server
func (s *Server) Close() error {
	return s.server.Shutdown(context.Background())
}
