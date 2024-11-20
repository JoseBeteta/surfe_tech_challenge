package application_test

import (
	"encoding/json"
	user_application "github.com/JoseBeteta/surfe/app/application"
	domainUser "github.com/JoseBeteta/surfe/app/domain"
	"github.com/JoseBeteta/surfe/test/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	common_http "github.com/JoseBeteta/surfe/app/infrastructure/common/http"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock of the UserReadRepository
type MockUserReadRepository struct {
	mock.Mock
}

func (m *MockUserReadRepository) GetByID(id int) (domainUser.User, error) {
	args := m.Called(id)
	return args.Get(0).(domainUser.User), args.Error(1)
}

func TestHandleGetUserInfo(t *testing.T) {
	mockRepo := new(MockUserReadRepository)

	user := domainUser.User{
		ID:        1, // Ensure this is an integer
		Name:      "John Doe",
		CreatedAt: time.Date(2022, time.December, 12, 0, 0, 0, 0, time.UTC), // Hardcoded date: 2022-12-12
	}

	mockRepo.On("GetByID", 1).Return(user, nil)

	logger := mocks.NewNullLogger()
	httpMapper := common_http.NewHttpMapper(logger)
	handler := user_application.NewUserHandler(mockRepo, logger, httpMapper)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	c.Params = gin.Params{
		{Key: "id", Value: "1"},
	}

	handler.HandleGetUserInfo(c)

	assert.Equal(t, http.StatusOK, rec.Code)

	var userInfo user_application.UserInfoResponse
	err := json.Unmarshal(rec.Body.Bytes(), &userInfo)
	assert.NoError(t, err)

	assert.Equal(t, user.ID, userInfo.ID)
	assert.Equal(t, user.Name, userInfo.Name)
	assert.Equal(t, user.CreatedAt.Format(time.RFC3339), userInfo.CreatedAt)

	mockRepo.AssertExpectations(t)
}
