package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Mapper struct {
	errorMap map[error]int
	Logger   slog.Logger
}

// NewMapper creates a http mapper
func NewMapper(errorMap map[error]int, logger slog.Logger) *Mapper {
	return &Mapper{errorMap: errorMap, Logger: logger}
}

func (e *Mapper) Initialize(r *gin.Engine) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// this ensures validation error uses tag name defined in json tag instead of struct name
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			return fld.Tag.Get("json")
		})
	}

	r.Use(e.Recovery())
}

// OkResponse writes http 200 ok and json response
func (e *Mapper) OkResponse(c *gin.Context, obj any) {
	okResponseJson(c, obj)
}

func (e *Mapper) ErrorResponse(c *gin.Context, err error) {
	statusCode := e.getStatusCode(err)

	if isClientError(statusCode) {
		e.Logger.Warn("error in http request")
	} else {
		e.Logger.Error("error in the server")
	}

	errorResponseJson(c, statusCode, getErrorMessage(err))
}

func isClientError(statusCode int) bool {
	return statusCode >= 400 && statusCode < 500
}

func (e *Mapper) getStatusCode(incomingErr error) int {
	if incomingErr == nil {
		return http.StatusOK
	}

	for err, status := range e.errorMap {
		if errors.Is(incomingErr, err) {
			return status
		}
	}

	if isBadRequest(incomingErr) {
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}

func isBadRequest(err error) bool {
	return isJSONSyntaxError(err) ||
		isValidationError(err) ||
		errors.Is(err, ErrEmptyBody)
}

func isJSONSyntaxError(err error) bool {
	var target *json.SyntaxError

	return errors.As(err, &target)
}

func isValidationError(err error) bool {
	var errs validator.ValidationErrors

	return errors.As(err, &errs)
}

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one
func (e *Mapper) Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(handleRecovery)
}

func getErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		e := errs[0]
		switch e.Tag() {
		case "required":
			return fmt.Sprintf("field '%s' is required", e.Field())
		case "min":
			if e.Kind() == reflect.Slice && e.Param() == "1" {
				return fmt.Sprintf("field '%s' is required", e.Field())
			}
		case "max":
			if e.Kind() == reflect.Slice {
				return fmt.Sprintf("field '%s' accepts up to %s items", e.Field(), e.Param())
			}
		}
	}

	return err.Error()
}
