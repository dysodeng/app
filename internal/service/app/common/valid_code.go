package common

import (
	"context"

	"github.com/dysodeng/app/internal/pkg/helper"
	"github.com/dysodeng/app/internal/pkg/logger"
	"github.com/dysodeng/app/internal/pkg/trace"
	"github.com/dysodeng/app/internal/pkg/validator"
	commonDo "github.com/dysodeng/app/internal/service/do/common"
	"github.com/dysodeng/app/internal/service/domain/common"
	"github.com/pkg/errors"
)

// ValidCodeAppService 验证码应用服务
type ValidCodeAppService struct {
	ctx                    context.Context
	validCodeDomainService *common.ValidCodeDomainService
	baseTraceSpanName      string
}

func NewValidCodeAppService(ctx context.Context) *ValidCodeAppService {
	baseTraceSpanName := "app.service.common.ValidCodeAppService"
	traceCtx := trace.New().NewSpan(ctx, baseTraceSpanName)
	return &ValidCodeAppService{
		ctx:                    traceCtx,
		validCodeDomainService: common.NewValidCodeDomainService(traceCtx),
		baseTraceSpanName:      baseTraceSpanName,
	}
}

func (vc *ValidCodeAppService) SendValidCode(sender, bizType, account string) error {
	if !helper.Contain(sender, []string{"sms", "email"}) {
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

	err := vc.validCodeDomainService.SendValidCode(commonDo.SenderType(sender), bizType, account)
	if err != nil {
		logger.Error(vc.ctx, "验证码发送失败", logger.ErrorField(err))
		return errors.New("验证码发送失败")
	}

	return nil
}

func (vc *ValidCodeAppService) VerifyValidCode(sender, bizType, account, code string) error {
	if !helper.Contain(sender, []string{"sms", "email"}) {
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

	err := vc.validCodeDomainService.VerifyValidCode(commonDo.SenderType(sender), bizType, account, code)
	if err != nil {
		return err
	}

	return nil
}
