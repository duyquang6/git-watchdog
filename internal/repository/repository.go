package repository

import (
	"github.com/duyquang6/git-watchdog/internal/model"
	"gorm.io/gorm"
)

type repositoryRepo struct{}

// RepoRepository provide interface interact with Scan model
type RepoRepository interface {
	Create(tx *gorm.DB, data *model.Repository) error
	GetByID(tx *gorm.DB, id uint) (*model.Repository, error)
	Update(tx *gorm.DB, data *model.Repository) error
	Delete(tx *gorm.DB, id uint) error
}

// NewRepoRepository create RepoRepository concrete object
func NewRepoRepository() *repositoryRepo {
	return &repositoryRepo{}
}

// Create new repository
func (s *repositoryRepo) Create(tx *gorm.DB, data *model.Repository) error {
	return tx.Create(data).Error
}

// GetByID repository
func (s *repositoryRepo) GetByID(tx *gorm.DB, id uint) (*model.Repository, error) {
	var data model.Repository
	return &data, tx.First(&data, id).Error
}

// Update repository
func (s *repositoryRepo) Update(tx *gorm.DB, data *model.Repository) error {
	res := tx.Model(data).Updates(*data)
	if res.Error == nil && res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return res.Error
}

// Delete repository
func (s *repositoryRepo) Delete(tx *gorm.DB, id uint) error {
	res := tx.Delete(&model.Repository{}, id)
	if res.Error == nil && res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return res.Error
}
