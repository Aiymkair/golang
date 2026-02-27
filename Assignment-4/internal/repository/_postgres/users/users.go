package users

import (
	"context"
	"database/sql"
	"fmt"

	"assignment-4/internal/repository/_postgres"
	"assignment-4/pkg/modules"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(pg *_postgres.Dialect) *UserRepository {
	return &UserRepository{db: pg.DB}
}

// GetAll
func (r *UserRepository) GetAll(ctx context.Context) ([]modules.User, error) {
	var users []modules.User
	query := `SELECT id, name, email, age, created_at FROM users`
	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, fmt.Errorf("UserRepository.GetAll: %w", err)
	}
	return users, nil
}

// GetByID
func (r *UserRepository) GetByID(ctx context.Context, id int) (*modules.User, error) {
	var user modules.User
	query := `SELECT id, name, email, age, created_at FROM users WHERE id = $1`
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("UserRepository.GetByID: %w", err)
	}
	return &user, nil
}

// Create
func (r *UserRepository) Create(ctx context.Context, name, email string, age *int) (int, error) {
	var id int
	query := `INSERT INTO users (name, email, age) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRowxContext(ctx, query, name, email, age).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("UserRepository.Create: %w", err)
	}
	return id, nil
}

// Update
func (r *UserRepository) Update(ctx context.Context, id int, name, email string, age *int) (int64, error) {
	query := `UPDATE users SET name = $1, email = $2, age = $3 WHERE id = $4`
	res, err := r.db.ExecContext(ctx, query, name, email, age, id)
	if err != nil {
		return 0, fmt.Errorf("UserRepository.Update: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("UserRepository.Update: rows affected: %w", err)
	}
	return rows, nil
}

// Delete
func (r *UserRepository) Delete(ctx context.Context, id int) (int64, error) {
	query := `DELETE FROM users WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return 0, fmt.Errorf("UserRepository.Delete: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("UserRepository.Delete: rows affected: %w", err)
	}
	return rows, nil
}
