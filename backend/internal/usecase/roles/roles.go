package roles

import . "github.com/saime-0/nice-pea-chat/internal/model/role"

type Params struct {
	IDs  []uint `json:"ids"`
	Name string `json:"name"`
}

func (p Params) Run() ([]Role, error) {

	return []Role{
		{
			ID:          1,
			Name:        "asd",
			Permissions: nil,
		},
	}, nil
}
