package common

import (
	"context"

	commonDao "github.com/dysodeng/app/internal/dal/dao/common"
	"github.com/dysodeng/app/internal/pkg/trace"
)

// MailDomainService 邮件领域服务
type MailDomainService struct {
	ctx               context.Context
	mailDao           *commonDao.MailDao
	baseTraceSpanName string
}

func NewMailDomainService(ctx context.Context) *MailDomainService {
	baseTraceSpanName := "domain.common.MailDomainService"
	traceCtx := trace.New().NewSpan(ctx, baseTraceSpanName)
	return &MailDomainService{
		ctx:               traceCtx,
		mailDao:           commonDao.NewMailDao(traceCtx),
		baseTraceSpanName: baseTraceSpanName,
	}
}
