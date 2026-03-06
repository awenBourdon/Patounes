package dto

import "github.com/google/uuid"

type UserResponseDto struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}
