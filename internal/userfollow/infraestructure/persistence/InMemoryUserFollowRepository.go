package persistence

import (
	"something/internal/userfollow/domain"
)

type repository struct {
	followers []*domain.UserFollow
}

var (
	userFollowInstance *repository
)

// NewInMemoryUserFollowRepository ...
func NewInMemoryUserFollowRepository() domain.UserFollowRepository {
	userFollowInstance = &repository{
		followers: make([]*domain.UserFollow, 15),
	}
	return userFollowInstance
}

func (r *repository) FindFollowing(id string) ([]*domain.UserFollow, error) {
	var following []*domain.UserFollow
	for _, follow := range r.followers {
		if follow == nil {
			continue
		}
		if follow.From == id {
			following = append(following, follow)
		}
	}
	return following, nil
}

func (r *repository) FindFollowers(id string) ([]*domain.UserFollow, error) {
	var followers []*domain.UserFollow
	for _, follow := range r.followers {
		if follow == nil {
			continue
		}
		if follow.To == id {
			followers = append(followers, follow)
		}
	}
	return followers, nil
}

func (r *repository) Follow(u *domain.UserFollow) error {
	r.followers = append(r.followers, u)
	return nil
}

func (r *repository) Unfollow(u *domain.UserFollow) error {
	var elementToRemove int
	for i, follow := range r.followers {
		if follow == nil {
			continue
		}
		if follow.From == u.From && follow.To == u.To {
			elementToRemove = i
			break
		}
	}
	r.followers[elementToRemove] = r.followers[len(r.followers)-1]
	r.followers[len(r.followers)-1] = nil
	r.followers = r.followers[:len(r.followers)-1]
	return nil
}
