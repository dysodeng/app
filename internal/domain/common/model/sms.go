package model

import "github.com/dysodeng/app/internal/infrastructure/persistence/model/common"

type SmsConfig struct {
	ID              uint64 `json:"id"`
	SmsType         string `json:"sms_type"`
	AppKey          string `json:"app_key"`
	SecretKey       string `json:"secret_key"`
	FreeSignName    string `json:"free_sign_name"`
	ValidCodeExpire uint   `json:"valid_code_expire"`
}

func (conf *SmsConfig) ToModel() *common.SmsConfig {
	dataModel := &common.SmsConfig{
		SmsType:         conf.SmsType,
		AppKey:          conf.AppKey,
		SecretKey:       conf.SecretKey,
		FreeSignName:    conf.FreeSignName,
		ValidCodeExpire: conf.ValidCodeExpire,
	}
	if conf.ID > 0 {
		dataModel.ID = conf.ID
	}
	return dataModel
}

func SmsConfigFromModel(dataModel *common.SmsConfig) *SmsConfig {
	return &SmsConfig{
		ID:              dataModel.ID,
		SmsType:         dataModel.SmsType,
		AppKey:          dataModel.AppKey,
		SecretKey:       dataModel.SecretKey,
		FreeSignName:    dataModel.FreeSignName,
		ValidCodeExpire: dataModel.ValidCodeExpire,
	}
}

type SmsTemplate struct {
	ID           uint64
	TemplateName string
	Template     string
	TemplateId   string
}

func (template *SmsTemplate) ToModel() *common.SmsTemplate {
	dataModel := &common.SmsTemplate{
		TemplateName: template.TemplateName,
		Template:     template.Template,
		TemplateId:   template.TemplateId,
	}
	if template.ID > 0 {
		dataModel.ID = template.ID
	}
	return dataModel
}

func SmsTemplateFromModel(dataModel *common.SmsTemplate) *SmsTemplate {
	return &SmsTemplate{
		ID:           dataModel.ID,
		TemplateName: dataModel.TemplateName,
		Template:     dataModel.Template,
		TemplateId:   dataModel.TemplateId,
	}
}
