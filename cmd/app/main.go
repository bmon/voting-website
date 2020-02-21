package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/bmon/voting-website/pkg/api"
	"github.com/bmon/voting-website/pkg/env"
)

func main() {
	// load dotenv file into process environ
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	config := env.LoadConfig()
	// create and load config from environ
	api.RegisterRoutes(config)
	// serve the application and log requests as they arrive
	log.Println("Starting server on :" + config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
