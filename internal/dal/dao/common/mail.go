package common

import (
	"context"

	"github.com/dysodeng/app/internal/infrastructure/persistence/model/common"

	"github.com/dysodeng/app/internal/pkg/db"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// MailDao 邮件配置数据访问层
type MailDao interface {
	Config(ctx context.Context) (*common.MailConfig, error)
	SaveConfig(ctx context.Context, config common.MailConfig) error
}

// mailDao 邮件配置数据访问对象
type mailDao struct {
	baseTraceSpanName string
}

func NewMailDao() MailDao {
	return &mailDao{
		baseTraceSpanName: "dal.dao.common.MailDao",
	}
}

func (dao *mailDao) Config(ctx context.Context) (*common.MailConfig, error) {
	var config common.MailConfig
	err := db.DB().WithContext(ctx).First(&config).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &config, nil
}

func (dao *mailDao) SaveConfig(ctx context.Context, config common.MailConfig) error {
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
		err = db.DB().WithContext(ctx).Create(&config).Error
	}
	return err
}
