package user

import (
	"time"

	"github.com/google/uuid"
)

type ID = uuid.UUID

type User struct {
	ID           ID
	UpdatedAt    time.Time
	Username     string
	PasswordHash []byte
	Firstname    string
	Lastname     string
}
