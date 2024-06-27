package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "grpc-server/proto" // Import generated protobuf package

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

const (
	port           = ":50051"
	dataSourceName = "host=localhost port=5432 user=postgres dbname=postgres password=1 sslmode=disable"
	driverName     = "postgres"
)

// User represents a user model
type User struct {
	Id       int64  `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

var db *sqlx.DB

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) CreateUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	user := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id"
	var id int64
	err := db.QueryRowContext(ctx, query, user.Name, user.Email, user.Password).Scan(&id)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Insert user with ID: %d\n", id)
	user.Id = id
	return &pb.UserResponse{User: &pb.User{
		Id:       user.Id,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}}, nil
}

func (s *server) GetUser(ctx context.Context, req *pb.UserID) (*pb.UserResponse, error) {
	var user User
	err := db.GetContext(ctx, &user, "SELECT id, name, email, password FROM users WHERE id=$1", req.Id)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Get user with ID: %d\n", user.Id)
	return &pb.UserResponse{User: &pb.User{
		Id:       user.Id,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}}, nil
}

func (s *server) UpdateUser(ctx context.Context, req *pb.User) (*pb.UserResponse, error) {
	user := &User{
		Id:       req.Id,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	_, err := db.ExecContext(ctx, "UPDATE users SET name=$1, email=$2, password=$3 WHERE id=$4",
		user.Name, user.Email, user.Password, user.Id)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Update user with ID: %d\n", user.Id)
	return &pb.UserResponse{User: &pb.User{
		Id:       user.Id,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *pb.UserID) (*pb.UserResponse, error) {

	_, err := db.ExecContext(ctx, "DELETE FROM users WHERE id=$1", req.Id)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Delete user with ID: %d\n", req.Id)
	return &pb.UserResponse{}, nil
}

func main() {
	// Connect to PostgreSQL database
	var err error
	db, err = sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create gRPC server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})

	log.Printf("gRPC server listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
