package mail

import (
	"bytes"
	htmlTemplate "html/template"

	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
)

// Sender 邮件发送器接口
type Sender interface {
	// SendMail 发送邮件
	SendMail() error
}

// Mail 邮件发送器
type Mail struct {
	message *gomail.Message
	option  *option
}

// NewMailSender 新建邮件发送器
// to []string 接收者邮箱地址列表
// template 邮件模板
// config Config 邮件配置
// opts []Option 邮件额外选项
func NewMailSender(to []string, template string, config Config, opts ...Option) (*Mail, error) {
	o := &option{}
	config(o)
	for _, opt := range opts {
		opt(o)
	}

	mailSender := &Mail{
		option: o,
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(o.fromUser, o.fromName))
	m.SetHeader("To", to...)
	if o.subject != "" {
		m.SetHeader("Subject", o.subject)
	}

	// 邮件内容
	tpl, err := htmlTemplate.ParseFiles("template/email/" + template + ".html")
	if err != nil {
		return nil, errors.Wrap(err, "邮件模板不存在")
	}

	var body bytes.Buffer
	err = tpl.Execute(&body, o.params)
	if err != nil {
		return nil, errors.Wrap(err, "邮件模板解析错误")
	}
	m.SetBody("text/html", body.String())

	// 附件
	if len(mailSender.option.attach) > 0 {
		for _, a := range mailSender.option.attach {
			m.Attach(a.filename, gomail.Rename(a.attachName))
		}
	}

	mailSender.message = m

	return mailSender, nil
}

// SendMail 发送邮件
func (sender *Mail) SendMail() error {
	d := gomail.NewDialer(sender.option.host, sender.option.port, sender.option.username, sender.option.password)
	err := d.DialAndSend(sender.message)
	if err != nil {
		return errors.Wrap(err, "邮件发送失败")
	}
	return nil
}
