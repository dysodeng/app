package common

import (
	"context"

	commonDao "github.com/dysodeng/app/internal/dal/dao/common"
	"github.com/dysodeng/app/internal/dal/model/common"
	"github.com/dysodeng/app/internal/pkg/logger"
	"github.com/dysodeng/app/internal/pkg/sms"
	"github.com/dysodeng/app/internal/pkg/sms/alicloud"
	"github.com/dysodeng/app/internal/pkg/trace"
	commonDo "github.com/dysodeng/app/internal/service/do/common"
	"github.com/pkg/errors"
)

// SmsDomainService 短信领域服务
type SmsDomainService struct {
	ctx               context.Context
	smsDao            *commonDao.SmsDao
	baseTraceSpanName string
}

func NewSmsDomainService(ctx context.Context) *SmsDomainService {
	baseTraceSpanName := "domain.common.SmsDomainService"
	traceCtx := trace.New().NewSpan(ctx, baseTraceSpanName)
	return &SmsDomainService{
		ctx:               traceCtx,
		smsDao:            commonDao.NewSmsDao(traceCtx),
		baseTraceSpanName: baseTraceSpanName,
	}
}

// Config 获取短信配置
func (ss *SmsDomainService) Config() (*commonDo.SmsConfig, error) {
	config, err := ss.smsDao.Config()
	if err != nil {
		logger.Error(ss.ctx, "短信配置获取失败", logger.ErrorField(err))
		return nil, errors.Wrap(err, "短信配置获取失败")
	}
	return &commonDo.SmsConfig{
		AppKey:          config.AppKey,
		FreeSignName:    config.FreeSignName,
		SecretKey:       config.SecretKey,
		SmsType:         config.SmsType,
		ValidCodeExpire: config.ValidCodeExpire,
	}, nil
}

// SaveSmsConfig 保存短信配置
func (ss *SmsDomainService) SaveSmsConfig(config commonDo.SmsConfig) error {
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
	err := ss.smsDao.SaveConfig(common.SmsConfig{
		AppKey:          config.AppKey,
		FreeSignName:    config.FreeSignName,
		SecretKey:       config.SecretKey,
		SmsType:         config.SmsType,
		ValidCodeExpire: config.ValidCodeExpire,
	})
	if err != nil {
		logger.Error(ss.ctx, "短信配置保存失败", logger.ErrorField(err))
		return errors.Wrap(err, "短信配置保存失败")
	}
	return nil
}

// Template 获取短信模板
func (ss *SmsDomainService) Template(template string) (*commonDo.SmsTemplate, error) {
	smsTemplate, err := ss.smsDao.Template(template)
	if err != nil {
		logger.Error(ss.ctx, "短信模板获取失败", logger.ErrorField(err))
		return nil, errors.Wrap(err, "短信模板获取失败")
	}
	return &commonDo.SmsTemplate{
		Template:     smsTemplate.Template,
		TemplateId:   smsTemplate.TemplateId,
		TemplateName: smsTemplate.TemplateName,
	}, nil
}

// SendSms 发送短信
func (ss *SmsDomainService) SendSms(telephone, template string, templateParams map[string]string) error {
	config, err := ss.smsDao.Config()
	if err != nil {
		return errors.Wrap(err, "短信配置获取失败")
	}
	if config.ID <= 0 {
		return errors.New("短信配置不存在")
	}

	smsTemplate, err := ss.smsDao.Template(template)
	if err != nil {
		return errors.Wrap(err, "短信模板获取失败")
	}
	if smsTemplate.ID <= 0 {
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
		return errors.Wrap(err, "短信发送失败")
	}

	return nil
}
