package common

import (
	"context"
	"github.com/dysodeng/app/internal/domain/common/model"
	"github.com/dysodeng/app/internal/domain/common/repository"
	"github.com/dysodeng/app/internal/infrastructure/persistence/model/common"
	"github.com/dysodeng/app/internal/infrastructure/transactions"
	"github.com/dysodeng/app/internal/pkg/db"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type smsRepository struct {
	baseTraceSpanName string
	txManager         transactions.TransactionManager
}

func NewSmsRepository(txManager transactions.TransactionManager) repository.SmsRepository {
	return &smsRepository{
		baseTraceSpanName: "infrastructure.persistence.repository.common.SmsRepository",
		txManager:         txManager,
	}
}

func (repo *smsRepository) Config(ctx context.Context) (*model.SmsConfig, error) {
	var config common.SmsConfig
	err := db.DB().WithContext(ctx).First(&config).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return model.SmsConfigFromModel(&config), nil
}

func (repo *smsRepository) SaveConfig(ctx context.Context, config *model.SmsConfig) error {
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
		dataModel := config.ToModel()
		err = db.DB().WithContext(ctx).Create(dataModel).Error
	}
	return err
}

func (repo *smsRepository) Template(ctx context.Context, template string) (*model.SmsTemplate, error) {
	var templateConfig common.SmsTemplate
	err := db.DB().WithContext(ctx).Where("template=?", template).First(&templateConfig).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return model.SmsTemplateFromModel(&templateConfig), nil
}
