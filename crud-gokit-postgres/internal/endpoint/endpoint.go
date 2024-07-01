package endpoint

import (
	"context"

	"crud-gokit-postgres/internal/middleware"
	"crud-gokit-postgres/internal/model"
	"crud-gokit-postgres/internal/proto"

	"github.com/go-kit/kit/endpoint"
	"google.golang.org/grpc"
)

func NewUserServiceClient(conn *grpc.ClientConn) proto.UserServiceClient {
	return proto.NewUserServiceClient(conn)
}

type Endpoints struct {
	CreateUserEndpoint endpoint.Endpoint
	GetUserEndpoint    endpoint.Endpoint
	UpdateUserEndpoint endpoint.Endpoint
	DeleteUserEndpoint endpoint.Endpoint
}

func MakeEndpoints(client proto.UserServiceClient, authUser, authPassword, authRealm string) Endpoints {
	authMiddleware := middleware.AuthMiddleware(authUser, authPassword, authRealm)
	createUserEndpoint := makeCreateUserEndpoint(client)
	getUserEndpoint := makeGetUserEndpoint(client)
	updateUserEndpoint := makeUpdateUserEndpoint(client)
	deleteUserEndpoint := makeDeleteUserEndpoint(client)

	// Apply authentication middleware to each endpoint
	return Endpoints{
		CreateUserEndpoint: authMiddleware(createUserEndpoint),
		GetUserEndpoint:    authMiddleware(getUserEndpoint),
		UpdateUserEndpoint: authMiddleware(updateUserEndpoint),
		DeleteUserEndpoint: authMiddleware(deleteUserEndpoint),
	}
}

func makeCreateUserEndpoint(client proto.UserServiceClient) endpoint.Endpoint {

	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(CreateUserRequest)
		grpcReq := &proto.UserRequest{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
		}
		grpcResp, err := client.CreateUser(ctx, grpcReq)
		if err != nil {
			return nil, err
		}
		return CreateUserResponse{Id: grpcResp.User.Id}, nil
	}
}

func makeGetUserEndpoint(client proto.UserServiceClient) endpoint.Endpoint {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserRequest)
		grpcReq := &proto.UserID{
			Id: req.Id,
		}
		grpcResp, err := client.GetUser(ctx, grpcReq)
		if err != nil {
			return nil, err
		}
		user := model.User{
			Id:       grpcResp.User.Id,
			Name:     grpcResp.User.Name,
			Email:    grpcResp.User.Email,
			Password: grpcResp.User.Password,
		}
		return GetUserResponse{User: user}, nil
	}
}

func makeUpdateUserEndpoint(client proto.UserServiceClient) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateUserRequest)
		grpcReq := &proto.User{
			Id:       req.Id,
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
		}
		_, err := client.UpdateUser(ctx, grpcReq)
		if err != nil {
			return nil, err
		}
		return UpdateUserResponse{Success: true}, nil
	}
}

func makeDeleteUserEndpoint(client proto.UserServiceClient) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteUserRequest)
		grpcReq := &proto.UserID{
			Id: req.Id,
		}
		_, err := client.DeleteUser(ctx, grpcReq)
		if err != nil {
			return nil, err
		}
		return DeleteUserResponse{Success: true}, nil
	}
}

// Request and Response structs
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	Id int64 `json:"id"`
}

type GetUserRequest struct {
	Id int64 `json:"id"`
}

type GetUserResponse struct {
	User model.User `json:"user"`
}

type UpdateUserRequest struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserResponse struct {
	Success bool `json:"success"`
}

type DeleteUserRequest struct {
	Id int64 `json:"id"`
}

type DeleteUserResponse struct {
	Success bool `json:"success"`
}
