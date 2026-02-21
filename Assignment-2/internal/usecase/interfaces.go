package usecase

import (
	"assignment-2/pkg/modules"
	"context"
)

// UserRepository
type UserRepository interface {
	GetAll(ctx context.Context) ([]modules.User, error)
	GetByID(ctx context.Context, id int) (*modules.User, error)
	Create(ctx context.Context, name, email string, age *int) (int, error)
	Update(ctx context.Context, id int, name, email string, age *int) (int64, error)
	Delete(ctx context.Context, id int) (int64, error)
}

// UserUsecase
type UserUsecase interface {
	GetUsers(ctx context.Context) ([]modules.User, error)
	GetUserByID(ctx context.Context, id int) (*modules.User, error)
	CreateUser(ctx context.Context, name, email string, age *int) (int, error)
	UpdateUser(ctx context.Context, id int, name, email string, age *int) error
	DeleteUser(ctx context.Context, id int) error
}
