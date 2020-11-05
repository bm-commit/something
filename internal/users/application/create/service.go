package create

import (
	"errors"
	"something/internal/users/domain"
	"something/pkg/crypto"
)

// Service ...
type Service interface {
	CreateUser(*UserCommand) error
}

type service struct {
	repository domain.UserRepository
	cryptoRepo crypto.Crypto
}

// NewService ...
func NewService(repository domain.UserRepository, cryptoInstance crypto.Crypto) Service {
	return &service{repository: repository, cryptoRepo: cryptoInstance}
}

func (s *service) CreateUser(command *UserCommand) error {
	err := usecaseValidations(command, s.repository)
	if err != nil {
		return err
	}

	hashedPassword, err := s.cryptoRepo.Hash(command.Password)
	if err != nil {
		return err
	}

	user, err := domain.NewUser(command.ID, command.Name,
		command.Username, command.Email, hashedPassword)
	if err != nil {
		return err
	}
	err = s.repository.Save(user)
	if err != nil {
		return err
	}
	return nil
}

func usecaseValidations(command *UserCommand, repo domain.UserRepository) error {
	existingUserID, _ := repo.FindByID(command.ID)
	if existingUserID != nil {
		return errors.New("user id already exists")
	}
	existingUsername, _ := repo.FindByUsername(command.Username)
	if existingUsername != nil {
		return errors.New("username already in use")
	}
	existingEmail, _ := repo.FindByEmail(command.Email)
	if existingEmail != nil {
		return errors.New("email already in use")
	}
	return nil
}
