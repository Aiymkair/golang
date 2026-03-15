package models

import (
	"time"
)

type User struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Gender    string     `json:"gender"`
	BirthDate time.Time  `json:"birth_date"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type PaginatedResponse struct {
	Data       []User `json:"data"`
	TotalCount int    `json:"totalCount"`
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
}
