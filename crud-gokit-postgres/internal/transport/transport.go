package httptransport

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "github.com/go-kit/kit/endpoint"
    myEndpoint "crud-gokit-postgres/internal/endpoint"
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
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        request, err := decodeRequest(r)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        response, err := e(ctx, request)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
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
    var req CreateUserRequest
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
    return GetUserRequest{ID: int64(id)}, nil
}

func decodeUpdateUserRequest(r *http.Request) (interface{}, error) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        return nil, err
    }
    var req UpdateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        return nil, err
    }
    req.ID = int64(id)
    return req, nil
}

func decodeDeleteUserRequest(r *http.Request) (interface{}, error) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        return nil, err
    }
    return DeleteUserRequest{ID: int64(id)}, nil
}

type CreateUserRequest = myEndpoint.CreateUserRequest
type GetUserRequest = myEndpoint.GetUserRequest
type UpdateUserRequest = myEndpoint.UpdateUserRequest
type DeleteUserRequest = myEndpoint.DeleteUserRequest