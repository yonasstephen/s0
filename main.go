package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yonasstephen/s0/server"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./bin/")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal("failed to read config:", err)
		}
	}

	conf := server.Config{
		LogLevel:     viper.GetString("LOG_LEVEL"),
		Port:         viper.GetInt("PORT"),
		ReadTimeout:  viper.GetDuration("SERVER_WRITE_TIMEOUT"),
		WriteTimeout: viper.GetDuration("SERVER_READ_TIMEOUT"),
	}

	// set log level
	logLevel, err := logrus.ParseLevel(conf.LogLevel)
	if err != nil {
		log.Printf("invalid log level '%s', defaulting to error level", conf.LogLevel)
		log.Printf("valid log levels are: panic, fatal, error, info, warn, debug, trace")
		logrus.SetLevel(logrus.ErrorLevel)
	} else {
		logrus.SetLevel(logLevel)
	}

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	// Start http server
	server := server.NewServer(conf)
	go func() {
		server.Start()
	}()

	<-shutdownChan
	ctx, cancel := context.ContextWithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logrus.Fatalf("Server shutdown failed:", err)
	}
}
