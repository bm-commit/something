package followers

import (
	"something/internal/userfollow/domain"
)

// Service ...
type Service interface {
	Follow(from, to string) error
	Unfollow(from, to string) error
}

type service struct {
	repository domain.UserFollowRepository
}

// NewService ...
func NewService(repo domain.UserFollowRepository) Service {
	return &service{repository: repo}
}

func (s *service) Follow(from, to string) error {
	userFollow, _ := domain.NewUserFollow(from, to)
	err := s.repository.Follow(userFollow)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Unfollow(from, to string) error {
	userFollow, _ := domain.NewUserFollow(from, to)
	err := s.repository.Unfollow(userFollow)
	if err != nil {
		return err
	}
	return nil
}
