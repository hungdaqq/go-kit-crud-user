package main

import (
	"context"
	"log"
	"net/http"

	"crud-gokit-postgres/internal/endpoint"
	"crud-gokit-postgres/internal/proto"
	httptransport "crud-gokit-postgres/internal/transport/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	initTracer()

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

func initTracer() {
	exporter, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint("0.0.0.0:4318"))
	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}

	// Create a new tracer provider with a batch span processor and the otlp exporter
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("hungdq30"),
		)),
	)

	// Set the global tracer provider
	otel.SetTracerProvider(tp)
}
