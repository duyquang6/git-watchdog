package repository

import (
	"encoding/json"
	"github.com/duyquang6/git-watchdog/pkg/dto"
	"github.com/duyquang6/git-watchdog/pkg/exception"
	_validator "github.com/duyquang6/git-watchdog/pkg/validator"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (s *Controller) HandleCreateRepository() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "read body failed").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		req := dto.CreateRepositoryRequest{}
		err = json.Unmarshal(data, &req)
		if err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "unmarshal failed").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		if err := _validator.GetValidate().Struct(req); err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "validation error").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		res, err := s.service.Create(ctx, req)
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

func (s *Controller) HandleUpdateRepository() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "read body failed").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		req := dto.UpdateRepositoryRequest{}
		err = json.Unmarshal(data, &req)
		if err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "unmarshal failed").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "parse failed").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		if err := _validator.GetValidate().Struct(req); err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "validation error").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		res, err := s.service.Update(ctx, uint(id), req)
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

func (s *Controller) HandleGetOneRepository() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)

		if err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "parse id failed").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		res, err := s.service.GetOne(ctx, uint(id))
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

func (s *Controller) HandleDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)

		if err != nil {
			appErr := exception.Wrap(http.StatusBadRequest, err, "parse id failed").(exception.AppError)
			c.JSON(appErr.GetHTTPStatusCode(), appErr.ToAppErrorResponse())
			return
		}

		res, err := s.service.Delete(ctx, uint(id))
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
