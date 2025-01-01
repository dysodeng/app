package common

import (
	"context"

	"github.com/dysodeng/app/internal/pkg/helper"
	"github.com/dysodeng/app/internal/pkg/logger"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/dysodeng/app/internal/pkg/validator"
	commonDo "github.com/dysodeng/app/internal/service/do/common"
	"github.com/dysodeng/app/internal/service/domain/common"
	"github.com/pkg/errors"
)

// ValidCodeAppService 验证码应用服务
type ValidCodeAppService interface {
	SendValidCode(ctx context.Context, sender, bizType, account string) error
	VerifyValidCode(ctx context.Context, sender, bizType, account, code string) error
}

// validCodeAppService 验证码应用服务
type validCodeAppService struct {
	baseTraceSpanName      string
	validCodeDomainService common.ValidCodeDomainService
}

var validCodeAppServiceInstance ValidCodeAppService

func NewValidCodeAppService(validCodeDomainService common.ValidCodeDomainService) ValidCodeAppService {
	if validCodeAppServiceInstance == nil {
		validCodeAppServiceInstance = &validCodeAppService{
			baseTraceSpanName:      "service.app.common.ValidCodeAppService",
			validCodeDomainService: validCodeDomainService,
		}
	}
	return validCodeAppServiceInstance
}

func (vc *validCodeAppService) SendValidCode(ctx context.Context, sender, bizType, account string) error {
	spanCtx, span := trace.Tracer().Start(ctx, vc.baseTraceSpanName+".SendValidCode")
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

	err := vc.validCodeDomainService.SendValidCode(spanCtx, commonDo.SenderType(sender), bizType, account)
	if err != nil {
		logger.Error(spanCtx, "验证码发送失败", logger.ErrorField(err))
		return errors.New("验证码发送失败")
	}

	return nil
}

func (vc *validCodeAppService) VerifyValidCode(ctx context.Context, sender, bizType, account, code string) error {
	spanCtx, span := trace.Tracer().Start(ctx, vc.baseTraceSpanName+".VerifyValidCode")
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

	err := vc.validCodeDomainService.VerifyValidCode(spanCtx, commonDo.SenderType(sender), bizType, account, code)
	if err != nil {
		return err
	}

	return nil
}
