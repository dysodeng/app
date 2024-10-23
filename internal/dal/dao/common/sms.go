package common

import (
	"context"

	"github.com/dysodeng/app/internal/dal/model/common"
	"github.com/dysodeng/app/internal/pkg/db"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type SmsDao struct {
	ctx context.Context
}

func NewSmsDao(ctx context.Context) *SmsDao {
	return &SmsDao{ctx: ctx}
}

func (dao *SmsDao) Config() (*common.SmsConfig, error) {
	var config common.SmsConfig
	err := db.DB().WithContext(dao.ctx).First(&config).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &config, nil
}

func (dao *SmsDao) SaveConfig(config common.SmsConfig) error {
	var conf common.SmsConfig
	db.DB().WithContext(dao.ctx).First(&conf)
	var err error
	if conf.ID > 0 {
		err = db.DB().WithContext(dao.ctx).Model(&common.SmsConfig{}).Where("id=?", conf.ID).
			Updates(map[string]interface{}{
				"sms_type":          config.SmsType,
				"app_key":           config.AppKey,
				"secret_key":        config.SecretKey,
				"free_sign_name":    config.FreeSignName,
				"valid_code_expire": config.ValidCodeExpire,
			}).Error
	} else {
		err = db.DB().WithContext(dao.ctx).Create(&config).Error
	}
	return err
}

func (dao *SmsDao) Template(template string) (*common.SmsTemplate, error) {
	var templateConfig common.SmsTemplate
	err := db.DB().WithContext(dao.ctx).Where("template=?", template).First(&templateConfig).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &templateConfig, nil
}
