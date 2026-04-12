package usecase

import (
	"Assignment-7/internal/entity"
	"github.com/google/uuid"
)

type UserInterface interface {
	RegisterUser(dto *entity.CreateUserDTO) (*entity.User, error)
	LoginUser(dto *entity.LoginUserDTO) (string, error)
	GetUserByID(userID uuid.UUID) (*entity.User, error)
	PromoteUser(userID uuid.UUID) error
}
