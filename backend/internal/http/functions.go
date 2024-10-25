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

	"github.com/saime-0/nice-pea-chat/internal/app/optional"
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

func uintOptionalParam(values url.Values, param string) (optional.Uint, error) {
	valStr := values.Get(param)
	if valStr == "" {
		return optional.Uint{}, nil
	}

	v, err := strconv.ParseUint(valStr, 10, 64)
	if err != nil {
		return optional.Uint{}, fmt.Errorf("некорректный uint параметр %s: %w", param, err)
	}

	return optional.NewUint(uint(v)), nil
}

func boolOptionalParam(values url.Values, param string) (optional.Bool, error) {
	valStr := values.Get(param)
	if valStr == "" {
		return optional.NoneBool, nil
	}

	valBool, err := strconv.ParseBool(valStr)
	if err != nil {
		return optional.NoneBool, fmt.Errorf("некорректный необязательный bool параметр %s: %w", param, err)
	}

	return optional.Boole(valBool, true), nil
}
