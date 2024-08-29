package mail

// attach 附件
type attach struct {
	filename   string
	attachName string
}

// option 邮件配置
type option struct {
	transport string
	username  string
	password  string
	host      string
	port      int
	fromUser  string
	fromName  string
	subject   string
	attach    []attach
	params    map[string]string
}

// Config 邮件主要配置
type Config func(*option)

// WithConfig 设置邮件主要配置
// @param host string 邮件提供商发送地址
// @param port int 邮件提供商发送端口
// @param transport string 邮件提供商发送协议
// @param username string 邮件提供商用户名
// @param password string 邮件提供商用户密码
// @param fromUser string 发送者邮箱地址
// @param fromName string 发送者名称
func WithConfig(host string, port int, transport, username, password, fromUser, fromName string) Config {
	return func(o *option) {
		o.host = host
		o.port = port
		o.transport = transport
		o.username = username
		o.password = password
		o.fromUser = fromUser
		o.fromName = fromName
	}
}

// Option 邮件额外配置
type Option func(*option)

// WithSubject 添加邮件主题
// @param subject string 邮件主题
func WithSubject(subject string) Option {
	return func(o *option) {
		o.subject = subject
	}
}

// WithParams 添加邮件参数
// @param params map[string]string 参数列表
func WithParams(params map[string]string) Option {
	return func(o *option) {
		o.params = params
	}
}

// WithAttach 添加邮件附件
// @param filename string 附件文件地址
// @param attachName string 附件名称
func WithAttach(filename, attachName string) Option {
	return func(o *option) {
		o.attach = append(o.attach, attach{
			filename:   filename,
			attachName: attachName,
		})
	}
}
