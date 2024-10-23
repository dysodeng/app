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
