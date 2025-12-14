package routes

import (
	"net/http"

	"go-tunnel/handlers"
)

func SetupRoutes() {
    http.HandleFunc("/login", handlers.LoginHandler)
}
