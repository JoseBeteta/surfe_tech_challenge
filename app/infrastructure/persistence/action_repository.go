package persistence

import (
	"encoding/json"
	domainAction "github.com/JoseBeteta/surfe/app/domain"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"sync"
)

// ActionJSONRepository is a repository that interacts with a JSON file
type ActionJSONRepository struct {
	filePath string
	mutex    sync.Mutex
}

// NewActionJSONRepository creates a new repository that uses a JSON file
func NewActionJSONRepository(filePath string) *ActionJSONRepository {
	return &ActionJSONRepository{filePath: filePath}
}

// CountByUserID returns the count of actions for a given user ID
func (r *ActionJSONRepository) CountByUserID(userID int) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	actions, err := r.readFromFile()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, action := range actions {
		if action.UserID == userID {
			count++
		}
	}

	return count, nil
}

// GetNextActionProbabilities calculates the probabilities of next actions after a given action type
func (r *ActionJSONRepository) GetNextActionProbabilities(actionType string) (map[string]float64, error) {
	actions, err := r.readFromFile()
	if err != nil {
		return nil, err
	}

	// Sort actions by userId and createdAt
	sort.Slice(actions, func(i, j int) bool {
		if actions[i].UserID == actions[j].UserID {
			return actions[i].CreatedAt.Before(actions[j].CreatedAt)
		}
		return actions[i].UserID < actions[j].UserID
	})

	// Count transitions from the given actionType
	transitionCounts := make(map[string]int)
	totalTransitions := 0

	for i := 0; i < len(actions)-1; i++ {
		current := actions[i]
		next := actions[i+1]

		if current.UserID == next.UserID && current.Type == actionType {
			transitionCounts[next.Type]++
			totalTransitions++
		}
	}

	// Calculate probabilities
	probabilities := make(map[string]float64)
	for action, count := range transitionCounts {
		probabilities[action] = math.Round((float64(count)/float64(totalTransitions))*100) / 100
	}

	return probabilities, nil
}

// GetAll retrieves all actions from the JSON file
func (r *ActionJSONRepository) GetAll() ([]domainAction.Action, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	actions, err := r.readFromFile()
	if err != nil {
		return nil, err
	}

	return actions, nil
}

// readFromFile reads the action data from the JSON file
func (r *ActionJSONRepository) readFromFile() ([]domainAction.Action, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []domainAction.Action{}, nil
		}
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var actions []domainAction.Action
	err = json.Unmarshal(data, &actions)
	if err != nil {
		return nil, err
	}

	return actions, nil
}
