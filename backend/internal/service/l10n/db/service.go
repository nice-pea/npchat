package db

import (
	"bytes"
	"errors"
	"html/template"
	"strings"

	"gorm.io/gorm"

	"github.com/saime-0/nice-pea-chat/internal/service/l10n/db/model"
)

type Service struct {
	DB *gorm.DB
}

func (s *Service) Localize(code, locale string, vars map[string]string) (string, error) {
	fields := strings.Split(code, ":")
	if len(fields) != 2 {
		return "", errors.New("invalid code")
	}

	cond := model.Localization{Category: fields[0], Code: fields[1]}
	var locs []model.Localization
	if err := s.DB.Find(&locs, cond).Error; err != nil {
		return "", err
	} else if len(locs) == 0 {
		return "", errors.New("invalid code")
	}

	var text string
	for _, l := range locs {
		if l.Locale == locale {
			text = l.Text
			break
		} else if l.Locale == "en_US" {
			text = l.Text
		}
	}

	tpl, err := template.New("").Parse(text)
	if err != nil {
		return "", err
	}

	var result bytes.Buffer
	if err = tpl.Execute(&result, vars); err != nil {
		return "", err
	}

	return result.String(), nil
}
