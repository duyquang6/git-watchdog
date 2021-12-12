package repository

import (
	"encoding/json"
	"errors"
	"github.com/duyquang6/git-watchdog/internal/service/mocks"
	"github.com/duyquang6/git-watchdog/pkg/dto"
	"github.com/duyquang6/git-watchdog/pkg/exception"
	"github.com/duyquang6/git-watchdog/pkg/null"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListScan(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(mocks.RepositoryService)
	c := NewController(mockService)
	r.GET("/api/v1/repositories/:id/scans", c.HandleListScan())

	t.Run("Success", func(t *testing.T) {
		mockResp := &dto.ListScanResponse{
			Meta: dto.PaginationMeta{
				Meta: dto.Meta{
					Code:    http.StatusOK,
					Message: http.StatusText(http.StatusOK),
				},
				Total: 10,
			},
			Data: []dto.Scan{
				{
					ID: 1,
				},
			},
		}

		mockService.On("ListScan", mock.Anything, null.NewUint(1), uint(1), uint(1)).Return(mockResp, nil)

		rr := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api/v1/repositories/1/scans?page=1&limit=1", nil)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		respBody, err := json.Marshal(mockResp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockService.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		mockService.On("ListScan", mock.Anything, null.NewUint(2), uint(1), uint(1)).
			Return(nil, exception.
				Wrap(exception.ErrInternalServer, errors.New("unexpected"), "list scan failed"))

		rr := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api/v1/repositories/2/scans?page=1&limit=1", nil)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("PageParamInvalid", func(t *testing.T) {
		rr := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api/v1/repositories/2/scans?page=0&limit=1", nil)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockService.AssertExpectations(t)
	})
	t.Run("LimitParamInvalid", func(t *testing.T) {
		rr := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api/v1/repositories/2/scans?page=1&limit=0", nil)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockService.AssertExpectations(t)
	})
}

func TestIssueScan(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(mocks.RepositoryService)
	c := NewController(mockService)
	r.POST("/api/v1/repositories/:id/scans", c.HandleIssueScan())

	t.Run("Success", func(t *testing.T) {
		mockResp := &dto.IssueScanResponse{
			Meta: dto.Meta{
				Code:    http.StatusCreated,
				Message: http.StatusText(http.StatusCreated),
			},
		}

		mockService.On("IssueScanRepo", mock.Anything, uint(1)).Return(mockResp, nil)

		rr := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/v1/repositories/1/scans", nil)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		respBody, err := json.Marshal(mockResp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockService.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		mockService.On("IssueScanRepo", mock.Anything, uint(2)).Return(nil,
			exception.Wrap(exception.ErrInternalServer, errors.New("unexpected"), "list scan failed"))

		rr := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/v1/repositories/2/scans", nil)
		assert.NoError(t, err)

		r.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockService.AssertExpectations(t)
	})
}
