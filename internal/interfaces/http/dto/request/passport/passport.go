package passport

type LoginRequest struct {
	UserType string `json:"user_type" binding:"required" msg:"缺少用户类型"`

	// 用户登录
	GrantType string `json:"grant_type"`
	WxCode    string `json:"wx_code"`
	Code      string `json:"code"`
	OpenId    string `json:"open_id"`

	// 管理员登录
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" msg:"缺少refresh_token"`
}
