package repository

import (
	"github.com/duyquang6/git-watchdog/internal/model"
	"github.com/duyquang6/git-watchdog/pkg/null"
	"gorm.io/gorm"
)

type scanRepo struct{}

// ScanRepository provide interface interact with Scan model
type ScanRepository interface {
	Create(tx *gorm.DB, scan *model.Scan) error
	Update(tx *gorm.DB, scan *model.Scan) error
	Delete(tx *gorm.DB, id uint) error
	List(tx *gorm.DB, repoID null.Uint, offset, limit uint) ([]model.Scan, error)
	GetByID(tx *gorm.DB, id uint) (model.Scan, error)
}

// NewScanRepository create ScanRepository concrete object
func NewScanRepository() *scanRepo {
	return &scanRepo{}
}

// Create new scan
func (s *scanRepo) Create(tx *gorm.DB, scan *model.Scan) error {
	res := tx.Create(scan)
	return res.Error
}

// Update scan
func (s *scanRepo) Update(tx *gorm.DB, data *model.Scan) error {
	res := tx.Model(data).Updates(data)
	if res.Error == nil && res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return res.Error
}

// List get scans result with pagination support
func (s *scanRepo) List(tx *gorm.DB, repoID null.Uint, offset, limit uint) ([]model.Scan, error) {
	var scans []model.Scan
	res := tx.Preload("Repository").Offset(int(offset)).Limit(int(limit))
	if repoID.Valid {
		res.Where("repository_id = ?", repoID.Uint)
	}
	return scans, res.Find(&scans).Error
}

func (s *scanRepo) Delete(tx *gorm.DB, id uint) error {
	return tx.Delete(&model.Scan{}, id).Error
}

func (s *scanRepo) GetByID(tx *gorm.DB, id uint) (model.Scan, error) {
	var data model.Scan
	return data, tx.Preload("Repository").First(&data, id).Error
}
