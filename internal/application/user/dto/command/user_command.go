package command

type UserCreateCommand struct {
	Telephone string `json:"telephone"`
	Password  string `json:"password"`
	Avatar    string `json:"avatar"`
	Nickname  string `json:"nickname"`
	RealName  string `json:"real_name"`
	Gender    uint8  `json:"gender"`
	Birthday  string `json:"birthday"`
	Status    uint8  `json:"status"`
}
