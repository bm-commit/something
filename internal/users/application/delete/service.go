package delete

import (
	"errors"
	"something/internal/users/domain"
)

// Service ...
type Service interface {
	DeleteUserByID(id string) error
}

type service struct {
	repository domain.UserRepository
}

// NewService ...
func NewService(repository domain.UserRepository) Service {
	return &service{repository: repository}
}

func (s *service) DeleteUserByID(id string) error {
	review, _ := s.repository.FindByID(id)
	if review == nil {
		return errors.New("user not found")
	}
	err := s.repository.Delete(id)
	return err
}
