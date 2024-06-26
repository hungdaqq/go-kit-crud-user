package main

import (
    "log"
    "net/http"

    _ "github.com/lib/pq"

    "crud-gokit-postgres/internal/db"
    "crud-gokit-postgres/internal/endpoint"
    "crud-gokit-postgres/internal/service"
    httptransport "crud-gokit-postgres/internal/transport"
)

func main() {
    // Initialize database connection
    db.InitDB("postgres://postgres:1@localhost/postgres?sslmode=disable")
    defer db.GetDB().Close()

    // Create service
    userService := service.NewUserService(db.GetDB())

    // Create endpoints
    endpoints := endpoint.MakeEndpoints(userService)

    // Create HTTP handler
    httpHandler := httptransport.NewHTTPHandler(endpoints)

    // Define server parameters
    addr := "0.0.0.0:8088" // specify the port you want to listen on

    // Start HTTP server
    log.Printf("Starting server on %s", addr)
    err := http.ListenAndServe(addr, httpHandler)
    if err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
