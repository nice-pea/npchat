package db

import (
	"bytes"
	"errors"
	"html/template"
	"strings"

	"gorm.io/gorm"

	"github.com/saime-0/nice-pea-chat/internal/service/l10n"
	"github.com/saime-0/nice-pea-chat/internal/service/l10n/db/model"
)

type Service struct {
	DB *gorm.DB
}

func (s *Service) Localize(code, locale string, vars map[string]string) (string, error) {
	codeParts := strings.Split(code, ":")
	if len(codeParts) != 2 {
		return "", errors.New("invalid code")
	}

	cond := model.Localization{Category: codeParts[0], Item: codeParts[1]}
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
		} else if l.Locale == l10n.LocaleDefault {
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
