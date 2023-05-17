package util

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
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
