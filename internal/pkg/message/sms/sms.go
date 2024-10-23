package sms

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/dysodeng/app/internal/dal/model/common"

	"github.com/dysodeng/app/internal/pkg/helper"

	"github.com/dysodeng/app/internal/pkg/api"
	"github.com/dysodeng/app/internal/pkg/db"
	"github.com/dysodeng/app/internal/pkg/message"
	"github.com/dysodeng/app/internal/pkg/message/sms/alicloud"
	"github.com/dysodeng/app/internal/pkg/redis"

	"github.com/pkg/errors"
)

// SendSmsCode 发送验证码
func SendSmsCode(phoneNumber, template string) error {
	var smsConfig common.SmsConfig
	db.DB().First(&smsConfig)
	if smsConfig.ID <= 0 {
		log.Println("SendSmsCode: 短信未配置")
		return errors.New("短信未配置")
	}

	var dbTemplate = "code"

	var templateConfig common.SmsTemplate
	db.DB().Where("template=?", dbTemplate).First(&templateConfig)
	if templateConfig.ID <= 0 {
		return errors.New("短信模版不存在")
	}

	templateId := templateConfig.TemplateId
	templateParam := make(map[string]string)

	templateParam["code"] = helper.GenValidateCode(6) // 验证码

	// 验证码缓存
	key := redis.Key("sms_code_" + template + ":" + phoneNumber)

	redisClient := redis.Client()

	redisClient.Del(context.Background(), key)

	smsCode := message.Code{
		Code:   templateParam["code"],
		Time:   time.Now().Unix(),
		Expire: int64(smsConfig.ValidCodeExpire),
	}

	redisClient.HMSet(context.Background(), key, map[string]interface{}{"Code": smsCode.Code, "Time": smsCode.Time, "Expire": smsCode.Expire})
	redisClient.Expire(context.Background(), key, time.Duration(smsConfig.ValidCodeExpire*60)*time.Second)

	log.Printf("send sms code:[template:%s, telephone:%s, code:%s]\n", template, phoneNumber, smsCode.Code)

	var sender Sender

	switch smsConfig.SmsType {
	case "ali_cloud":
		sender = alicloud.NewAliCloudSmsSender(
			phoneNumber,
			templateId,
			alicloud.WithConfig(
				smsConfig.AppKey,
				smsConfig.SecretKey,
				smsConfig.FreeSignName,
			),
			alicloud.WithParams(templateParam),
		)
		break
	default:
		return errors.New("sms sender error:" + smsConfig.SmsType)
	}

	_, err := sender.SendSms()
	if err != nil {
		return err
	}

	return nil
}

// ValidSmsCode 验证短信验证码
func ValidSmsCode(phoneNumber, template, smsCode string) error {
	// 验证码缓存
	key := redis.Key("sms_code_" + template + ":" + phoneNumber)

	redisClient := redis.Client()

	code, err := redisClient.HGet(context.Background(), key, "Code").Result()
	if err != nil {
		log.Println(err)
		return api.EMValidCodeExpireError
	}
	expire, _ := redisClient.HGet(context.Background(), key, "Expire").Result()
	codeTime, _ := redisClient.HGet(context.Background(), key, "Time").Result()

	expireInt, err := strconv.ParseInt(expire, 10, 64)
	if err != nil {
		log.Println(err)
		return api.EMValidCodeExpireError
	}
	codeTimeInt, err := strconv.ParseInt(codeTime, 10, 64)
	if err != nil {
		log.Println(err)
		return api.EMValidCodeExpireError
	}

	if codeTimeInt+expireInt*60 > time.Now().Unix() {
		if code != smsCode {
			return api.EMValidCodeError
		}

		redisClient.Del(context.Background(), key)
		// 验证码速率限制key
		limitKey := redis.Key("sms_code_limit_" + template + ":" + phoneNumber)
		redisClient.Del(context.Background(), limitKey)
	} else {
		return api.EMValidCodeExpireError
	}

	return nil
}

// SendSmsMessage 发送普通短信消息
func SendSmsMessage(phoneNumber, template string, params map[string]string) error {
	var smsConfig common.SmsConfig
	db.DB().First(&smsConfig)
	if smsConfig.ID <= 0 {
		log.Println("SendSmsCode: 短信未配置")
		return errors.New("短信未配置")
	}

	var templateConfig common.SmsTemplate
	db.DB().Where("template=?", template).First(&templateConfig)
	if templateConfig.ID <= 0 {
		return errors.New("短信模版不存在")
	}

	templateId := templateConfig.TemplateId

	var sender Sender

	switch smsConfig.SmsType {
	case "ali_cloud":
		sender = alicloud.NewAliCloudSmsSender(
			phoneNumber,
			templateId,
			alicloud.WithConfig(
				smsConfig.AppKey,
				smsConfig.SecretKey,
				smsConfig.FreeSignName,
			),
			alicloud.WithParams(params),
		)
		break
	default:
		return errors.New("sms sender error:" + smsConfig.SmsType)
	}

	_, err := sender.SendSms()
	if err != nil {
		return err
	}

	return nil
}
