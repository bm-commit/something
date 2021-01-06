package find

import (
	"something/internal/users/application"
	"something/internal/users/domain"
)

// Service ...
type Service interface {
	FindUsers() ([]*application.UserResponse, error)
	FindUserByID(id string) (*application.UserResponse, error)
	FindUserByUsername(username string) (*application.UserResponse, error)
}

type service struct {
	repository domain.UserRepository
}

// NewService ...
func NewService(repository domain.UserRepository) Service {
	return &service{repository: repository}
}

func (s *service) FindUsers() ([]*application.UserResponse, error) {
	users, err := s.repository.Find()
	if err != nil {
		return nil, err
	}
	return application.NewUsersResponse(users), nil
}

func (s *service) FindUserByID(id string) (*application.UserResponse, error) {
	user, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	return application.NewUserResponse(user), nil
}

func (s *service) FindUserByUsername(username string) (*application.UserResponse, error) {
	user, err := s.repository.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	return application.NewUserResponse(user), nil
}
