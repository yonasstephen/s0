package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/yonasstephen/s0/store/handler"
)

// Config stores configuration for server
type Config struct {
	LogLevel     string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Server is the server that serves API endpoints
type Server struct {
	config Config
}

// NewServer instantiate a new server with the given config
func NewServer(config Config) *Server {
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

	log.Printf("Running s0 on port %d...", s.config.Port)
	log.Fatal(srv.ListenAndServe())
}
