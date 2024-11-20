package application

import (
	"github.com/JoseBeteta/surfe/app/domain"
	"github.com/JoseBeteta/surfe/app/infrastructure/common/http"
	"github.com/gin-gonic/gin"
	"log/slog"
	"strconv"
)

const (
	userIdParameterKey   = "id"
	actionIdParameterKey = "action"
	referUser            = "REFER_USER"
)

type ReferralGraph map[int][]int

// ActionHandler of action handler http requests
type ActionHandler struct {
	actionReadRepository domain.ActionReadRepository
	logger               slog.Logger
	httpMapper           *http.Mapper
}

// NewActionHandler creates a new handler for action info
func NewActionHandler(
	actionReadRepository domain.ActionReadRepository,
	logger slog.Logger,
	httpMapper *http.Mapper,
) *ActionHandler {
	return &ActionHandler{
		actionReadRepository,
		logger,
		httpMapper,
	}
}

func (h *ActionHandler) Initialize(r *gin.Engine, middlewares ...gin.HandlerFunc) {
	group := r.Group("api/actions")

	group.Use(
		http.Consume(http.V1),
		http.Produce(http.V1),
	)
	group.Use(middlewares...)

	group.GET("users/:id", h.HandleGetActionCountInfo)
	group.GET("probability/users/:action", h.HandleGetNextActionProbability)
	group.GET("referral", h.HandleCalculationReferralIndex)
}

type CountResponse struct {
	Count int `json:"count"`
}

// HandleGetActionCountInfo retrieves action count info
func (h *ActionHandler) HandleGetActionCountInfo(c *gin.Context) {
	idStr := c.Param(userIdParameterKey)

	userId, err := strconv.Atoi(idStr)
	if err != nil {
		h.httpMapper.ErrorResponse(c, err)
		return
	}

	count, err := h.actionReadRepository.CountByUserID(userId)
	if err != nil {
		h.logger.Warn("actions not found for this user", "id", userId)
		h.httpMapper.ErrorResponse(c, err)
		return
	}

	h.httpMapper.OkResponse(c, CountResponse{
		Count: count,
	})
}

// HandleGetNextActionProbability retrieves next action probability info
func (h *ActionHandler) HandleGetNextActionProbability(c *gin.Context) {
	action := c.Param(actionIdParameterKey)

	probabilities, err := h.actionReadRepository.GetNextActionProbabilities(action)
	if err != nil {
		h.logger.Warn("probabilities not found for this action", "action", action)
		h.httpMapper.ErrorResponse(c, err)
		return
	}

	h.httpMapper.OkResponse(c, probabilities)
}

// buildReferralGraph group users by referral into a graph
func buildReferralGraph(actions []domain.Action) ReferralGraph {
	graph := make(ReferralGraph)
	for _, action := range actions {
		if action.Type == referUser {
			graph[action.UserID] = append(graph[action.UserID], action.TargetUser)
		}
	}
	return graph
}

func (h *ActionHandler) HandleCalculationReferralIndex(c *gin.Context) {
	actions, err := h.actionReadRepository.GetAll()
	if err != nil {
		h.logger.Warn("actions not found")
		h.httpMapper.ErrorResponse(c, err)
		return
	}

	graph := buildReferralGraph(actions)

	referralIndex := make(map[int]int)
	processed := make(map[int]bool)

	var dfs func(userID int) int
	dfs = func(userID int) int {
		// If already calculated, return cached value
		if count, found := referralIndex[userID]; found {
			return count
		}

		count := 0
		processed[userID] = true

		for _, referredUser := range graph[userID] {
			// Avoid recalculation for already processed nodes
			if !processed[referredUser] {
				// Use recursivity to visit the nodes of the referred user
				count += 1 + dfs(referredUser)
				continue
			}
			// in case the referral it's already calculated
			count += 1 + referralIndex[referredUser]
		}

		referralIndex[userID] = count
		return count
	}

	// Compute referral index for all users in the graph
	for userID := range graph {
		if !processed[userID] {
			dfs(userID)
		}
	}

	// Return the computed referral index
	h.httpMapper.OkResponse(c, referralIndex)
}
