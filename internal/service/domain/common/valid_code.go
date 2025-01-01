package common

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dysodeng/app/internal/pkg/helper"
	"github.com/dysodeng/app/internal/pkg/redis"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/dysodeng/app/internal/service"
	commonDo "github.com/dysodeng/app/internal/service/do/common"
	"github.com/pkg/errors"
)

// ValidCodeDomainService 验证码领域服务
type ValidCodeDomainService interface {
	SendValidCode(ctx context.Context, sender commonDo.SenderType, bizType, account string) error
	VerifyValidCode(ctx context.Context, sender commonDo.SenderType, bizType, account, code string) error
}

// validCodeDomainService 验证码领域服务
type validCodeDomainService struct {
	baseTraceSpanName string
	smsDomainService  SmsDomainService
	mailDomainService MailDomainService
}

var validCodeDomainServiceInstance ValidCodeDomainService

func NewValidCodeDomainService(smsDomainService SmsDomainService, mailDomainService MailDomainService) ValidCodeDomainService {
	if validCodeDomainServiceInstance == nil {
		validCodeDomainServiceInstance = &validCodeDomainService{
			baseTraceSpanName: "service.domain.common.ValidCodeDomainService",
			smsDomainService:  smsDomainService,
			mailDomainService: mailDomainService,
		}
	}
	return validCodeDomainServiceInstance
}

// SendValidCode 发送验证码
func (vc *validCodeDomainService) SendValidCode(ctx context.Context, sender commonDo.SenderType, bizType, account string) error {
	spanCtx, span := trace.Tracer().Start(ctx, vc.baseTraceSpanName+".SendValidCode")
	defer span.End()

	switch sender {
	case commonDo.SmsSender:
		if account == "" {
			return errors.New("缺少手机号")
		}
	case commonDo.EmailSender:
		if account == "" {
			return errors.New("缺少邮箱地址")
		}
	default:
		return errors.New("消息发送类型错误")
	}
	if bizType == "" {
		return errors.New("缺少业务类型")
	}

	client := redis.Client()

	// 验证码速率限制key
	limitKey := redis.Key(fmt.Sprintf("%s_code_limit:%s:%s", sender, bizType, account))
	var limitTotal int = 0
	var limitExpire float64 = 3600
	if client.Exists(spanCtx, limitKey).Val() > 0 {
		limitTotal, _ = client.Get(spanCtx, limitKey).Int()
		if limitTotal >= 5 {
			return service.EMValidCodeLimitError
		}
		ttl := client.TTL(spanCtx, limitKey).Val()
		limitExpire = ttl.Seconds()
	}

	templateParam := make(map[string]string)
	var expire int64 = 10
	templateParam["code"] = helper.RandomNumberString(6) // 验证码
	templateParam["time"] = strconv.FormatInt(expire, 10)

	codeCacheKey := redis.Key(fmt.Sprintf("%s_code_%s:%s", sender, bizType, account))

	smsCode := commonDo.ValidCode{
		Code:   templateParam["code"],
		Time:   time.Now().Unix(),
		Expire: 10,
	}

	client.HMSet(spanCtx, codeCacheKey, map[string]interface{}{
		"Code":   smsCode.Code,
		"Time":   smsCode.Time,
		"Expire": smsCode.Expire,
	})
	client.Expire(spanCtx, codeCacheKey, 600*time.Second)

	// 发送验证码
	var err error
	if sender == commonDo.SmsSender {
		err = vc.smsDomainService.SendSms(spanCtx, account, "code", templateParam)
	} else {
		err = vc.mailDomainService.SendMail(spanCtx, []string{account}, "验证码", "code", templateParam)
	}
	if err != nil {
		return err
	}

	// 设置发送次数
	limitTotal += 1
	client.Set(spanCtx, limitKey, limitTotal, time.Duration(limitExpire)*time.Second)

	return nil
}

// VerifyValidCode 验证码验证
func (vc *validCodeDomainService) VerifyValidCode(ctx context.Context, sender commonDo.SenderType, bizType, account, code string) error {
	spanCtr, span := trace.Tracer().Start(ctx, vc.baseTraceSpanName+".VerifyValidCode")
	defer span.End()

	switch sender {
	case commonDo.SmsSender:
		if account == "" {
			return errors.New("缺少手机号")
		}
	case commonDo.EmailSender:
		if account == "" {
			return errors.New("缺少邮箱地址")
		}
	default:
		return errors.New("消息发送类型错误")
	}
	if bizType == "" {
		return errors.New("缺少业务类型")
	}
	if code == "" {
		return errors.New("缺少验证码")
	}

	codeCacheKey := redis.Key(fmt.Sprintf("%s_code_%s:%s", sender, bizType, account))

	client := redis.Client()
	cacheCode, err := client.HGet(spanCtr, codeCacheKey, "Code").Result()
	if err != nil {
		return service.EMValidCodeExpireError
	}
	expire, _ := client.HGet(spanCtr, codeCacheKey, "Expire").Result()
	codeTime, _ := client.HGet(spanCtr, codeCacheKey, "Time").Result()
	expireInt, err := strconv.ParseInt(expire, 10, 64)
	if err != nil {
		return service.EMValidCodeExpireError
	}
	codeTimeInt, err := strconv.ParseInt(codeTime, 10, 64)
	if err != nil {
		return service.EMValidCodeExpireError
	}

	if codeTimeInt+expireInt*60 > time.Now().Unix() {
		if code != cacheCode {
			return service.EMValidCodeError
		}
		client.Del(spanCtr, codeCacheKey)
	} else {
		return service.EMValidCodeExpireError
	}

	return nil
}
