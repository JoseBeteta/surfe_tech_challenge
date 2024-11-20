package persistence

import (
	"encoding/json"
	"errors"
	domainUser "github.com/JoseBeteta/surfe/app/domain"
	"io/ioutil"
	"os"
	"sync"
)

// UserJSONRepository is a repository that interacts with a JSON file
type UserJSONRepository struct {
	filePath string
	mutex    sync.Mutex
}

// NewUserJSONRepository creates a new repository that uses a JSON file
func NewUserJSONRepository(filePath string) *UserJSONRepository {
	return &UserJSONRepository{filePath: filePath}
}

// GetByID retrieves a user by ID
func (r *UserJSONRepository) GetByID(id int) (domainUser.User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	users, err := r.readFromFile()
	if err != nil {
		return domainUser.User{}, err
	}

	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}
	return domainUser.User{}, errors.New("user not found")
}

// readFromFile reads the user data from the JSON file
func (r *UserJSONRepository) readFromFile() ([]domainUser.User, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []domainUser.User{}, nil
		}
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	users := []domainUser.User{}
	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}
