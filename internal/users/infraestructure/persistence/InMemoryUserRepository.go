package persistence

import (
	"errors"
	"something/internal/users/domain"
)

type repository struct {
	users map[string]*domain.User
}

var (
	userInstance *repository
)

// NewInMemoryUserRepository ...
func NewInMemoryUserRepository() domain.UserRepository {
	userInstance = &repository{
		users: make(map[string]*domain.User),
	}
	return userInstance
}

func (r *repository) Find() ([]*domain.User, error) {
	var users []*domain.User
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

func (r *repository) FindByID(id string) (*domain.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *repository) FindByEmail(email string) (*domain.User, error) {
	var user *domain.User
	found := false
	for _, u := range r.users {
		if email == u.Email {
			user = u
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("email not found")
	}
	return user, nil
}
func (r *repository) FindByUsername(username string) (*domain.User, error) {
	var user *domain.User
	found := false
	for _, u := range r.users {
		if username == u.Username {
			user = u
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("username not found")
	}
	return user, nil
}

func (r *repository) Update(user *domain.User) error {
	r.users[user.ID] = user
	return nil
}

func (r *repository) UpdateInterests(userID, bookID, status string) error {
	r.users[userID].Interests[bookID] = status
	return nil
}

func (r *repository) Save(user *domain.User) error {
	r.users[user.ID] = user
	return nil
}

func (r *repository) Delete(id string) error {
	_, ok := r.users[id]
	if ok {
		delete(r.users, id)
	}
	return nil
}
