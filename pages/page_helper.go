package pages

import "reflect"

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
