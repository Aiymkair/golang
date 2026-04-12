package repo

import (
	"Assignment-7/internal/entity"
	"errors"
	"sync"

	"github.com/google/uuid"
)

type UserRepo struct {
	users map[uuid.UUID]*entity.User
	mu    sync.RWMutex
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		users: make(map[uuid.UUID]*entity.User),
	}
}

func (r *UserRepo) CreateUser(user *entity.User) (*entity.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return nil, errors.New("user already exists")
	}
	r.users[user.ID] = user
	return user, nil
}

func (r *UserRepo) FindByUsername(username string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *UserRepo) FindByID(id uuid.UUID) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *UserRepo) UpdateRole(id uuid.UUID, newRole string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return errors.New("user not found")
	}
	user.Role = newRole
	return nil
}
