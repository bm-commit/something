package domain

//UserCriteria ...
type UserCriteria struct {
	Page    int64
	PerPage int64
	Query   string
}

// NewUserCriteria ...
func NewUserCriteria(page, perPage int, query string) *UserCriteria {
	return &UserCriteria{
		Page:    int64(page),
		PerPage: int64(perPage),
		Query:   query,
	}
}
