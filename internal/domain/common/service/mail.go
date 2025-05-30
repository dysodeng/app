package service

import (
	"context"
	"github.com/dysodeng/app/internal/domain/common/model"
	"github.com/dysodeng/app/internal/domain/common/repository"
	"github.com/dysodeng/app/internal/pkg/mail"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/pkg/errors"
)

// MailDomainService 邮件领域服务
type MailDomainService interface {
	Config(ctx context.Context) (*model.MailConfig, error)
	SaveMailConfig(ctx context.Context, config model.MailConfig) error
	SendMail(ctx context.Context, email []string, subject, template string, templateParams map[string]string) error
}

type mailDomainService struct {
	baseTraceSpanName string
	mailRepo          repository.MailRepository
}

func NewMailDomainService(mailRepo repository.MailRepository) MailDomainService {
	return &mailDomainService{
		baseTraceSpanName: "domain.common.service.MailDomainService",
		mailRepo:          mailRepo,
	}
}

func (svc *mailDomainService) Config(ctx context.Context) (*model.MailConfig, error) {
	config, err := svc.mailRepo.Config(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "邮件配置获取失败")
	}
	return config, nil
}

func (svc *mailDomainService) SaveMailConfig(ctx context.Context, config model.MailConfig) error {
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

	err := svc.mailRepo.SaveConfig(ctx, &model.MailConfig{
		ID:        config.ID,
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

func (svc *mailDomainService) SendMail(ctx context.Context, email []string, subject, template string, templateParams map[string]string) error {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".SendMail")
	defer span.End()

	config, err := svc.mailRepo.Config(spanCtx)
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
