package permissions

import (
	"github.com/saime-0/nice-pea-chat/internal/model"
	permsL10n "github.com/saime-0/nice-pea-chat/internal/model/permission/l10n"
	"github.com/saime-0/nice-pea-chat/internal/service/l10n"
)

type Params struct {
	Locale string
	L10n   l10n.Service
}

func (p Params) Run() (_ []model.Permission, err error) {
	perms := make([]model.Permission, len(model.Permissions))
	for i, id := range model.Permissions {
		perms[i] = model.Permission{ID: id}
		perms[i].Name, err = p.L10n.Localize(permsL10n.Name(id), p.Locale, nil)
		if err != nil {
			return nil, err
		}
		perms[i].Desc, err = p.L10n.Localize(permsL10n.Desc(id), p.Locale, nil)
		if err != nil {
			return nil, err
		}
	}

	return perms, nil
}
