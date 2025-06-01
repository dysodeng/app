package query

import "github.com/dysodeng/app/internal/pkg/form"

type UserListQuery struct {
	form.Pagination
	Telephone string `json:"telephone"`
	Nickname  string `json:"nickname"`
	RealName  string `json:"real_name"`
	Status    uint8  `json:"status"`
}
