package command

type LoginCommand struct {
	UserType string

	// 管理员登录
	Username string
	Password string

	// 用户登录
	PlatformType string
	GrantType    string
	WxCode       string
	Code         string
	OpenId       string
}

type VerifyTokenCommand struct {
	UserType string
	Token    string
}
