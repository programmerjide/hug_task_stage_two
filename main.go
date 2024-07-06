package main

import (
	"hng_stage_two_task/config"
	"hng_stage_two_task/internal/api"
	"log"
)

func main() {

	cfg, err := config.SetupEnv()

	if err != nil {
		log.Fatalf("config file is not loaded properly %v\n", err)
	}

	api.StartServer(cfg)
}
