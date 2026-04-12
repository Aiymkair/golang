package usecase

import (
	"Assignment-7/internal/entity"
	"Assignment-7/internal/usecase/repo"
	"Assignment-7/utils"
	"errors"

	"github.com/google/uuid"
)

type UserUseCase struct {
	userRepo repo.UserRepository
}

func NewUserUseCase(userRepo repo.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (uc *UserUseCase) RegisterUser(dto *entity.CreateUserDTO) (*entity.User, error) {
	// проверяем уникальность username
	existing, _ := uc.userRepo.FindByUsername(dto.Username)
	if existing != nil {
		return nil, errors.New("username already taken")
	}

	hashedPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	role := dto.Role
	if role == "" {
		role = "user"
	}

	user := &entity.User{
		ID:       uuid.New(),
		Username: dto.Username,
		Email:    dto.Email,
		Password: hashedPassword,
		Role:     role,
	}

	created, err := uc.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (uc *UserUseCase) LoginUser(dto *entity.LoginUserDTO) (string, error) {
	user, err := uc.userRepo.FindByUsername(dto.Username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPassword(user.Password, dto.Password) {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (uc *UserUseCase) GetUserByID(userID uuid.UUID) (*entity.User, error) {
	return uc.userRepo.FindByID(userID)
}

func (uc *UserUseCase) PromoteUser(userID uuid.UUID) error {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user.Role == "admin" {
		return errors.New("user is already admin")
	}
	return uc.userRepo.UpdateRole(userID, "admin")
}
