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
type MailDomainService interface {
	Config(ctx context.Context) (*commonDo.MailConfig, error)
	SaveMailConfig(ctx context.Context, config commonDo.MailConfig) error
	SendMail(ctx context.Context, email []string, subject, template string, templateParams map[string]string) error
}

// mailDomainService 邮件领域服务
type mailDomainService struct {
	baseTraceSpanName string
	mailDao           commonDao.MailDao
}

var mailDomainServiceInstance MailDomainService

func NewMailDomainService(mailDao commonDao.MailDao) MailDomainService {
	if mailDomainServiceInstance == nil {
		mailDomainServiceInstance = &mailDomainService{
			baseTraceSpanName: "service.domain.common.MailDomainService",
			mailDao:           mailDao,
		}
	}
	return mailDomainServiceInstance
}

// Config 获取邮件配置
func (ms *mailDomainService) Config(ctx context.Context) (*commonDo.MailConfig, error) {
	config, err := ms.mailDao.Config(ctx)
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
func (ms *mailDomainService) SaveMailConfig(ctx context.Context, config commonDo.MailConfig) error {
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

	err := ms.mailDao.SaveConfig(ctx, common.MailConfig{
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
func (ms *mailDomainService) SendMail(ctx context.Context, email []string, subject, template string, templateParams map[string]string) error {
	spanCtx, span := trace.Tracer().Start(ctx, ms.baseTraceSpanName+".SendMail")
	defer span.End()

	config, err := ms.mailDao.Config(spanCtx)
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
