package repository

import (
	"context"

	"github.com/dysodeng/app/internal/domain/common/model"
)

// MailRepository 邮件仓储接口
type MailRepository interface {
	Config(ctx context.Context) (*model.MailConfig, error)
	SaveConfig(ctx context.Context, config *model.MailConfig) error
}
