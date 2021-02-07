package find

import (
	"something/internal/users/application"
	"something/internal/users/domain"
)

// PAGE Default pagination page
const PAGE int = 1

// PERPAGE Default page size (the number of items to return per page).
const PERPAGE int = 50

// Service ...
type Service interface {
	FindUsers(criteria *Criteria) ([]*application.UserResponse, error)
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

func (s *service) FindUsers(criteria *Criteria) ([]*application.UserResponse, error) {
	if criteria.Page == 0 {
		criteria.Page = PAGE
	}
	if criteria.PerPage == 0 || criteria.PerPage > 1000 {
		criteria.PerPage = PERPAGE
	}
	newUserCriteria := domain.NewUserCriteria(
		criteria.Page, criteria.PerPage, criteria.Query,
	)
	users, err := s.repository.Find(newUserCriteria)
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
