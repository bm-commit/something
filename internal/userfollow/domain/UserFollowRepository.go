package domain

// UserFollowRepository ...
type UserFollowRepository interface {
	FindFollowing(string) ([]*UserFollow, error)
	FindFollowers(string) ([]*UserFollow, error)
	Follow(*UserFollow) error
	Unfollow(*UserFollow) error
}
