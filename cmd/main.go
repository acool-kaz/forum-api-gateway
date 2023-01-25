package main

import (
	"log"

	"github.com/acool-kaz/forum-api-gateway/internal/app"
	"github.com/acool-kaz/forum-api-gateway/internal/config"
)

func main() {
	cfg, err := config.InitConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	app, err := app.InitApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app.Run()
}
