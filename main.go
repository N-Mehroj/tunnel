package main

import (
    "fmt"
    "net/http"
    "go-tunnel/routes"
    "go-tunnel/services"

)

func main() {
    services.ConnectDatabase()
    routes.SetupRoutes()
    // fmt.Println("Server 8080 portda ishlayapti")
    fmt.Println("Server http://localhost:8080 portda ishlayapti")
    http.ListenAndServe(":8080", nil)
}
