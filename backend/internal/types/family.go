package types

import "time"

type Family struct {
	ID          string
	Name        string
	Description string
	UpdatedAt   time.Time
}
