package repository

import (
	"context"

	"github.com/dysodeng/app/internal/domain/common/model"
)

// SmsRepository 短信配置仓储接口
type SmsRepository interface {
	Config(ctx context.Context) (*model.SmsConfig, error)
	SaveConfig(ctx context.Context, config *model.SmsConfig) error
	Template(ctx context.Context, template string) (*model.SmsTemplate, error)
}
