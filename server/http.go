package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/yonasstephen/s0/store/handler"
)

// Server is the server that serves API endpoints
type HTTPServer struct {
	config Config
	srv    *http.Server
}

// NewServer instantiate a new server with the given config
func NewHTTPServer(config Config) *HTTPServer {
	// setup static directory
	fs := http.FileServer(http.Dir("./static"))

	// setup router
	r := mux.NewRouter()
	r.HandleFunc("/upload", handler.UploadHandler).Methods(http.MethodPost)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("127.0.0.1:%d", config.Port),
		WriteTimeout: config.WriteTimeout,
		ReadTimeout:  config.ReadTimeout,
	}

	return &HTTPServer{
		config: config,
		srv:    srv,
	}
}

// Start sets up the server routing & start the server.
// This is a blocking function.
func (s *HTTPServer) Start() {
	logrus.Infof("Running s0 on port %d...", s.config.Port)
	if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
		logrus.Fatal(err)
	}
}

// Shutdown stops the http server gracefully
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
