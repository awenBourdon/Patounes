package user

import (
	"context"
	"patoune-api/database/sqlc"

	"github.com/google/uuid"
)

type userRepository struct {
	queries *sqlc.Queries
}

func NewUserRepository(queries *sqlc.Queries) UserRepository {
	return &userRepository{queries: queries}
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	sqlcUser, err := r.queries.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        sqlcUser.ID,
		Name:      sqlcUser.Name,
		Email:     sqlcUser.Email,
		CreatedAt: sqlcUser.CreatedAt,
		UpdatedAt: sqlcUser.UpdatedAt,
		DeletedAt: sqlcUser.DeletedAt,
	}, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	sqlcUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        sqlcUser.ID,
		Name:      sqlcUser.Name,
		Email:     sqlcUser.Email,
		CreatedAt: sqlcUser.CreatedAt,
		UpdatedAt: sqlcUser.UpdatedAt,
		DeletedAt: sqlcUser.DeletedAt,
	}, nil
}

func (r *userRepository) List(ctx context.Context) ([]*User, error) {
	sqlcUsers, err := r.queries.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]*User, len(sqlcUsers))
	for i, sqlcUser := range sqlcUsers {
		users[i] = &User{
			ID:        sqlcUser.ID,
			Name:      sqlcUser.Name,
			Email:     sqlcUser.Email,
			CreatedAt: sqlcUser.CreatedAt,
			UpdatedAt: sqlcUser.UpdatedAt,
			DeletedAt: sqlcUser.DeletedAt,
		}
	}

	return users, nil
}

func (r *userRepository) Create(ctx context.Context, user *User) (*User, error) {
	sqlcUser, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        sqlcUser.ID,
		Name:      sqlcUser.Name,
		Email:     sqlcUser.Email,
		Password:  sqlcUser.Password,
		CreatedAt: sqlcUser.CreatedAt,
		UpdatedAt: sqlcUser.UpdatedAt,
		DeletedAt: sqlcUser.DeletedAt,
	}, nil
}

func (r *userRepository) Update(ctx context.Context, user *User) (*User, error) {
	sqlcUser, err := r.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        sqlcUser.ID,
		Name:      sqlcUser.Name,
		Email:     sqlcUser.Email,
		CreatedAt: sqlcUser.CreatedAt,
		UpdatedAt: sqlcUser.UpdatedAt,
		DeletedAt: sqlcUser.DeletedAt,
	}, nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteUser(ctx, id)
}
