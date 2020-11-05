package find

import (
	"something/internal/userfollow/application"
	"something/internal/userfollow/domain"
)

// Service ...
type Service interface {
	Following(userID string) ([]*application.UserFollowResponse, error)
	Followers(userID string) ([]*application.UserFollowResponse, error)
}

type service struct {
	repository domain.UserFollowRepository
}

// NewService ...
func NewService(repository domain.UserFollowRepository) Service {
	return &service{repository: repository}
}

func (s *service) Following(id string) ([]*application.UserFollowResponse, error) {
	following, err := s.repository.FindFollowing(id)
	if err != nil {
		return nil, err
	}
	return application.NewFollowsResponse(following), nil
}

func (s *service) Followers(id string) ([]*application.UserFollowResponse, error) {
	followers, err := s.repository.FindFollowers(id)
	if err != nil {
		return nil, err
	}
	return application.NewFollowsResponse(followers), nil
}
