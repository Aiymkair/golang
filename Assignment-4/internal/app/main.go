package app

import (
	"context"
	"fmt"
	"time"

	"assignment-4/internal/repository"
	"assignment-4/internal/repository/_postgres"
	"assignment-4/pkg/modules"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConfig := initPostgreSQL()
	pgDialect := _postgres.NewPGXDialect(ctx, dbConfig)

	repositories := repository.NewRepositories(pgDialect)

	users, err := repositories.Users.GetAll(ctx)
	if err != nil {
		fmt.Printf("Error fetching users: %v\n", err)
		return
	}
	fmt.Printf("Users: %+v\n", users)
}

func initPostgreSQL() *modules.PostgreSQL {
	return &modules.PostgreSQL{
		Host:        "localhost",
		Port:        "5432",
		Username:    "postgres",
		Password:    "postgres",
		DBName:      "mydb",
		SSLMode:     "disable",
		ExecTimeout: 5 * time.Second,
	}
}
