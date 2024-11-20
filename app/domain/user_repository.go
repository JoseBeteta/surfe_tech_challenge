package domain

// UserReadRepository is the interface for the Repository used to fetch data from storage
type UserReadRepository interface {
	GetByID(id int) (User, error)
}
