package application

import (
	domainUser "github.com/JoseBeteta/surfe/app/domain"
	"github.com/JoseBeteta/surfe/app/infrastructure/common/http"
	"log/slog"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	userIDParameterKey = "id"
)

// UserHandler of user handler http requests
type UserHandler struct {
	userReadRepository domainUser.UserReadRepository
	logger             slog.Logger
	httpMapper         *http.Mapper
}

// NewUserHandler creates a new handler for user info
func NewUserHandler(
	userReadRepository domainUser.UserReadRepository,
	logger slog.Logger,
	httpMapper *http.Mapper,
) *UserHandler {
	return &UserHandler{
		userReadRepository,
		logger,
		httpMapper,
	}
}

func (h *UserHandler) Initialize(r *gin.Engine, middlewares ...gin.HandlerFunc) {
	group := r.Group("api/users")

	group.Use(
		http.Consume(http.V1),
		http.Produce(http.V1),
	)
	group.Use(middlewares...)

	group.GET("/:id", h.HandleGetUserInfo)
}

type UserInfoResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

// HandleGetUserInfo retrieves user info
func (h *UserHandler) HandleGetUserInfo(c *gin.Context) {
	idStr := c.Param(userIDParameterKey)

	userId, err := strconv.Atoi(idStr)
	if err != nil {
		h.httpMapper.ErrorResponse(c, err)
		return
	}

	user, err := h.userReadRepository.GetByID(userId)
	if err != nil {
		h.logger.Warn("user not found", "id", userId)
		h.httpMapper.ErrorResponse(c, err)
		return
	}

	h.httpMapper.OkResponse(c, UserInfoResponse{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	})
}
