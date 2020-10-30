package login

import (
	"errors"
	"something/internal/users/application"
	"something/internal/users/domain"
	"something/pkg/crypto"
)

// Service ...
type Service interface {
	Login(*Command) (*application.UserResponse, error)
}

type service struct {
	repository domain.UserRepository
	cryptoRepo crypto.Crypto
}

// NewService ...
func NewService(repository domain.UserRepository, cryptoInstance crypto.Crypto) Service {
	return &service{repository: repository, cryptoRepo: cryptoInstance}
}

func (s *service) Login(c *Command) (*application.UserResponse, error) {
	user, err := s.repository.FindByEmail(c.Email)
	if err != nil {
		return nil, err
	}
	valid := s.cryptoRepo.CompareHashAndText(c.Password, user.Password)
	if !valid {
		return nil, errors.New("invalid email or password")
	}
	return application.NewUserResponse(user), nil
}
