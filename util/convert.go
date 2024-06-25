package util

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func ToString(value interface{}) (string, error) {
	switch reflect.ValueOf(value).Kind() {
	default:
		return "", errors.New("convert fail")
	case reflect.String:
		return value.(string), nil
	case reflect.Int:
		return strconv.Itoa(value.(int)), nil
	case reflect.Float32:
		return fmt.Sprintf("%v", value.(float32)), nil
	case reflect.Float64:
		return fmt.Sprintf("%v", value.(float64)), nil
	case reflect.Bool:
		if value.(bool) == true {
			return "true", nil
		} else {
			return "false", nil
		}
	}
}

func ToInt64(value interface{}) (int64, error) {
	switch reflect.ValueOf(value).Kind() {
	default:
		return 0, errors.New("convert fail")
	case reflect.Float64:
		return int64(value.(float64)), nil
	case reflect.String:
		return strconv.ParseInt(value.(string), 10, 64)
	}
}

func ToMonth(value int) (m time.Month, e error) {
	switch value {
	case 1:
		m = time.January
		return
	case 2:
		m = time.February
		return
	case 3:
		m = time.March
		return
	case 4:
		m = time.April
		return
	case 5:
		m = time.May
		return
	case 6:
		m = time.June
		return
	case 7:
		m = time.July
		return
	case 8:
		m = time.August
		return
	case 9:
		m = time.September
		return
	case 10:
		m = time.October
		return
	case 11:
		m = time.November
		return
	case 12:
		m = time.December
		return
	default:
		return m, errors.New("value is not month number")
	}
}
