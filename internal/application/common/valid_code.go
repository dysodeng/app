package common

import (
	"context"
	"github.com/dysodeng/app/internal/domain/common/model"
	"github.com/dysodeng/app/internal/domain/common/service"
	"github.com/dysodeng/app/internal/pkg/helper"
	"github.com/dysodeng/app/internal/pkg/logger"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/dysodeng/app/internal/pkg/validator"
	"github.com/pkg/errors"
)

// ValidCodeApplicationService 验证码应用服务
type ValidCodeApplicationService interface {
	SendValidCode(ctx context.Context, sender, bizType, account string) error
	VerifyValidCode(ctx context.Context, sender, bizType, account, code string) error
}

// validCodeApplicationService 验证码应用服务
type validCodeApplicationService struct {
	baseTraceSpanName      string
	validCodeDomainService service.ValidCodeDomainService
}

func NewValidCodeAppService(validCodeDomainService service.ValidCodeDomainService) ValidCodeApplicationService {
	return &validCodeApplicationService{
		baseTraceSpanName:      "application.common.ValidCodeApplicationService",
		validCodeDomainService: validCodeDomainService,
	}
}

func (svc *validCodeApplicationService) SendValidCode(ctx context.Context, sender, bizType, account string) error {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".SendValidCode")
	defer span.End()

	if !helper.Contain([]string{"sms", "email"}, sender) {
		return errors.New("消息发送类型错误")
	}

	if bizType == "" {
		return errors.New("业务类型不能为空")
	}
	if sender == "sms" {
		if account == "" {
			return errors.New("手机号不能为空")
		} else {
			if !validator.IsMobile(account) {
				return errors.New("手机号格式错误")
			}
		}
	} else {
		if account == "" {
			return errors.New("邮箱不能为空")
		} else {
			if !validator.IsEmail(account) {
				return errors.New("邮箱格式错误")
			}
		}
	}

	err := svc.validCodeDomainService.SendValidCode(spanCtx, model.SenderType(sender), bizType, account)
	if err != nil {
		logger.Error(spanCtx, "验证码发送失败", logger.ErrorField(err))
		return errors.New("验证码发送失败")
	}

	return nil
}

func (svc *validCodeApplicationService) VerifyValidCode(ctx context.Context, sender, bizType, account, code string) error {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".VerifyValidCode")
	defer span.End()

	if !helper.Contain([]string{"sms", "email"}, sender) {
		return errors.New("消息发送类型错误")
	}

	if bizType == "" {
		return errors.New("业务类型不能为空")
	}
	if account == "" {
		if sender == "sms" {
			return errors.New("手机号不能为空")
		} else {
			return errors.New("邮箱不能为空")
		}
	}
	if code == "" {
		return errors.New("验证码不能为空")
	}

	err := svc.validCodeDomainService.VerifyValidCode(spanCtx, model.SenderType(sender), bizType, account, code)
	if err != nil {
		return err
	}

	return nil
}
