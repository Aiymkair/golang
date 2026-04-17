package service

import (
	"Assignment-8/ex2/repository"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Aiym"}
	mockRepo.EXPECT().GetUserByID(1).Return(user, nil)

	result, err := userService.GetUserByID(1)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Aiym"}
	mockRepo.EXPECT().CreateUser(user).Return(nil)

	err := userService.CreateUser(user)
	assert.NoError(t, err)
}

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	user := &repository.User{ID: 2, Name: "New User", Email: "new@example.com"}

	t.Run("User already exists", func(t *testing.T) {
		mockRepo.EXPECT().GetByEmail(user.Email).Return(&repository.User{ID: 3}, nil)
		err := service.RegisterUser(user)
		assert.EqualError(t, err, "user with this email already exists")
	})

	t.Run("New User -> Success", func(t *testing.T) {
		mockRepo.EXPECT().GetByEmail(user.Email).Return(nil, nil)
		mockRepo.EXPECT().CreateUser(user).Return(nil)
		err := service.RegisterUser(user)
		assert.NoError(t, err)
	})

	t.Run("Repository error on GetByEmail", func(t *testing.T) {
		mockRepo.EXPECT().GetByEmail(user.Email).Return(nil, errors.New("db error"))
		err := service.RegisterUser(user)
		assert.ErrorContains(t, err, "error getting user by email")
	})

	t.Run("Repository error on CreateUser", func(t *testing.T) {
		mockRepo.EXPECT().GetByEmail(user.Email).Return(nil, nil)
		mockRepo.EXPECT().CreateUser(user).Return(errors.New("db error"))
		err := service.RegisterUser(user)
		assert.ErrorContains(t, err, "db error")
	})
}

func TestUpdateUserName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	t.Run("Empty name", func(t *testing.T) {
		err := service.UpdateUserName(1, "")
		assert.EqualError(t, err, "name cannot be empty")
	})

	t.Run("User not found/repo error", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(1).Return(nil, errors.New("not found"))
		err := service.UpdateUserName(1, "New Name")
		assert.EqualError(t, err, "not found")
	})

	t.Run("Successful update", func(t *testing.T) {
		user := &repository.User{ID: 1, Name: "Old"}
		mockRepo.EXPECT().GetUserByID(1).Return(user, nil)
		mockRepo.EXPECT().UpdateUser(gomock.Any()).DoAndReturn(func(u *repository.User) error {
			assert.Equal(t, "New Name", u.Name)
			return nil
		})
		err := service.UpdateUserName(1, "New Name")
		assert.NoError(t, err)
	})

	t.Run("UpdateUser fails", func(t *testing.T) {
		user := &repository.User{ID: 1, Name: "Old"}
		mockRepo.EXPECT().GetUserByID(1).Return(user, nil)
		mockRepo.EXPECT().UpdateUser(gomock.Any()).Return(errors.New("db error"))
		err := service.UpdateUserName(1, "New Name")
		assert.EqualError(t, err, "db error")
	})
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	t.Run("Attempt to delete admin", func(t *testing.T) {
		err := service.DeleteUser(1)
		assert.EqualError(t, err, "it is not allowed to delete admin user")
	})

	t.Run("Successful delete", func(t *testing.T) {
		mockRepo.EXPECT().DeleteUser(2).Return(nil)
		err := service.DeleteUser(2)
		assert.NoError(t, err)
	})

	t.Run("Repository error", func(t *testing.T) {
		mockRepo.EXPECT().DeleteUser(3).Return(errors.New("db error"))
		err := service.DeleteUser(3)
		assert.EqualError(t, err, "db error")
	})
}
