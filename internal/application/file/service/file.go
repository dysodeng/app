package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/dysodeng/app/internal/application/file/dto/response"
	"github.com/dysodeng/app/internal/domain/file/service"
	"github.com/dysodeng/app/internal/domain/shared/errors"
	"github.com/dysodeng/app/internal/infrastructure/shared/logger"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/trace"
)

// FileApplicationService 文件应用服务
type FileApplicationService interface {
	// FileInfo 获取文件信息
	FileInfo(ctx context.Context, id string) (*response.FileResponse, error)
}

type fileApplicationService struct {
	baseTraceSpanName string
	fileDomainService service.FileDomainService
}

func NewFileApplicationService(fileDomainService service.FileDomainService) FileApplicationService {
	return &fileApplicationService{
		baseTraceSpanName: "application.file.FileApplicationService",
		fileDomainService: fileDomainService,
	}
}

func (svc *fileApplicationService) FileInfo(ctx context.Context, id string) (*response.FileResponse, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".FileInfo")
	defer span.End()

	fileId, err := uuid.Parse(id)
	if err != nil {
		logger.Warn(spanCtx, "文件ID格式错误", logger.ErrorField(err))
		return nil, errors.NewFileError("FILE_ID_INVALID", "文件ID格式错误", nil).Wrap(err)
	}

	info, err := svc.fileDomainService.Info(spanCtx, fileId)
	if err != nil {
		return nil, err
	}

	return &response.FileResponse{
		ID:        info.ID,
		Name:      info.Name.String(),
		NameIndex: info.NameIndex,
		Path:      info.Path,
		Size:      info.Size,
		Ext:       info.Ext,
		MediaType: info.MediaType.ToInt(),
		MimeType:  info.MimeType,
		Status:    info.Status,
		CreatedAt: info.CreatedAt,
	}, nil
}
