package response

type LoginResponse struct {
	Registered         bool  `json:"registered"` // 用户是否已注册
	Token              any   `json:"token"`
	Expire             int64 `json:"expire"`
	RefreshToken       any   `json:"refresh_token"`
	RefreshTokenExpire int64 `json:"refresh_token_expire"`
	Attach             any   `json:"attach,omitempty"`
}
