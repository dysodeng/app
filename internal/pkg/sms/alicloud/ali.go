package alicloud

import (
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/pkg/errors"
)

// AliCloud 阿里云短信
type AliCloud struct {
	option *option
}

// NewAliCloudSmsSender 创建阿里去短信发送
// @param phoneNumber string 收信人手机号
// @param templateCode string 短信模板ID
// config Config 短信配置
// opts []Option 短信额外选项
func NewAliCloudSmsSender(phoneNumber string, templateId string, config Config, opts ...Option) *AliCloud {
	o := &option{
		phoneNumber: phoneNumber,
		templateId:  templateId,
	}
	config(o)
	for _, opt := range opts {
		opt(o)
	}

	sender := &AliCloud{
		option: o,
	}

	return sender
}

// SendSms 发送短信
func (aliCloud *AliCloud) SendSms() (bool, error) {
	client, err := sdk.NewClientWithOptions(
		"default",
		sdk.NewConfig(),
		&credentials.AccessKeyCredential{
			AccessKeyId:     aliCloud.option.accessKey,
			AccessKeySecret: aliCloud.option.accessSecret,
		},
	)
	if err != nil {
		return false, err
	}

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https"
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"

	params := make(map[string]string)
	params["PhoneNumbers"] = aliCloud.option.phoneNumber
	params["SignName"] = aliCloud.option.signName
	params["TemplateCode"] = aliCloud.option.templateId
	if aliCloud.option.templateParam != "" {
		params["TemplateParam"] = aliCloud.option.templateParam
	}

	request.QueryParams = params

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		log.Println("AliCloudSender->Send: " + err.Error())
		return false, errors.New("短信发送失败")
	}

	r := response.GetHttpContentString()
	log.Println("send success")
	log.Println(r)

	return true, nil
}
