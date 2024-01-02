package main

import (
	"log"
	"os"

	"L0/internal/caches/ordercache/ordercachemap"
	"L0/internal/config"
	"L0/internal/logger"
)

func main() {
	config, err := config.LoadConfig(os.Getenv("CONF_PATH"))
	if err != nil {
		log.Fatalln(err.Error())
	}

	logger, err := logger.NewLogger(config)
	if err != nil {
		log.Fatalln(err.Error())
	}

	orderCache := ordercachemap.NewCache(logger)
}
