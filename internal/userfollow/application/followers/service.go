package followers

import (
	"errors"
	"something/internal/userfollow/domain"
	userDomain "something/internal/users/domain"
)

// Service ...
type Service interface {
	Follow(from, to string) error
	Unfollow(from, to string) error
}

type service struct {
	repository     domain.UserFollowRepository
	userRepository userDomain.UserRepository
}

// NewService ...
func NewService(repo domain.UserFollowRepository, userRepo userDomain.UserRepository) Service {
	return &service{repository: repo, userRepository: userRepo}
}

func (s *service) Follow(from, to string) error {
	existingUserID, _ := s.userRepository.FindByID(to)
	if existingUserID == nil {
		return errors.New("user not found")
	}
	userFollow, _ := domain.NewUserFollow(from, to)
	err := s.repository.Follow(userFollow)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Unfollow(from, to string) error {
	existingUserID, _ := s.userRepository.FindByID(to)
	if existingUserID == nil {
		return errors.New("user not found")
	}
	userFollow, _ := domain.NewUserFollow(from, to)
	err := s.repository.Unfollow(userFollow)
	if err != nil {
		return err
	}
	return nil
}
