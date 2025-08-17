package main

import (
	"log"
	"net/http"

	"github.com/fhsmendes/deploy-cloud-run/handler"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (for local development)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it, using environment variables")
	}

	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	r.Get("/temperature", handler.TemperatureHandler)

	port := "8080"
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
