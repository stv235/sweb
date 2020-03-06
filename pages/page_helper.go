package pages

import (
	"errors"
	"reflect"
	"sweb/database"
	"sweb/form"
	"time"
)

var ErrUnsupportedType = errors.New("unsupported type")

type PageHelper struct {
}

func (PageHelper) FormatEqual(v1, v2 interface{}, str string) string {
	if reflect.DeepEqual(v1, v2) {
		return str
	}

	return ""
}

func (PageHelper) FormatIs(str string, values ...bool) string {
	for _, val := range values {
		if !val {
			return ""
		}
	}

	return str
}

func (PageHelper) FormatFormDate(t interface{}) string {
	switch t.(type) {
	case database.Time:
		return t.(database.Time).Format(form.FormDateFormat)
	case time.Time:
		return t.(time.Time).Format(form.FormDateFormat)
	}

	panic(ErrUnsupportedType)
}

func (PageHelper) FormatFormTime(t interface{}) string {
	switch t.(type) {
	case database.Time:
		return t.(database.Time).Format(form.FormTimeFormat)
	case time.Time:
		return t.(time.Time).Format(form.FormTimeFormat)
	}

	panic(ErrUnsupportedType)
}