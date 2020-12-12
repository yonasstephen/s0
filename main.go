package main

import (
	"log"

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

	server := server.NewServer(conf)
	server.Start()
}
