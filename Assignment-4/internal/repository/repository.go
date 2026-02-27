package repository

import (
	"assignment-4/internal/repository/_postgres"
	"assignment-4/internal/repository/_postgres/users"
)

type Repositories struct {
	Users *users.UserRepository
}

func NewRepositories(pg *_postgres.Dialect) *Repositories {
	return &Repositories{
		Users: users.NewUserRepository(pg),
	}
}
