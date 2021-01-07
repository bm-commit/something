package application

import (
	"something/internal/users/domain"
	"time"
)

// UserResponse ...
type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedOn time.Time `json:"created_on"`
}

// NewUserResponse ...
func NewUserResponse(User *domain.User) *UserResponse {
	return &UserResponse{
		ID:        User.ID,
		Name:      User.Name,
		Username:  User.Username,
		Email:     User.Email,
		Role:      User.Role,
		CreatedOn: User.CreatedOn,
	}
}

// NewUsersResponse ...
func NewUsersResponse(Users []*domain.User) []*UserResponse {
	usersResponse := []*UserResponse{}
	for _, user := range Users {
		usersResponse = append(usersResponse, NewUserResponse(user))
	}
	return usersResponse
}
