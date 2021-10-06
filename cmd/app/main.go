package main

import (
	"github.com/Ivanhahanov/ad-infrastructure-api/config"
	"github.com/Ivanhahanov/ad-infrastructure-api/internal/app"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

func main() {
	var cfg config.Config

	err := cleanenv.ReadConfig("../../config/config.yml", &cfg)
	if err != nil {log.Fatalf("Config error: %s", err)}

	app.Run(&cfg)
}