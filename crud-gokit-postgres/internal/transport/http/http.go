package httptransport

import (
	"context"
	myEndpoint "crud-gokit-postgres/internal/endpoint"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// NewHTTPHandler creates a new HTTP handler for the endpoints.
func NewHTTPHandler(endpoints myEndpoint.Endpoints) http.Handler {
	r := mux.NewRouter()

	r.Methods("POST").Path("/api/users").Handler(httpHandler(endpoints.CreateUserEndpoint, decodeCreateUserRequest))
	r.Methods("GET").Path("/api/users/{id}").Handler(httpHandler(endpoints.GetUserEndpoint, decodeGetUserRequest))
	r.Methods("PUT").Path("/api/users/{id}").Handler(httpHandler(endpoints.UpdateUserEndpoint, decodeUpdateUserRequest))
	r.Methods("DELETE").Path("/api/users/{id}").Handler(httpHandler(endpoints.DeleteUserEndpoint, decodeDeleteUserRequest))

	return r
}

// httpHandler wraps a Go-Kit endpoint with the necessary HTTP decoding and encoding.
func httpHandler(e endpoint.Endpoint, decodeRequest func(r *http.Request) (interface{}, error)) http.Handler {
	tracer := otel.Tracer("crud-gokit-postgres")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization token missing", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "basic" {
			http.Error(w, "Invalid Authorization token format", http.StatusUnauthorized)
			return
		}

		authToken := parts[1]

		// Create a new context with the token
		ctx := context.WithValue(r.Context(), httptransport.ContextKeyRequestAuthorization, authToken)
		ctx, span := tracer.Start(ctx, r.URL.Path)
		defer span.End()

		request, err := decodeRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := e(ctx, request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			// Record the status in the span
			span.SetStatus(codes.Error, err.Error())
			span.SetAttributes(attribute.Int("http.status_code", http.StatusInternalServerError))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

// Decode functions for each request type.

func decodeCreateUserRequest(r *http.Request) (interface{}, error) {
	var req myEndpoint.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeGetUserRequest(r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return nil, err
	}
	return myEndpoint.GetUserRequest{Id: int64(id)}, nil
}

func decodeUpdateUserRequest(r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return nil, err
	}
	var req myEndpoint.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	req.Id = int64(id)
	return req, nil
}

func decodeDeleteUserRequest(r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return nil, err
	}
	return myEndpoint.DeleteUserRequest{Id: int64(id)}, nil
}
