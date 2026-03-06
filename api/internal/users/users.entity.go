package user

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
}
