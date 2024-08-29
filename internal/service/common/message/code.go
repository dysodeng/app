package message

import (
	"context"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/pkg/api"
	"github.com/dysodeng/app/internal/pkg/message/mail"
	"github.com/dysodeng/app/internal/pkg/message/sms"
	"github.com/dysodeng/app/internal/pkg/redis"
	"github.com/dysodeng/app/internal/service"

	"time"

	messageModel "github.com/dysodeng/app/internal/model/common"

	"github.com/pkg/errors"
)

// CodeMessageService 验证码消息服务
type CodeMessageService struct{}

// NewCodeMessageService 创建验证码消息服务
func NewCodeMessageService() *CodeMessageService {
	return &CodeMessageService{}
}

// SendValidCode 发送验证码消息
// @param senderType sms-短信 email-邮件
// @param bizType 自定义业务类型
// @param userType 用户类型 user-终端用户
// @param account 用户账号 senderType=sms时为手机号 senderType=email时为邮箱地址
// @param userId 用户ID 对应于userType的用户ID
func (ms CodeMessageService) SendValidCode(
	senderType messageModel.SenderType,
	bizType,
	userType,
	account string,
) service.Error {
	switch senderType {
	case messageModel.SmsSender:
		if account == "" {
			return service.Error{Code: api.CodeFail, Error: errors.New("缺少手机号")}
		}
		break
	case messageModel.EmailSender:
		if account == "" {
			return service.Error{Code: api.CodeFail, Error: errors.New("缺少邮箱地址")}
		}
		break
	default:
		return service.Error{Code: api.CodeFail, Error: errors.New("消息发送类型错误")}
	}
	if bizType == "" {
		return service.Error{Code: api.CodeFail, Error: errors.New("缺少业务类型")}
	}

	// 验证码速率限制key
	limitKey := redis.Key("sms_code_limit:" + userType + ":" + account + ":" + bizType)
	client := redis.Client()
	var total int = 0
	var expire float64 = 3600
	if client.Exists(context.Background(), limitKey).Val() > 0 {
		total, _ = client.Get(context.Background(), limitKey).Int()
		if total >= 5 {
			return service.Error{Code: api.CodeFail, Error: api.EMValidCodeLimitError}
		}
		ttl := client.TTL(context.Background(), limitKey).Val()
		expire = ttl.Seconds()
	}

	if senderType == messageModel.SmsSender {
		err := sms.SendSmsCode(account, bizType)
		if err != nil {
			return service.Error{Code: api.CodeFail, Error: errors.Wrap(err, "验证码发送失败，请稍候再试。")}
		}
	} else {
		err := mail.SendMailCode(account, bizType)
		if err != nil {
			return service.Error{Code: api.CodeFail, Error: errors.Wrap(err, "验证码发送失败，请稍候再试。")}
		}
	}

	// 设置发送次数
	total += 1
	client.Set(context.Background(), limitKey, total, time.Duration(expire)*time.Second)

	return service.Error{Code: api.CodeOk}
}

// VerifyValidCode 验证码验证
// @param senderType sms-短信消息 email-邮件消息
// @param bizType 自定义业务类型
// @param account 用户账号 senderType=sms时为手机号 senderType=email时为邮箱地址
// @param code 验证码
func (ms CodeMessageService) VerifyValidCode(
	senderType messageModel.SenderType,
	bizType,
	account,
	code string,
) service.Error {
	// 非生产环境下可以使用万能验证码
	if config.App.Env != config.Release {
		if code == "8848" {
			return service.Error{Code: api.CodeOk}
		}
	}

	switch senderType {
	case messageModel.SmsSender:
		if account == "" {
			return service.Error{Code: api.CodeFail, Error: errors.New("手机号为空")}
		}
		err := sms.ValidSmsCode(account, bizType, code)
		if err != nil {
			return service.Error{Code: api.CodeFail, Error: errors.New(err.Error())}
		}
		break

	case messageModel.EmailSender:
		if account == "" {
			return service.Error{Code: api.CodeFail, Error: errors.New("邮箱地址为空")}
		}
		err := mail.ValidMailCode(account, bizType, code)
		if err != nil {
			return service.Error{Code: api.CodeFail, Error: errors.New(err.Error())}
		}
		break

	default:
		return service.Error{Code: api.CodeFail, Error: errors.New("消息发送类型错误")}
	}

	return service.Error{Code: api.CodeOk}
}
