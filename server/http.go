package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/yonasstephen/s0/store/handler"
)

// Server is the server that serves API endpoints
type HTTPServer struct {
	config Config
}

// NewServer instantiate a new server with the given config
func NewHTTPServer(config Config) *HTTPServer {
	return &Server{
		config: config,
	}
}

// Start sets up the server routing & start the server.
// will throw panic if the server fails to start
func (s *Server) Start() {
	// setup static directory
	fs := http.FileServer(http.Dir("./static"))

	// setup router
	r := mux.NewRouter()
	r.HandleFunc("/upload", handler.UploadHandler).Methods(http.MethodPost)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("127.0.0.1:%d", s.config.Port),
		WriteTimeout: s.config.WriteTimeout,
		ReadTimeout:  s.config.ReadTimeout,
	}

	logLevel, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		log.Printf("invalid log level '%s', defaulting to error level", s.config.LogLevel)
		log.Printf("valid log levels are: panic, fatal, error, info, warn, debug, trace")
		logrus.SetLevel(logrus.ErrorLevel)
	} else {
		logrus.SetLevel(logLevel)
	}

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		logrus.Infof("Running s0 on port %d...", s.config.Port)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logrus.Fatal(err)
		}
	}()

	logrus.Info("Server started")

	<-shutdownChan
	logrus.Info("Server stopped")

	ctx, cancel := context.ContextWithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown failed:", err)
	}
	logrus.Info("Server stopped gracefully")
}
