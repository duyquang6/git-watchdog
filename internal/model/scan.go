package model

import (
	"github.com/duyquang6/git-watchdog/pkg/customtypes"
	"github.com/duyquang6/git-watchdog/pkg/null"
	"gorm.io/datatypes"
)

// Scan data model
type Scan struct {
	BaseModel
	RepositoryID uint                     `gorm:"column:repository_id"`
	Repository   Repository               `gorm:"foreignKey:RepositoryID" validate:"-"`
	Status       customtypes.ResultStatus `gorm:"column:status"`
	QueuedAt     null.Time                `gorm:"column:queued_at"`
	ScanningAt   null.Time                `gorm:"column:scanning_at"`
	FinishedAt   null.Time                `gorm:"column:finished_at"`
	Findings     datatypes.JSON           `gorm:"column:findings"`
	Note         null.String              `gorm:"column:note"`
}
