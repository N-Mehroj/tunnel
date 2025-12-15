package main

import (
	"log"
	"net/http"
	"os"

	"go-tunnel/routes"
	"go-tunnel/services"
)

func main() {
	services.ConnectDatabase()
	
	// Setup routes (Laravel-style)
	router := routes.SetupRoutes()

	host := os.Getenv("APP_URL")
	if host == "" {
		host = "0.0.0.0:8080"
	}

	log.Printf("Server running at http://%s", host)
	log.Printf("Available routes: %v", len(router.ListRoutes()))

	// Use the router as HTTP handler
	log.Fatal(http.ListenAndServe(host, router))
}
