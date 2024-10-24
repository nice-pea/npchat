package http

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func uintsParam(values url.Values, param string) ([]uint, error) {
	valStr := values.Get(param)
	if valStr == "" {
		return nil, nil
	}

	slice := strings.Split(valStr, ",")
	var res = make([]uint, 0, len(slice))
	for _, str := range slice {
		val, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("некорректный параметр %s (массив uint через запятую): %w", param, err)
		}
		res = append(res, uint(val))
	}

	return res, nil
}

func uintOptionalParam(values url.Values, param string) (uint, bool, error) {
	valStr := values.Get(param)
	if valStr == "" {
		return 0, false, nil
	}

	v, err := strconv.ParseUint(valStr, 10, 64)
	if err != nil {
		return 0, false, fmt.Errorf("некорректный параметр %s (uint): %w", param, err)
	}

	return uint(v), true, nil
}
