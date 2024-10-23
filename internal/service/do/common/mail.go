package common

// MailConfig 邮件配置领域对象
type MailConfig struct {
	User      string `json:"user"`
	FromName  string `json:"from_name"`
	Transport string `json:"transport"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
}
