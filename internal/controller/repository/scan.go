package repository

import (
	"github.com/duyquang6/git-watchdog/pkg/exception"
	"github.com/duyquang6/git-watchdog/pkg/null"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

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

func (s *Controller) HandleListScan() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		page, err := strconv.ParseUint(c.Param("page"), 10, 64)
		if err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "parse page field failed").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}
		limit, err := strconv.ParseUint(c.Param("limit"), 10, 64)
		if err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "parse limit field failed").(exception.AppError)
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
