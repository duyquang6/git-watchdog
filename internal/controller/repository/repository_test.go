package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/duyquang6/git-watchdog/internal/service/mocks"
	"github.com/duyquang6/git-watchdog/pkg/dto"
	"github.com/duyquang6/git-watchdog/pkg/exception"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateRepository(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(mocks.RepositoryService)
	c := NewController(mockService)
	r.POST("/api/v1/repositories", c.HandleCreateRepository())

	t.Run("Success", func(t *testing.T) {
		mockResp := &dto.CreateRepositoryResponse{
			Meta: dto.Meta{
				Code:    http.StatusCreated,
				Message: http.StatusText(http.StatusCreated),
			},
		}
		mockReq := dto.CreateRepositoryRequest{
			Name: "1",
			URL:  "http://google.com",
		}

		mockService.On("Create", mock.Anything, mockReq).Return(mockResp, nil)

		rr := httptest.NewRecorder()

		payload, _ := json.Marshal(mockReq)
		body := bytes.NewReader(payload)
		request, err := http.NewRequest(http.MethodPost, "/api/v1/repositories", body)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		respBody, err := json.Marshal(mockResp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockService.AssertExpectations(t)
	})

	t.Run("CreateFailed", func(t *testing.T) {
		mockReq := dto.CreateRepositoryRequest{
			Name: "2",
			URL:  "http://google.com",
		}

		mockService.On("Create", mock.Anything, mockReq).Return(nil,
			exception.Wrap(exception.ErrInternalServer, errors.New("unexpected"), "create repository failed"))

		rr := httptest.NewRecorder()

		payload, _ := json.Marshal(mockReq)
		body := bytes.NewReader(payload)
		request, err := http.NewRequest(http.MethodPost, "/api/v1/repositories", body)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("URLValidationFailed", func(t *testing.T) {
		mockReq := dto.CreateRepositoryRequest{
			Name: "2",
			URL:  "1",
		}

		rr := httptest.NewRecorder()

		payload, _ := json.Marshal(mockReq)
		body := bytes.NewReader(payload)
		request, err := http.NewRequest(http.MethodPost, "/api/v1/repositories", body)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUpdateRepository(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(mocks.RepositoryService)
	c := NewController(mockService)
	r.PUT("/api/v1/repositories/:id", c.HandleUpdateRepository())

	t.Run("Success", func(t *testing.T) {
		mockReq := dto.UpdateRepositoryRequest{
			Name: "1",
			URL:  "http://google.com",
		}
		mockResp := &dto.UpdateRepositoryResponse{Meta: dto.Meta{
			Code:    200,
			Message: http.StatusText(200),
		}}

		mockService.On("Update", mock.Anything, uint(1), mockReq).Return(mockResp, nil)

		rr := httptest.NewRecorder()

		payload, _ := json.Marshal(mockReq)
		body := bytes.NewReader(payload)
		request, err := http.NewRequest(http.MethodPut, "/api/v1/repositories/1", body)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		respBody, err := json.Marshal(mockResp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockService.AssertExpectations(t)
	})

	t.Run("CreateFailed", func(t *testing.T) {
		mockReq := dto.UpdateRepositoryRequest{
			Name: "1",
			URL:  "http://google.com",
		}

		mockService.On("Update", mock.Anything, uint(2), mockReq).Return(nil,
			exception.Wrap(exception.ErrInternalServer, errors.New("unexpected"), "update repository failed"))

		rr := httptest.NewRecorder()

		payload, _ := json.Marshal(mockReq)
		body := bytes.NewReader(payload)
		request, err := http.NewRequest(http.MethodPut, "/api/v1/repositories/2", body)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("URLValidationFailed", func(t *testing.T) {
		mockReq := dto.UpdateRepositoryRequest{
			Name: "1",
			URL:  "2",
		}

		rr := httptest.NewRecorder()

		payload, _ := json.Marshal(mockReq)
		body := bytes.NewReader(payload)
		request, err := http.NewRequest(http.MethodPut, "/api/v1/repositories/3", body)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetOneRepository(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(mocks.RepositoryService)
	c := NewController(mockService)
	r.GET("/api/v1/repositories/:id", c.HandleGetOneRepository())

	t.Run("Success", func(t *testing.T) {
		mockResp := &dto.GetOneRepositoryResponse{
			Meta: dto.Meta{
				Code:    200,
				Message: http.StatusText(200),
			},
			Data: dto.Repository{
				ID:   1,
				Name: "1",
				URL:  "1",
			},
		}

		mockService.On("GetOne", mock.Anything, uint(1)).Return(mockResp, nil)

		rr := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api/v1/repositories/1", nil)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		respBody, err := json.Marshal(mockResp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockService.AssertExpectations(t)
	})

	t.Run("GetFailed", func(t *testing.T) {
		mockService.On("GetOne", mock.Anything, uint(2)).Return(nil,
			exception.Wrap(exception.ErrInternalServer, errors.New("unexpected"), "get repository failed"))

		rr := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api/v1/repositories/2", nil)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockService.AssertExpectations(t)
	})
}

func TestDeleteRepository(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(mocks.RepositoryService)
	c := NewController(mockService)
	r.DELETE("/api/v1/repositories/:id", c.HandleDelete())

	t.Run("Success", func(t *testing.T) {
		mockResp := &dto.DeleteRepositoryResponse{Meta: dto.Meta{
			Code:    200,
			Message: http.StatusText(200),
		}}

		mockService.On("Delete", mock.Anything, uint(1)).Return(mockResp, nil)

		rr := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodDelete, "/api/v1/repositories/1", nil)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		respBody, err := json.Marshal(mockResp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockService.AssertExpectations(t)
	})

	t.Run("DeleteFailed", func(t *testing.T) {
		mockService.On("Delete", mock.Anything, uint(2)).Return(nil,
			exception.Wrap(exception.ErrInternalServer, errors.New("unexpected"), "delete repository failed"))

		rr := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodDelete, "/api/v1/repositories/2", nil)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockService.AssertExpectations(t)
	})
}
