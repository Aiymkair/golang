package _postgres

import (
	"context"
	"fmt"
	"log"

	"assignment-2/pkg/modules"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Dialect struct {
	DB *sqlx.DB
}

func NewPGXDialect(ctx context.Context, cfg *modules.PostgreSQL) *Dialect {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to db: %v", err))
	}

	if err = db.Ping(); err != nil {
		panic(fmt.Sprintf("db ping failed: %v", err))
	}

	log.Println("database connection established")
	return &Dialect{DB: db}
}

// AutoMigrate – опционально, если нужно применять миграции
func AutoMigrate(cfg *modules.PostgreSQL) {
	sourceURL := "file://database/migrations"
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		panic(fmt.Sprintf("migrate init failed: %v", err))
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(fmt.Sprintf("migration up failed: %v", err))
	}

	log.Println("database migrations applied successfully")
}
