package update

import (
	"encoding/json"
	"errors"
	"something/internal/users/domain"
	"strings"
)

// Service ...
type Service interface {
	UpdateUserByID(*UserCommand) error
	UpdateUserInterests(*UserInterestsCommand) error
}

type service struct {
	repository domain.UserRepository
}

// NewService ...
func NewService(repository domain.UserRepository) Service {
	return &service{repository: repository}
}

func (s *service) UpdateUserByID(user *UserCommand) error {
	existingUser, _ := s.repository.FindByID(user.ID)
	if existingUser == nil {
		return errors.New("user not found")
	}

	out, err := json.Marshal(user)
	if err != nil {
		return err
	}
	// Merge new book fields with existing book fields
	err = json.NewDecoder(strings.NewReader(string(out))).Decode(existingUser)
	if err != nil {
		return err
	}

	updatedUser, err := domain.NewUser(
		existingUser.ID, existingUser.Name, existingUser.Username,
		existingUser.Email, existingUser.Password)
	if err != nil {
		return err
	}
	updatedUser.CreatedOn = existingUser.CreatedOn

	err = s.repository.Update(updatedUser)
	return err
}

func (s *service) UpdateUserInterests(interestCommand *UserInterestsCommand) error {
	err := s.repository.UpdateInterests(
		interestCommand.UserID,
		interestCommand.BookID,
		interestCommand.Status,
	)
	return err
}
