package null

import (
	"database/sql"
	"encoding/json"
)

type String struct {
	sql.NullString
}

func NewString(arg string) String {
	temp := String{}
	temp.Valid = true
	temp.String = arg
	return temp
}

// MarshalJSON for String
func (ni String) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.String)
}

// UnmarshalJSON for String
func (ni *String) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ni.String)
	ni.Valid = err == nil
	return err
}
