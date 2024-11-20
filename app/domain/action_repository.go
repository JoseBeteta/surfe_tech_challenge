package domain

// ActionReadRepository is the interface for the Repository used to fetch data from storage
type ActionReadRepository interface {
	CountByUserID(userID int) (int, error)
	GetNextActionProbabilities(actionType string) (map[string]float64, error)
	GetAll() ([]Action, error)
}
