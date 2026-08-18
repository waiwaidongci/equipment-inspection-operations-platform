package database

import (
	"database/sql"
	"time"
)

const timeFormat = time.RFC3339Nano

func FormatTime(t time.Time) string {
	return t.Format(timeFormat)
}

func ParseTime(value string) (time.Time, error) {
	return time.Parse(timeFormat, value)
}

func NullTime(value sql.NullString) (*time.Time, error) {
	if !value.Valid {
		return nil, nil
	}
	t, err := ParseTime(value.String)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func NullInt64(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}
	v := value.Int64
	return &v
}

func Int64Pointer(value *int64) sql.NullInt64 {
	if value == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *value, Valid: true}
}
