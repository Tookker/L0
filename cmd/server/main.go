package main

import (
	"log"
	"os"

	"L0/internal/config"
	"L0/internal/logger"
	"L0/internal/router"
	"L0/internal/server"
	"L0/internal/store/postgre"
	"L0/internal/subscriber"
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

	db, err := postgre.NewPostgre(config, logger)
	if err != nil {
		log.Fatalln(err.Error())
	}

	broker, err := subscriber.NewBroker(config, db, logger)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = broker.Subscribe()
	if err != nil {
		log.Fatalln(err.Error())
	}

	router := router.NewChiRouter(db, logger)
	server := server.NewServer(router, db, config)

	err = server.StartServer()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
