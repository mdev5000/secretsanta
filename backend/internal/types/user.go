package types

import (
	"time"
)

type ID = string

type User struct {
	ID           ID
	UpdatedAt    time.Time
	Username     string
	PasswordHash []byte
	Firstname    string
	Lastname     string
	FamilyIDs    []string
}
