package service

import (
	"context"
	"github.com/dysodeng/app/internal/domain/common/model"
	"github.com/dysodeng/app/internal/domain/common/repository"
	"github.com/dysodeng/app/internal/pkg/logger"
	"github.com/dysodeng/app/internal/pkg/sms"
	"github.com/dysodeng/app/internal/pkg/sms/alicloud"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/pkg/errors"
)

// SmsDomainService 短信领域服务
type SmsDomainService interface {
	Config(ctx context.Context) (*model.SmsConfig, error)
	SaveSmsConfig(ctx context.Context, config *model.SmsConfig) error
	Template(ctx context.Context, template string) (*model.SmsTemplate, error)
	SendSms(ctx context.Context, telephone, template string, templateParams map[string]string) error
}

// smsDomainService 短信领域服务
type smsDomainService struct {
	baseTraceSpanName string
	smsRepo           repository.SmsRepository
}

func NewSmsDomainService(smsRepo repository.SmsRepository) SmsDomainService {
	return &smsDomainService{
		baseTraceSpanName: "domain.common.service.SmsDomainService",
		smsRepo:           smsRepo,
	}
}

func (svc *smsDomainService) Config(ctx context.Context) (*model.SmsConfig, error) {
	config, err := svc.smsRepo.Config(ctx)
	if err != nil {
		logger.Error(ctx, "短信配置获取失败", logger.ErrorField(err))
		return nil, errors.Wrap(err, "短信配置获取失败")
	}
	return config, nil
}

func (svc *smsDomainService) SaveSmsConfig(ctx context.Context, config *model.SmsConfig) error {
	if config.SmsType == "" {
		return errors.New("短信类型不能为空")
	}
	if config.SmsType == "ali_cloud" {
		if config.AppKey == "" || config.SecretKey == "" || config.FreeSignName == "" {
			return errors.New("阿里云短信配置信息不完整")
		}
	} else {
		return errors.New("暂不支持该短信类型")
	}

	err := svc.smsRepo.SaveConfig(ctx, &model.SmsConfig{
		ID:              config.ID,
		AppKey:          config.AppKey,
		FreeSignName:    config.FreeSignName,
		SecretKey:       config.SecretKey,
		SmsType:         config.SmsType,
		ValidCodeExpire: config.ValidCodeExpire,
	})
	if err != nil {
		logger.Error(ctx, "短信配置保存失败", logger.ErrorField(err))
		return errors.Wrap(err, "短信配置保存失败")
	}
	return nil
}

func (svc *smsDomainService) Template(ctx context.Context, template string) (*model.SmsTemplate, error) {
	smsTemplate, err := svc.smsRepo.Template(ctx, template)
	if err != nil {
		logger.Error(ctx, "短信模板获取失败", logger.ErrorField(err))
		return nil, errors.Wrap(err, "短信模板获取失败")
	}
	return smsTemplate, nil
}

func (svc *smsDomainService) SendSms(ctx context.Context, telephone, template string, templateParams map[string]string) error {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".SendSms")
	defer span.End()

	config, err := svc.smsRepo.Config(spanCtx)
	if err != nil {
		trace.Error(errors.New("短信配置获取失败"), span)
		return errors.Wrap(err, "短信配置获取失败")
	}
	if config.ID <= 0 {
		trace.Error(errors.New("短信配置不存在"), span)
		return errors.New("短信配置不存在")
	}

	smsTemplate, err := svc.smsRepo.Template(ctx, template)
	if err != nil {
		trace.Error(errors.New("短信模板获取失败"), span)
		return errors.Wrap(err, "短信模板获取失败")
	}
	if smsTemplate.ID <= 0 {
		trace.Error(errors.New("短信模板不存在"), span)
		return errors.New("短信模板不存在")
	}

	var sender sms.Sender
	switch config.SmsType {
	case "ali_cloud":
		sender = alicloud.NewAliCloudSmsSender(
			telephone,
			smsTemplate.TemplateId,
			alicloud.WithConfig(
				config.AppKey,
				config.SecretKey,
				config.FreeSignName,
			),
			alicloud.WithParams(templateParams),
		)
	default:
		return errors.New("暂不支持该短信类型")
	}

	_, err = sender.SendSms()
	if err != nil {
		trace.Error(errors.New("短信发送失败"), span)
		return errors.Wrap(err, "短信发送失败")
	}

	return nil
}
