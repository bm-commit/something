package find

import (
	"errors"
	"something/internal/userfollow/application"
	"something/internal/userfollow/domain"
	userDomain "something/internal/users/domain"
)

// Service ...
type Service interface {
	Following(userID string) ([]*application.UserFollowResponse, error)
	Followers(userID string) ([]*application.UserFollowResponse, error)
}

type service struct {
	repository     domain.UserFollowRepository
	userRepository userDomain.UserRepository
}

// NewService ...
func NewService(repository domain.UserFollowRepository, userRepo userDomain.UserRepository) Service {
	return &service{repository: repository, userRepository: userRepo}
}

func (s *service) Following(id string) ([]*application.UserFollowResponse, error) {
	existingUserID, _ := s.userRepository.FindByID(id)
	if existingUserID == nil {
		return nil, errors.New("user not found")
	}
	following, err := s.repository.FindFollowing(id)
	if err != nil {
		return nil, err
	}
	return application.NewFollowsResponse(following), nil
}

func (s *service) Followers(id string) ([]*application.UserFollowResponse, error) {
	existingUserID, _ := s.userRepository.FindByID(id)
	if existingUserID == nil {
		return nil, errors.New("user not found")
	}
	followers, err := s.repository.FindFollowers(id)
	if err != nil {
		return nil, err
	}
	return application.NewFollowsResponse(followers), nil
}
