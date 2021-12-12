package repository

import (
	"github.com/duyquang6/git-watchdog/internal/service"
)

type Controller struct {
	service service.RepositoryService
}

// NewController creates a new controller.
func NewController(purchaseService service.RepositoryService) *Controller {
	return &Controller{
		service: purchaseService,
	}
}
