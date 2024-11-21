package common

import (
	"context"

	commonDao "github.com/dysodeng/app/internal/dal/dao/common"
	"github.com/dysodeng/app/internal/dal/model/common"
	"github.com/dysodeng/app/internal/pkg/mail"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	commonDo "github.com/dysodeng/app/internal/service/do/common"
	"github.com/pkg/errors"
)

// MailDomainService 邮件领域服务
type MailDomainService struct {
	ctx               context.Context
	baseTraceSpanName string
}

func NewMailDomainService(ctx context.Context) *MailDomainService {
	return &MailDomainService{
		ctx:               ctx,
		baseTraceSpanName: "service.domain.common.MailDomainService",
	}
}

// Config 获取邮件配置
func (ms *MailDomainService) Config() (*commonDo.MailConfig, error) {
	mailDao := commonDao.NewMailDao(ms.ctx)
	config, err := mailDao.Config()
	if err != nil {
		return nil, errors.Wrap(err, "邮件配置获取失败")
	}
	return &commonDo.MailConfig{
		Host:      config.Host,
		Port:      config.Port,
		FromName:  config.FromName,
		Password:  config.Password,
		Transport: config.Transport,
		User:      config.User,
		Username:  config.Username,
	}, nil
}

// SaveMailConfig 保存邮件配置
func (ms *MailDomainService) SaveMailConfig(config commonDo.MailConfig) error {
	if config.Host == "" {
		return errors.New("缺少邮件服务器地址")
	}
	if config.Port == 0 {
		return errors.New("缺少邮件服务器端口")
	}
	if config.User == "" {
		return errors.New("缺少发件邮箱地址")
	}
	if config.FromName == "" {
		return errors.New("缺少发送者名称")
	}
	if config.Username == "" {
		return errors.New("缺少邮箱用户名")
	}
	if config.Password == "" {
		return errors.New("缺少邮箱授权码")
	}
	if config.Transport != "smtp" {
		config.Transport = "smtp"
	}

	mailDao := commonDao.NewMailDao(ms.ctx)
	err := mailDao.SaveConfig(common.MailConfig{
		Host:      config.Host,
		Port:      config.Port,
		FromName:  config.FromName,
		Password:  config.Password,
		Transport: config.Transport,
		User:      config.User,
		Username:  config.Username,
	})
	if err != nil {
		return errors.Wrap(err, "邮件配置保存失败")
	}
	return nil
}

// SendMail 发送邮件
func (ms *MailDomainService) SendMail(email []string, subject, template string, templateParams map[string]string) error {
	spanCtx, span := trace.Tracer().Start(ms.ctx, ms.baseTraceSpanName+".SendMail")
	defer span.End()

	mailDao := commonDao.NewMailDao(spanCtx)
	config, err := mailDao.Config()
	if err != nil {
		trace.Error(errors.Wrap(err, "邮件配置获取失败"), span)
		return errors.Wrap(err, "邮件配置获取失败")
	}
	if config.ID <= 0 {
		trace.Error(errors.New("邮件配置不存在"), span)
		return errors.New("邮件配置不存在")
	}

	opts := []mail.Option{
		mail.WithParams(templateParams),
		mail.WithSubject(subject),
	}

	sender, err := mail.NewMailSender(
		email,
		template,
		mail.WithConfig(
			config.Host,
			config.Port,
			config.Transport,
			config.Username,
			config.Password,
			config.User,
			config.FromName,
		),
		opts...,
	)
	if err != nil {
		trace.Error(err, span)
		return errors.Wrap(err, "创建邮件发送器失败")
	}

	err = sender.SendMail()
	if err != nil {
		trace.Error(err, span)
		return errors.Wrap(err, "邮件发送失败")
	}

	return nil
}
