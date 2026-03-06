package auth

import (
	"context"
)

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*AuthResponse, error)
	Login(ctx context.Context, input LoginInput) (*AuthResponse, error)
}
