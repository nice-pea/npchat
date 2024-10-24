package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/text/language"
)

// Функция для парсинга JSON из тела запроса
func parseJSONRequest(body io.ReadCloser, v any) error {
	// Декодируйте JSON из тела запроса
	decoder := json.NewDecoder(body)
	defer body.Close() // Закрываем тело запроса после декодирования
	return decoder.Decode(v)
}

func locale(acceptLanguage, defaults string) string {
	matcher := language.NewMatcher([]language.Tag{
		language.AmericanEnglish,
		language.English,
		language.Russian,
	})
	if tags, _, err := language.ParseAcceptLanguage(acceptLanguage); err != nil {
		log.Printf("locale: language.ParseAcceptLanguage: %v", err)
		return defaults
	} else {
		tag, _, _ := matcher.Match(tags...)
		base, _ := tag.Base()
		reg, _ := tag.Region()
		return base.String() + "_" + reg.String()
	}
}

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

//nolint:unused
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
