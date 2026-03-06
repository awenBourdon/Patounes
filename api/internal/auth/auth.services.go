package auth

import (
	"context"
	"errors"
	"patoune-api/database/sqlc"

	"github.com/google/uuid"
)

type authService struct {
	queries *sqlc.Queries
}

func NewAuthService(queries *sqlc.Queries) AuthService {
	return &authService{
		queries: queries,
	}
}

func (s *authService) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
	if input.Name == "" {
		return nil, errors.New("Name is required.")
	}
	if input.Email == "" {
		return nil, errors.New("Email is required.")
	}
	if input.Password == "" {
		return nil, errors.New("Password is required.")
	}

	if err := ValidatePassword(input.Password); err != nil {
		return nil, err
	}

	existingUser, err := s.queries.GetUserByEmail(ctx, input.Email)
	if err == nil && existingUser.ID != uuid.Nil {
		return nil, errors.New("This email is already used.")
	}

	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		return nil, errors.New("A error appear during the hash of the password.")
	}

	newUser, err := s.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	})
	if err != nil {
		return nil, errors.New("The creation failed.")
	}

	token, err := GenerateToken(newUser.ID, newUser.Email)
	if err != nil {
		return nil, errors.New("The token generation failed.")
	}

	return &AuthResponse{
		Token: token,
		User: UserResponse{
			ID:    newUser.ID.String(),
			Name:  newUser.Name,
			Email: newUser.Email,
		},
	}, nil
}

func (s *authService) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {

	if input.Email == "" {
		return nil, errors.New("Email is required.")
	}
	if input.Password == "" {
		return nil, errors.New("Password i required.")
	}

	user, err := s.queries.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.New("Email and/or password failed.")
	}

	match, err := VerifyPassword(input.Password, user.Password)
	if err != nil || !match {
		return nil, errors.New("Email and/or password failed.")
	}

	token, err := GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("Token creation failed.")
	}

	return &AuthResponse{
		Token: token,
		User: UserResponse{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}
