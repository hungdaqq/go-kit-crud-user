package service

import (
    "context"
    "fmt"

    "github.com/jmoiron/sqlx"
    "crud-gokit-postgres/internal/model"
)

type User = model.User
// UserService describes the service.
type UserService interface {
    Create(ctx context.Context, user User) error
    Get(ctx context.Context, id int64) (User, error)
    Update(ctx context.Context, id int64, user User) error
    Delete(ctx context.Context, id int64) error
}

// // User represents a user entity.
// type User struct {
//     ID       int64  `json:"id" db:"id"`
//     Name     string `json:"name" db:"name"`
//     Email    string `json:"email" db:"email"`
//     Password string `json:"password" db:"password"`
// }

// userService implements UserService interface.
type userService struct {
    db *sqlx.DB
}

// NewUserService creates a new instance of userService.
func NewUserService(db *sqlx.DB) UserService {
    return &userService{db: db}
}

func (s *userService) Create(ctx context.Context, user User) (error) {
    query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id"
    row := s.db.QueryRowContext(ctx, query, user.Name, user.Email, user.Password)
    err := row.Scan(&user.ID)
    if err != nil {
        return fmt.Errorf("error creating user: %v", err)
    }
    return nil
}

func (s *userService) Get(ctx context.Context, id int64) (User, error) {
    var user User
    query := "SELECT id, name, email, password FROM users WHERE id = $1"
    err := s.db.GetContext(ctx, &user, query, id)
    if err != nil {
        return User{}, fmt.Errorf("error getting user: %v", err)
    }
    return user, nil
}

func (s *userService) Update(ctx context.Context, id int64, user User) error {
    query := "UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4"
    _, err := s.db.ExecContext(ctx, query, user.Name, user.Email, user.Password, id)
    if err != nil {
        return fmt.Errorf("error updating user: %v", err)
    }
    return nil
}

func (s *userService) Delete(ctx context.Context, id int64) error {
    query := "DELETE FROM users WHERE id = $1"
    _, err := s.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("error deleting user: %v", err)
    }
    return nil
}
