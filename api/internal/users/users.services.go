package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
	if id == uuid.Nil {
		return nil, errors.New("user id is required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context) ([]*User, error) {
	return s.repo.List(ctx)
}

func (s *userService) CreateUser(ctx context.Context, name, email string) (*User, error) {
	// 🎯 Validation métier
	if name == "" {
		return nil, errors.New("name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}

	user := &User{
		Name:  name,
		Email: email,
	}

	return s.repo.Create(ctx, user)
}

func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, name, email string) (*User, error) {
	if id == uuid.Nil {
		return nil, errors.New("user id is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}

	user := &User{
		ID:    id,
		Name:  name,
		Email: email,
	}

	return s.repo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("user id is required")
	}
	return s.repo.Delete(ctx, id)
}
