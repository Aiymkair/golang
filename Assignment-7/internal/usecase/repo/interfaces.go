package repo

import (
	"Assignment-7/internal/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(user *entity.User) (*entity.User, error)
	FindByUsername(username string) (*entity.User, error)
	FindByID(id uuid.UUID) (*entity.User, error)
	UpdateRole(id uuid.UUID, newRole string) error
}
