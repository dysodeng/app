package common

type SmsConfig struct {
	SmsType         string `json:"sms_type"`
	AppKey          string `json:"app_key"`
	SecretKey       string `json:"secret_key"`
	FreeSignName    string `json:"free_sign_name"`
	ValidCodeExpire uint   `json:"valid_code_expire"`
}

type SmsTemplate struct {
	TemplateName string `json:"template_name"`
	Template     string `json:"template"`
	TemplateId   string `json:"template_id"`
}
