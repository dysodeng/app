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

type mailRepository struct {
	baseTraceSpanName string
	txManager         transactions.TransactionManager
}

func NewMailRepository(txManager transactions.TransactionManager) repository.MailRepository {
	return &mailRepository{
		baseTraceSpanName: "infrastructure.persistence.repository.common.MailRepository",
		txManager:         txManager,
	}
}

func (repo *mailRepository) Config(ctx context.Context) (*model.MailConfig, error) {
	var config common.MailConfig
	err := db.DB().WithContext(ctx).First(&config).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return model.MailConfigFormModel(&config), nil
}

func (repo *mailRepository) SaveConfig(ctx context.Context, config *model.MailConfig) error {
	var conf common.MailConfig
	db.DB().WithContext(ctx).First(&conf)
	var err error
	if conf.ID > 0 {
		err = db.DB().WithContext(ctx).Model(&common.MailConfig{}).Where("id=?", conf.ID).
			Updates(map[string]interface{}{
				"from_name": config.FromName,
				"host":      config.Host,
				"password":  config.Password,
				"port":      config.Port,
				"transport": config.Transport,
				"user":      config.User,
				"username":  config.Username,
			}).Error
	} else {
		dataModel := config.ToModel()
		err = db.DB().WithContext(ctx).Create(dataModel).Error
	}
	return err
}
