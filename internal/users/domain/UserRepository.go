package domain

// UserRepository ...
type UserRepository interface {
	Find(*UserCriteria) ([]*User, error)
	FindByID(string) (*User, error)
	FindByEmail(string) (*User, error)
	FindByUsername(string) (*User, error)
	Update(*User) error
	UpdateInterests(string, string, string) error
	Save(*User) error
	Delete(string) error
	DeleteInterest(string, string) error
}
