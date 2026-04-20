package domain

import "time"

type Repository struct {
	FullName    string
	Description string
	Stars       int32
	Forks       int64
	CreatedAt   time.Time
	Language    string
}
