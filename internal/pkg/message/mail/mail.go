package mail

import (
	"context"
	"log"
	"strconv"
	"time"

	messageModel "github.com/dysodeng/app/internal/dal/model/common"

	"github.com/dysodeng/app/internal/pkg/helper"

	"github.com/dysodeng/app/internal/pkg/api"
	"github.com/dysodeng/app/internal/pkg/db"
	"github.com/dysodeng/app/internal/pkg/message"
	"github.com/dysodeng/app/internal/pkg/redis"

	"github.com/pkg/errors"
)

// SendMailCode 发送邮件验证码
// @param email 接收者邮箱地址
func SendMailCode(email, template string) error {

	templateParam := make(map[string]string)

	var expire int64 = 10

	templateParam["code"] = helper.GenValidateCode(6) // 验证码
	templateParam["time"] = strconv.FormatInt(expire, 10)

	// 验证码缓存
	key := redis.Key("email_code_" + template + ":" + email)

	redisClient := redis.Client()

	redisClient.Del(context.Background(), key)

	smsCode := message.Code{
		Code:   templateParam["code"],
		Time:   time.Now().Unix(),
		Expire: expire,
	}

	redisClient.HMSet(context.Background(), key, map[string]interface{}{"Code": smsCode.Code, "Time": smsCode.Time, "Expire": smsCode.Expire})
	redisClient.Expire(context.Background(), key, time.Duration(expire*60)*time.Second)

	log.Printf("send email code:[template:%s, email:%s, code:%s]\n", template, email, smsCode.Code)

	var config messageModel.MailConfig
	db.DB().First(&config)
	if config.ID <= 0 {
		return errors.New("未配置邮件信息")
	}

	var sender Sender

	sender, err := NewMailSender(
		[]string{email},
		"code",
		WithConfig(
			config.Host,
			config.Port,
			config.Transport,
			config.Username,
			config.Password,
			config.User,
			config.FromName,
		),
		WithSubject("验证码"),
		WithParams(templateParam),
	)
	if err != nil {
		return errors.Wrap(err, "验证码发送失败")
	}

	err = sender.SendMail()
	if err != nil {
		return errors.Wrap(err, "验证码发送失败")
	}

	return nil
}

// ValidMailCode 验证邮件验证码
// @param email 接收者邮箱地址
// @param emailCode 验证码
func ValidMailCode(email, template, emailCode string) error {
	// 验证码缓存
	key := redis.Key("email_code_" + template + ":" + email)

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
		if code != emailCode {
			return api.EMValidCodeError
		}

		redisClient.Del(context.Background(), key)
	} else {
		return api.EMValidCodeExpireError
	}

	return nil
}
