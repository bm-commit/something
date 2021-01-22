package application

import (
	"something/internal/userfollow/domain"
	"time"
)

// UserFollowResponse ...
type UserFollowResponse struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	CreatedOn time.Time `json:"created_on"`
}

// UserFollowResponseLong allow return user follow details
type UserFollowResponseLong struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
	FollowAt time.Time `json:"follow_at"`
}

// NewFollowResponse ...
func newFollowResponse(uf *domain.UserFollow) *UserFollowResponse {
	return &UserFollowResponse{
		From:      uf.From,
		To:        uf.To,
		CreatedOn: uf.CreatedOn,
	}
}

// NewFollowsResponse ...
func NewFollowsResponse(follows []*domain.UserFollow) []*UserFollowResponse {
	userFollowResponse := []*UserFollowResponse{}
	for _, follow := range follows {
		userFollowResponse = append(userFollowResponse, newFollowResponse(follow))
	}
	return userFollowResponse
}
