package usecase

import (
	"context"
	"errors"
	"fmt"

	"assignment-2/pkg/modules"
)

type userUsecase struct {
	userRepo UserRepository
}

func NewUserUsecase(repo UserRepository) UserUsecase {
	return &userUsecase{userRepo: repo}
}

func (u *userUsecase) GetUsers(ctx context.Context) ([]modules.User, error) {
	return u.userRepo.GetAll(ctx)
}

func (u *userUsecase) GetUserByID(ctx context.Context, id int) (*modules.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user with id %d not found", id)
	}
	return user, nil
}

func (u *userUsecase) CreateUser(ctx context.Context, name, email string, age *int) (int, error) {
	if name == "" {
		return 0, errors.New("user name cannot be empty")
	}
	return u.userRepo.Create(ctx, name, email, age)
}

func (u *userUsecase) UpdateUser(ctx context.Context, id int, name, email string, age *int) error {
	if name == "" {
		return errors.New("user name cannot be empty")
	}
	rows, err := u.userRepo.Update(ctx, id, name, email, age)
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("user with id %d does not exist", id)
	}
	return nil
}

func (u *userUsecase) DeleteUser(ctx context.Context, id int) error {
	rows, err := u.userRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("user with id %d does not exist", id)
	}
	return nil
}
