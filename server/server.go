package server

import (
	"context"
	"time"
)

// Config stores configuration for server
type Config struct {
	LogLevel     string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Server interface {
	Start()
	Shutdown(ctx context.Context) error
}
