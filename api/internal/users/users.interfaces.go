package user

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context) ([]*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserService interface {
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)
	ListUsers(ctx context.Context) ([]*User, error)
	CreateUser(ctx context.Context, name, email string) (*User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, name, email string) (*User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
