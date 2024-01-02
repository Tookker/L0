package main

import (
	"L0/internal/config"
	"log"
	"os"
)

func main() {
	config, err := config.LoadConfig(os.Getenv("CONF_PATH"))
	if err != nil {
		log.Fatalln(err.Error())
	}
}
