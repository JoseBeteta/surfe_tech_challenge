package application_test

import (
	"encoding/json"
	application_action "github.com/JoseBeteta/surfe/app/application"
	domain_action "github.com/JoseBeteta/surfe/app/domain"
	common_http "github.com/JoseBeteta/surfe/app/infrastructure/common/http"
	"github.com/JoseBeteta/surfe/test/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock ActionReadRepository
type MockActionReadRepository struct {
	mock.Mock
}

func (m *MockActionReadRepository) CountByUserID(userID int) (int, error) {
	args := m.Called(userID)
	return args.Int(0), args.Error(1)
}

func (m *MockActionReadRepository) GetNextActionProbabilities(actionType string) (map[string]float64, error) {
	args := m.Called(actionType)
	return args.Get(0).(map[string]float64), args.Error(1)
}

func (m *MockActionReadRepository) GetAll() ([]domain_action.Action, error) {
	args := m.Called()
	return args.Get(0).([]domain_action.Action), args.Error(1)
}

// Test HandleGetActionCountInfo
func TestHandleGetActionCountInfo(t *testing.T) {
	mockRepo := new(MockActionReadRepository)
	logger := mocks.NewNullLogger()                 // Assuming you have a NullLogger for testing
	httpMapper := common_http.NewHttpMapper(logger) // Use real httpMapper, not mock

	handler := application_action.NewActionHandler(mockRepo, logger, httpMapper)

	mockRepo.On("CountByUserID", 1).Return(10, nil)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{
		{Key: "id", Value: "1"},
	}

	handler.HandleGetActionCountInfo(c)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response application_action.CountResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 10, response.Count)

	mockRepo.AssertExpectations(t)
}

// Test HandleGetNextActionProbability
func TestHandleGetNextActionProbability(t *testing.T) {
	mockRepo := new(MockActionReadRepository)
	logger := mocks.NewNullLogger()                 // Assuming you have a NullLogger for testing
	httpMapper := common_http.NewHttpMapper(logger) // Use real httpMapper, not mock

	handler := application_action.NewActionHandler(mockRepo, logger, httpMapper)

	mockRepo.On("GetNextActionProbabilities", "REFER_USER").Return(map[string]float64{
		"REFER_USER":   0.75,
		"VIEW_PROFILE": 0.25,
	}, nil)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{
		{Key: "action", Value: "REFER_USER"},
	}

	handler.HandleGetNextActionProbability(c)

	assert.Equal(t, http.StatusOK, rec.Code)

	var probabilities map[string]float64
	err := json.Unmarshal(rec.Body.Bytes(), &probabilities)
	assert.NoError(t, err)
	assert.Equal(t, map[string]float64{
		"REFER_USER":   0.75,
		"VIEW_PROFILE": 0.25,
	}, probabilities)

	mockRepo.AssertExpectations(t)
}

// Test HandleCalculationReferralIndex
func TestHandleCalculationReferralIndex(t *testing.T) {
	mockRepo := new(MockActionReadRepository)
	logger := mocks.NewNullLogger()                 // Assuming you have a NullLogger for testing
	httpMapper := common_http.NewHttpMapper(logger) // Use real httpMapper, not mock

	handler := application_action.NewActionHandler(mockRepo, logger, httpMapper)

	mockRepo.On("GetAll").Return([]domain_action.Action{
		{UserID: 1, Type: "REFER_USER", TargetUser: 2},
		{UserID: 2, Type: "REFER_USER", TargetUser: 3},
		{UserID: 1, Type: "REFER_USER", TargetUser: 4},
	}, nil)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	handler.HandleCalculationReferralIndex(c)

	assert.Equal(t, http.StatusOK, rec.Code)

	var referralIndex map[int]int
	err := json.Unmarshal(rec.Body.Bytes(), &referralIndex)
	assert.NoError(t, err)
	assert.Equal(t, map[int]int{
		1: 3,
		2: 1,
		3: 0,
		4: 0,
	}, referralIndex)

	mockRepo.AssertExpectations(t)
}
