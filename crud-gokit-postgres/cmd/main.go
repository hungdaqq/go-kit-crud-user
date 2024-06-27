package main

import (
	"log"
	"net/http"

	"crud-gokit-postgres/internal/endpoint"
	"crud-gokit-postgres/internal/proto"
	httptransport "crud-gokit-postgres/internal/transport/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Authentication details
	authUser := "IOT"
	authPassword := "1"
	authRealm := "ProtectedArea"

	// gRPC connection setup
	grpcAddr := "localhost:50051" // Address of the gRPC UserService server
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// Create gRPC client
	userServiceClient := proto.NewUserServiceClient(conn)

	// Create endpoints with authentication middleware
	endpoints := endpoint.MakeEndpoints(userServiceClient, authUser, authPassword, authRealm)

	// Create HTTP handler
	httpHandler := httptransport.NewHTTPHandler(endpoints)

	// Define server parameters
	addr := ":8080" // specify the port you want to listen on

	// Start HTTP server
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, httpHandler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
