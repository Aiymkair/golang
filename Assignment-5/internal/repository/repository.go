package repository

import (
	"Assignment-5/internal/models"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	_ "time"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) DB() *sql.DB {
	return r.db
}

var allowedColumns = map[string]bool{
	"id":         true,
	"name":       true,
	"email":      true,
	"gender":     true,
	"birth_date": true,
}

// GetPaginatedUsers
func (r *Repository) GetPaginatedUsers(page, pageSize int, filters map[string]interface{}, orderBy string) (models.PaginatedResponse, error) {
	offset := (page - 1) * pageSize

	selectQuery := `SELECT id, name, email, gender, birth_date, deleted_at FROM users`
	countQuery := `SELECT COUNT(*) FROM users`

	whereClause, args := buildWhereClause(filters)

	fullSelect := selectQuery
	fullCount := countQuery
	if whereClause != "" {
		fullSelect += " WHERE " + whereClause
		fullCount += " WHERE " + whereClause
	}

	// сортировка
	if orderBy != "" && allowedColumns[orderBy] {
		fullSelect += " ORDER BY " + orderBy
	} else {
		fullSelect += " ORDER BY id" // по умолчанию
	}

	// LIMIT и OFFSET
	fullSelect += " LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, pageSize, offset)

	// count
	var total int
	var countArgs []interface{}
	if len(args) > 2 {
		countArgs = args[:len(args)-2]
	}
	err := r.db.QueryRow(fullCount, countArgs...).Scan(&total)
	if err != nil {
		return models.PaginatedResponse{}, fmt.Errorf("count query failed: %w", err)
	}

	// select
	rows, err := r.db.Query(fullSelect, args...)
	if err != nil {
		return models.PaginatedResponse{}, fmt.Errorf("select query failed: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		var deletedAt sql.NullTime
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate, &deletedAt)
		if err != nil {
			return models.PaginatedResponse{}, err
		}
		if deletedAt.Valid {
			u.DeletedAt = &deletedAt.Time
		}
		users = append(users, u)
	}

	return models.PaginatedResponse{
		Data:       users,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// buildWhereClause
func buildWhereClause(filters map[string]interface{}) (string, []interface{}) {
	var clauses []string
	var args []interface{}
	idx := 1

	if val, ok := filters["deleted"]; ok {
		delete(filters, "deleted")
		showDeleted := false
		switch v := val.(type) {
		case bool:
			showDeleted = v
		case string:
			showDeleted = v == "true" || v == "1"
		}
		if showDeleted {
			clauses = append(clauses, "deleted_at IS NOT NULL")
		} else {
			clauses = append(clauses, "deleted_at IS NULL")
		}
	} else {
		// по умолчанию активные
		clauses = append(clauses, "deleted_at IS NULL")
	}

	for key, val := range filters {
		if !allowedColumns[key] {
			continue
		}
		switch key {
		case "name", "email":
			clauses = append(clauses, key+" ILIKE $"+strconv.Itoa(idx))
			args = append(args, "%"+val.(string)+"%")
		default:
			clauses = append(clauses, key+" = $"+strconv.Itoa(idx))
			args = append(args, val)
		}
		idx++
	}

	return strings.Join(clauses, " AND "), args
}

// GetCommonFriends
func (r *Repository) GetCommonFriends(userID1, userID2 int) ([]models.User, error) {
	query := `
        WITH friends_of_a AS (
            SELECT friend_id FROM user_friends WHERE user_id = $1
            UNION
            SELECT user_id FROM user_friends WHERE friend_id = $1
        ), friends_of_b AS (
            SELECT friend_id FROM user_friends WHERE user_id = $2
            UNION
            SELECT user_id FROM user_friends WHERE friend_id = $2
        ), common_friends AS (
            SELECT friend_id FROM friends_of_a
            INTERSECT
            SELECT friend_id FROM friends_of_b
        )
        SELECT u.id, u.name, u.email, u.gender, u.birth_date, u.deleted_at
        FROM users u
        JOIN common_friends cf ON u.id = cf.friend_id
        WHERE u.deleted_at IS NULL
    `

	rows, err := r.db.Query(query, userID1, userID2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []models.User
	for rows.Next() {
		var u models.User
		var deletedAt sql.NullTime
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate, &deletedAt); err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			u.DeletedAt = &deletedAt.Time
		}
		friends = append(friends, u)
	}
	return friends, nil
}
