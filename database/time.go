package database

import (
	"database/sql/driver"
	"errors"
	"sweb/form"
	"time"
)

type Time struct {
	time.Time
}

const DefaultTimeFormat = "2006-01-02 15:04"
const DefaultDateFormat = "2006-01-02"

func (t Time) Value() (driver.Value, error) {
	return t.UTC().Unix(), nil
}

func (t *Time) Scan(src interface{}) error {
	if src != nil {
		switch src.(type) {
		case int64:
			*t = Time{ Time: time.Unix(src.(int64), 0) }
			return nil
		default:
			return errors.New("cannot convert to database.Time")
		}
	}

	return nil
}

func (t Time) String() string {
	return t.Time.Format(DefaultTimeFormat)
}

func FromTime(t time.Time) Time {
	return Time{ Time: t }
}

func NewTime(t time.Time) *Time {
	return &Time{ Time: t }
}

func OptionalTime(t *time.Time) *Time {
	if t == nil {
		return nil
	}

	return NewTime(*t)
}

func (t Time) FormatForm() string {
	return t.Time.Format(form.FormDateTimeFormat)
}

func (t Time) FormatFormDate() string {
	return t.Time.Format(form.FormDateFormat)
}

func (t Time) FormatFormTime() string {
	return t.Time.Format(form.FormTimeFormat)
}