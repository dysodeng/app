package common

type SendValidCodeBody struct {
	Type      string `json:"type"`
	BizType   string `json:"biz_type"`
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
}

type VerifyValidCodeBody struct {
	Type      string `json:"type"`
	BizType   string `json:"biz_type"`
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
	ValidCode string `json:"valid_code"`
}
