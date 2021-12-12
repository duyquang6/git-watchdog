package dto

import (
	"github.com/duyquang6/git-watchdog/pkg/null"
	"gorm.io/datatypes"
)

// IssueScanResponse ...
type IssueScanResponse struct {
	Meta Meta `json:"meta"`
}

type ListScanResponse struct {
	Meta Meta   `json:"meta"`
	Data []Scan `json:"data"`
}

type Scan struct {
	ID         uint           `json:"id"`
	Repository Repository     `json:"repository"`
	Status     string         `json:"status"`
	QueuedAt   null.Uint      `json:"queuedAt"`
	ScanningAt null.Uint      `json:"scanningAt"`
	FinishedAt null.Uint      `json:"finishedAt"`
	Findings   datatypes.JSON `json:"findings"`
	Note       null.String    `json:"note"`
}
