package repository

import (
	"github.com/duyquang6/git-watchdog/pkg/exception"
	"github.com/duyquang6/git-watchdog/pkg/null"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// HandleIssueScan godoc
// @Summary      create scan task
// @Description  create task scan repository
// @Tags         repository
// @Accept       json
// @Produce      json
// @Param   id      path   uint     true  "Repository ID"
//@Success      200  {object}  dto.IssueScanResponse
//@Failure      400  {object}  exception.AppErrorResponse
//@Failure      404  {object}  exception.AppErrorResponse
//@Failure      500  {object}  exception.AppErrorResponse
// @Router       /repositories/{id}/scans [post]
func (s *Controller) HandleIssueScan() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		repoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "parse id field failed").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		res, err := s.service.IssueScanRepo(ctx, uint(repoID))
		if err != nil {
			if appErr, ok := err.(exception.AppError); ok {
				c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
				return
			}
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.JSON(res.Meta.Code, res)
	}
}

// HandleListScan godoc
// @Summary      get task scans
// @Description  get task scans
// @Tags         repository
// @Accept       json
// @Produce      json
// @Param   id      path   uint     true  "Repository ID"
// @Param   page      query   uint     true  "Page number" default(1)
// @Param   limit      query   uint     true  "Page size" default(10)
//@Success      200  {object}  dto.IssueScanResponse
//@Failure      400  {object}  exception.AppErrorResponse
//@Failure      404  {object}  exception.AppErrorResponse
//@Failure      500  {object}  exception.AppErrorResponse
// @Router       /repositories/{id}/scans [get]
func (s *Controller) HandleListScan() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		page, err := strconv.ParseUint(c.Query("page"), 10, 64)
		if err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "parse page field failed").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}
		if page == 0 {
			appErr := exception.Wrap(http.StatusBadRequest, err, "page field must be larger than 0").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}
		limit, err := strconv.ParseUint(c.Query("limit"), 10, 64)
		if err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "parse limit field failed").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}
		if limit == 0 {
			appErr := exception.Wrap(http.StatusBadRequest, err, "limit field must be larger than 0").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}
		repoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "parse id field failed").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		if err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "validation error").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		res, err := s.service.ListScan(ctx, null.NewUint(uint(repoID)), uint(page), uint(limit))
		if err != nil {
			if appErr, ok := err.(exception.AppError); ok {
				c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
				return
			}
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.JSON(res.Meta.Code, res)
	}
}
