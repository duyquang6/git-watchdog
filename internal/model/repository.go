package model

import (
	"fmt"

	"github.com/duyquang6/git-watchdog/pkg/validator"
	"gorm.io/gorm"
)

// Repository data model
type Repository struct {
	BaseModel
	Name string `gorm:"column:name" validate:"required"`
	URL  string `gorm:"column:url" validate:"url"`
}

// BeforeSave validate Repository model
func (t *Repository) BeforeSave(tx *gorm.DB) error {
	err := t.validateModel()
	if err != nil {
		return fmt.Errorf("can't save invalid wager: %w", err)
	}
	return nil
}

func (t *Repository) validateModel() error {
	err := validator.GetValidate().Struct(t)
	if err != nil {
		return fmt.Errorf("db field validation failed: %w", err)
	}

	return nil
}
