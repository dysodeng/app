package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dysodeng/app/internal/domain/common/model"
	"github.com/dysodeng/app/internal/pkg/helper"
	"github.com/dysodeng/app/internal/pkg/redis"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/pkg/errors"
)

// ValidCodeDomainService 验证码领域服务
type ValidCodeDomainService interface {
	SendValidCode(ctx context.Context, sender model.SenderType, bizType, account string) error
	VerifyValidCode(ctx context.Context, sender model.SenderType, bizType, account, code string) error
}

// validCodeDomainService 验证码领域服务
type validCodeDomainService struct {
	baseTraceSpanName string
	smsDomainService  SmsDomainService
	mailDomainService MailDomainService
}

func NewValidCodeDomainService(smsDomainService SmsDomainService, mailDomainService MailDomainService) ValidCodeDomainService {
	return &validCodeDomainService{
		baseTraceSpanName: "domain.common.service.ValidCodeDomainService",
		smsDomainService:  smsDomainService,
		mailDomainService: mailDomainService,
	}
}

func (svc *validCodeDomainService) SendValidCode(ctx context.Context, sender model.SenderType, bizType, account string) error {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".SendValidCode")
	defer span.End()

	switch sender {
	case model.SmsSender:
		if account == "" {
			return errors.New("缺少手机号")
		}
	case model.EmailSender:
		if account == "" {
			return errors.New("缺少邮箱地址")
		}
	default:
		return errors.New("消息发送类型错误")
	}
	if bizType == "" {
		return errors.New("缺少业务类型")
	}

	client := redis.MainClient()

	// 验证码速率限制key
	limitKey := redis.MainKey(fmt.Sprintf("%s_code_limit:%s:%s", sender, bizType, account))
	var limitTotal int = 0
	var limitExpire float64 = 3600
	if client.Exists(spanCtx, limitKey).Val() > 0 {
		limitTotal, _ = client.Get(spanCtx, limitKey).Int()
		if limitTotal >= 5 {
			return errors.New("操作太频繁，请稍后再尝试")
		}
		ttl := client.TTL(spanCtx, limitKey).Val()
		limitExpire = ttl.Seconds()
	}

	templateParam := make(map[string]string)
	var expire int64 = 10
	templateParam["code"] = helper.RandomNumberString(6) // 验证码
	templateParam["time"] = strconv.FormatInt(expire, 10)

	codeCacheKey := redis.MainKey(fmt.Sprintf("%s_code_%s:%s", sender, bizType, account))

	smsCode := model.ValidCode{
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
	if sender == model.SmsSender {
		err = svc.smsDomainService.SendSms(spanCtx, account, "code", templateParam)
	} else {
		err = svc.mailDomainService.SendMail(spanCtx, []string{account}, "验证码", "code", templateParam)
	}
	if err != nil {
		return err
	}

	// 设置发送次数
	limitTotal += 1
	client.Set(spanCtx, limitKey, limitTotal, time.Duration(limitExpire)*time.Second)

	return nil
}

func (svc *validCodeDomainService) VerifyValidCode(ctx context.Context, sender model.SenderType, bizType, account, code string) error {
	spanCtr, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".VerifyValidCode")
	defer span.End()

	switch sender {
	case model.SmsSender:
		if account == "" {
			return errors.New("缺少手机号")
		}
	case model.EmailSender:
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

	codeCacheKey := redis.MainKey(fmt.Sprintf("%s_code_%s:%s", sender, bizType, account))

	client := redis.MainClient()
	cacheCode, err := client.HGet(spanCtr, codeCacheKey, "Code").Result()
	if err != nil {
		return errors.New("验证码已过期")
	}
	expire, _ := client.HGet(spanCtr, codeCacheKey, "Expire").Result()
	codeTime, _ := client.HGet(spanCtr, codeCacheKey, "Time").Result()
	expireInt, err := strconv.ParseInt(expire, 10, 64)
	if err != nil {
		return errors.New("验证码已过期")
	}
	codeTimeInt, err := strconv.ParseInt(codeTime, 10, 64)
	if err != nil {
		return errors.New("验证码已过期")
	}

	if codeTimeInt+expireInt*60 > time.Now().Unix() {
		if code != cacheCode {
			return errors.New("验证码错误")
		}
		client.Del(spanCtr, codeCacheKey)
	} else {
		return errors.New("验证码已过期")
	}

	return nil
}
