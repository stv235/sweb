package sweb

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sweb/form"
	"time"
)

func RequirePanic(v ...interface{}) {
	str := fmt.Sprintln(v...)

	log.Output(3, str)
	panic(errors.New(str))
}

func Exists(values url.Values, key string) bool {
	_, ok := values[key]
	return ok
}

func RequireInt64(values url.Values, key string) int64 {
	valStr := values.Get(key)

	if valStr == "" {
		RequirePanic("required variable " + key + " empty")
	}

	valInt, err := strconv.ParseInt(valStr, 10, 64)

	if err != nil {
		log.Panicln(err)
	}

	return valInt
}

func RequireInt64s(values url.Values, key string) []int64 {
	valStrs, ok := values[key]

	if !ok {
		RequirePanic("required variable " + key + " empty")
	}

	valInts := make([]int64, 0)

	for _, valStr := range valStrs {
		if valStr == "" {
			RequirePanic("required variable " + key + " empty")
		}

		valInt, err := strconv.ParseInt(valStr, 10, 64)

		if err != nil {
			log.Panicln(err)
		}

		valInts = append(valInts, valInt)
	}

	return valInts
}

func OptionalInt64(values url.Values, key string) *int64 {
	valStr := values.Get(key)

	if valStr == "" {
		return nil
	}

	valInt := RequireInt64(values, key)

	return &valInt
}

func OptionalDate(values url.Values, key string) *time.Time {
	str := OptionalString(values, key)

	if str == "" {
		return nil
	}

	t, err := time.ParseInLocation(form.FormDateFormat, str, time.Local)

	if err != nil {
		log.Panicln("invalid date time format")
	}

	return &t
}

func OptionalDateTime(values url.Values, key1, key2 string) *time.Time {
	str1 := OptionalString(values, key1)

	if str1 == "" {
		return nil
	}

	str2 := OptionalString(values, key2)

	if str2 == "" {
		return nil
	}

	t, err := time.ParseInLocation(form.FormDateTimeFormat, str1 + "T" + str2, time.Local)

	if err != nil {
		log.Panicln("invalid date time format")
	}

	return &t
}

func RequireFloat64(values url.Values, key string) float64 {
	str := RequireString(values, key)

	val, err := strconv.ParseFloat(str, 64)

	if err != nil {
		log.Panicln(err)
	}

	return val
}

func RequireString(values url.Values, key string) string {
	valStr := values.Get(key)
	valStr = strings.TrimSpace(valStr)

	if valStr == "" {
		log.Panicln("required variable", key, "empty")
	}

	return valStr
}

func OptionalString(values url.Values, key string) string {
	valStr := values.Get(key)
	valStr = strings.TrimSpace(valStr)

	return valStr
}

func OptionalFloat64(values url.Values, key string) *float64 {
	str := OptionalString(values, key)

	if str == "" {
		return nil
	}

	val, err := strconv.ParseFloat(str, 64)

	if err != nil {
		return nil
	}

	return &val
}

func OptionalBool(values url.Values, key string) bool {
	str := OptionalString(values, key)

	return str != ""
}

func RequireDateTime(values url.Values, key string) time.Time {
	str := RequireString(values, key)

	t, err := time.ParseInLocation(form.FormDateTimeFormat, str, time.Local)

	if err != nil {
		log.Panicln("invalid date time format")
	}

	return t
}

func RequireDateTime2(values url.Values, dateKey string, timeKey string) time.Time {
	dateStr := RequireString(values, dateKey)
	timeStr := RequireString(values, timeKey)

	t, err := time.ParseInLocation(form.FormDateTimeFormat, dateStr + "T" + timeStr, time.Local)

	if err != nil {
		log.Panicln("invalid date time format")
	}

	return t
}

func RequireDate(values url.Values, key string) time.Time {
	val := OptionalDate(values, key)

	if val == nil {
		log.Panicln("required variable", key, "empty")
	}

	return *val
}
