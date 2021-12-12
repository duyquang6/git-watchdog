package null

import (
	"database/sql"
	"time"
)

type Time struct {
	sql.NullTime
}

func NewTime(arg time.Time) Time {
	_time := Time{}
	_time.Valid = true
	_time.Time = arg
	return _time
}
