package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/bmon/voting-website/pkg/api"
)

func main() {
	// load dotenv file into process environ
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Instantiate the API object and serve the application
	a := api.New()

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", a.Config.Port),
		Handler: a.Router(),
	}

	closeSignal := make(chan os.Signal, 1)
	signal.Notify(closeSignal, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-closeSignal
		log.Println("Recieved ", sig, ", shutting down")

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
	}()

	log.Println("Starting HTTP Server. Listening on", srv.Addr)
	log.Println(srv.ListenAndServe())
}
