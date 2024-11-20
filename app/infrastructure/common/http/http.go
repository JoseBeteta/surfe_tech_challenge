package http

import (
	"github.com/JoseBeteta/surfe/app/domain"
	"log/slog"
	"net/http"
)

var errorMap = map[error]int{
	domain.InvalidArgument: http.StatusBadRequest,
}

// NewHttpMapper creates a http mapper for inventory handlers
func NewHttpMapper(logger slog.Logger) *Mapper {
	return NewMapper(errorMap, logger)
}
