package http_test

import (
	"encoding/json"
	"errors"
	appHTTP "github.com/JoseBeteta/surfe/app/infrastructure/common/http"
	"net/http"
	"net/http/httptest"
	"testing"

	mocks "github.com/JoseBeteta/surfe/test/mocks"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/assert"
)

var (
	ErrFakeErrorBadRequest = errors.New("bad request")
	ErrFakeErrorNotFound   = errors.New("not found")
	ErrServerError         = errors.New("server error")
)

var errorMap = map[error]int{
	ErrFakeErrorBadRequest: http.StatusBadRequest,
	ErrFakeErrorNotFound:   http.StatusNotFound,
	ErrServerError:         http.StatusInternalServerError,
}

type stockExample struct {
	SKU      string `json:"sku"`
	Quantity uint   `json:"quantity"`
}

func TestOkResponse(t *testing.T) {
	httpMapper := appHTTP.NewMapper(errorMap, mocks.NewNullLogger())
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{}

	se := stockExample{"sku1", 10}

	httpMapper.OkResponse(c, se)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, http.StatusOK, c.Writer.Status())
	assert.JSONEq(t, `{"sku":"sku1","quantity":10}`, w.Body.String())
}

type dummy struct {
	X int `binding:"required"`
}

type dummy2 struct {
	X []int `binding:"required,min=1,max=2,dive,required"`
}

type dummy3 struct {
	X int `binding:"required" json:"x_json_name"`
}

func TestErrorHandling(t *testing.T) {
	r := gin.New()
	gin.SetMode(gin.TestMode)

	httpMapper := appHTTP.NewMapper(errorMap, mocks.NewNullLogger())
	httpMapper.Initialize(r)

	err1 := binding.Validator.ValidateStruct(&dummy{})
	err2 := binding.Validator.ValidateStruct(&dummy2{X: []int{}})
	err3 := binding.Validator.ValidateStruct(&dummy2{X: []int{1, 2, 3}})
	err4 := binding.Validator.ValidateStruct(&dummy3{})

	tests := []struct {
		name         string
		error        error
		expectedCode int
		expectedBody string
	}{
		{
			"testing nil",
			nil,
			http.StatusOK,
			`{"error":""}`,
		},
		{
			"testing bad request",
			ErrFakeErrorBadRequest,
			http.StatusBadRequest,
			`{"error":"bad request"}`,
		},
		{
			"testing syntax error",
			&json.SyntaxError{},
			http.StatusBadRequest,
			`{"error":""}`,
		},
		{
			"testing not found",
			ErrFakeErrorNotFound,
			http.StatusNotFound,
			`{"error":"not found"}`,
		},
		{
			"testing server error",
			ErrServerError,
			http.StatusInternalServerError,
			`{"error":"server error"}`,
		},
		{
			"testing any other error",
			errors.New("some new error"),
			http.StatusInternalServerError,
			`{"error":"some new error"}`,
		},
		{
			"custom message: field required",
			err1,
			http.StatusBadRequest,
			`{"error":"field 'X' is required"}`,
		},
		{
			"custom message: array field empty",
			err2,
			http.StatusBadRequest,
			`{"error":"field 'X' is required"}`,
		},
		{
			"custom message: array field too many items",
			err3,
			http.StatusBadRequest,
			`{"error":"field 'X' accepts up to 2 items"}`,
		},
		{
			"custom message: getting name from json tag",
			err4,
			http.StatusBadRequest,
			`{"error":"field 'x_json_name' is required"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = &http.Request{}

			httpMapper.ErrorResponse(c, tt.error)
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedCode, c.Writer.Status())
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestRecovery(t *testing.T) {

	t.Run("recovery with string message", func(t *testing.T) {

		httpMapper := appHTTP.NewMapper(errorMap, mocks.NewNullLogger())

		w := httptest.NewRecorder()
		_, engine := gin.CreateTestContext(w)

		engine.
			Use(httpMapper.Recovery()).
			Handle(http.MethodGet, "/test", func(c *gin.Context) { panic("run!") })

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"run!"}`, w.Body.String())
	})

	t.Run("recovery with not string message", func(t *testing.T) {

		httpMapper := appHTTP.NewMapper(errorMap, mocks.NewNullLogger())

		w := httptest.NewRecorder()
		_, engine := gin.CreateTestContext(w)

		engine.
			Use(httpMapper.Recovery()).
			Handle(http.MethodGet, "/test", func(c *gin.Context) { panic(1234) })

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Empty(t, w.Body.String())
	})
}
