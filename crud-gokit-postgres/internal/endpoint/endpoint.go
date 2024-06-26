package endpoint

import (
    "context"

    "github.com/go-kit/kit/endpoint"
    "crud-gokit-postgres/internal/service"
    "crud-gokit-postgres/internal/model"
)

type Endpoints struct {
    CreateUserEndpoint  endpoint.Endpoint
    GetUserEndpoint     endpoint.Endpoint
    UpdateUserEndpoint  endpoint.Endpoint
    DeleteUserEndpoint  endpoint.Endpoint
}

// type Endpoint interface {
//     ServeHTTP(context.Context, interface{}) (interface{}, error)
// }

func MakeEndpoints(s service.UserService) Endpoints {
    return Endpoints{
        CreateUserEndpoint:  makeCreateUserEndpoint(s),
        GetUserEndpoint:     makeGetUserEndpoint(s),
        UpdateUserEndpoint:  makeUpdateUserEndpoint(s),
        DeleteUserEndpoint:  makeDeleteUserEndpoint(s),
    }
}

func makeCreateUserEndpoint(s service.UserService) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(CreateUserRequest)
        user := model.User{
            Name:     req.Name,
            Email:    req.Email,
            Password: req.Password,
        }

        err := s.Create(ctx, user)
        if err != nil {
            return nil, err
        }
        return CreateUserResponse{ID: user.ID}, nil
    }
}

func makeGetUserEndpoint(s service.UserService) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(GetUserRequest)
        user, err := s.Get(ctx, req.ID)
        if err != nil {
            return nil, err
        }
        return GetUserResponse{User: user}, nil
    }
}

func makeUpdateUserEndpoint(s service.UserService) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(UpdateUserRequest)
        userUpdate := model.User{
            ID:       req.ID,
            Name:     req.Name,
            Email:    req.Email,
            Password: req.Password,
        }
        err := s.Update(ctx, req.ID, userUpdate)
        if err != nil {
            return nil, err
        }
        return UpdateUserResponse{Success: true}, nil
    }
}

func makeDeleteUserEndpoint(s service.UserService) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(DeleteUserRequest)
        err := s.Delete(ctx, req.ID)
        if err != nil {
            return nil, err
        }
        return DeleteUserResponse{Success: true}, nil
    }
}

// Request and Response structs
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Password string `json:"password"`
}

type CreateUserResponse struct {
    ID int64 `json:"id"`
}

type GetUserRequest struct {
    ID int64 `json:"id"`
}

type GetUserResponse struct {
    User model.User `json:"user"`
}

type UpdateUserRequest struct {
    ID    int64    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Password string `json:"password"`
}

type UpdateUserResponse struct {
    Success bool `json:"success"`
}

type DeleteUserRequest struct {
    ID int64 `json:"id"`
}

type DeleteUserResponse struct {
    Success bool `json:"success"`
}
