package common

import (
	"context"

	"github.com/dysodeng/app/internal/dal/model/common"
	"github.com/dysodeng/app/internal/pkg/db"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// MailDao 邮件配置数据访问对象
type MailDao struct {
	ctx context.Context
}

func NewMailDao(ctx context.Context) *MailDao {
	return &MailDao{ctx: ctx}
}

func (dao *MailDao) Config() (*common.MailConfig, error) {
	var config common.MailConfig
	err := db.DB().WithContext(dao.ctx).First(&config).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &config, nil
}

func (dao *MailDao) SaveConfig(config common.MailConfig) error {
	var conf common.MailConfig
	db.DB().WithContext(dao.ctx).First(&conf)
	var err error
	if conf.ID > 0 {
		err = db.DB().WithContext(dao.ctx).Model(&common.MailConfig{}).Where("id=?", conf.ID).
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
		err = db.DB().WithContext(dao.ctx).Create(&config).Error
	}
	return err
}
