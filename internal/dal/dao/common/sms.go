package common

import (
	"context"

	"github.com/dysodeng/app/internal/infrastructure/persistence/model/common"

	"github.com/dysodeng/app/internal/pkg/db"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// SmsDao 短信配置数据访问层
type SmsDao interface {
	Config(ctx context.Context) (*common.SmsConfig, error)
	SaveConfig(ctx context.Context, config common.SmsConfig) error
	Template(ctx context.Context, template string) (*common.SmsTemplate, error)
}

type smsDao struct {
	baseTraceSpanName string
}

func NewSmsDao() SmsDao {
	return &smsDao{
		baseTraceSpanName: "dal.dao.common.SmsDao",
	}
}

func (dao *smsDao) Config(ctx context.Context) (*common.SmsConfig, error) {
	var config common.SmsConfig
	err := db.DB().WithContext(ctx).First(&config).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &config, nil
}

func (dao *smsDao) SaveConfig(ctx context.Context, config common.SmsConfig) error {
	var conf common.SmsConfig
	db.DB().WithContext(ctx).First(&conf)
	var err error
	if conf.ID > 0 {
		err = db.DB().WithContext(ctx).Model(&common.SmsConfig{}).Where("id=?", conf.ID).
			Updates(map[string]interface{}{
				"sms_type":          config.SmsType,
				"app_key":           config.AppKey,
				"secret_key":        config.SecretKey,
				"free_sign_name":    config.FreeSignName,
				"valid_code_expire": config.ValidCodeExpire,
			}).Error
	} else {
		err = db.DB().WithContext(ctx).Create(&config).Error
	}
	return err
}

func (dao *smsDao) Template(ctx context.Context, template string) (*common.SmsTemplate, error) {
	var templateConfig common.SmsTemplate
	err := db.DB().WithContext(ctx).Where("template=?", template).First(&templateConfig).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &templateConfig, nil
}
