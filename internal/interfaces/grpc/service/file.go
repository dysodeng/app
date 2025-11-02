package service

import (
	"context"
	"time"

	"github.com/dysodeng/rpc/metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonV1 "github.com/dysodeng/app/api/generated/go/proto/common/v1"
	v1 "github.com/dysodeng/app/api/generated/go/proto/file/v1"
	fileApplicationService "github.com/dysodeng/app/internal/application/file/service"
	"github.com/dysodeng/app/internal/infrastructure/config"
)

type FileService struct {
	metadata.UnimplementedServiceRegister
	v1.UnimplementedFileServiceServer
	fileApplicationService fileApplicationService.FileApplicationService
}

func NewFileService(
	fileApplicationService fileApplicationService.FileApplicationService,
) *FileService {
	return &FileService{
		fileApplicationService: fileApplicationService,
	}
}

func (svc *FileService) RegisterMetadata() metadata.ServiceRegisterMetadata {
	return metadata.ServiceRegisterMetadata{
		AppName:     config.GlobalConfig.App.Name,
		ServiceName: "file.FileService",
		Version:     metadata.DefaultVersion,
		Env:         config.GlobalConfig.App.Environment,
		Tags:        []string{"file", "uploader", "file.storage"},
	}
}

func (svc *FileService) Metadata(_ context.Context, _ *commonV1.MetadataRequest) (*commonV1.MetadataResponse, error) {
	meta := svc.RegisterMetadata()
	return &commonV1.MetadataResponse{
		Code:    commonV1.Code_SUCCESS,
		Message: "success",
		Metadata: &commonV1.Metadata{
			AppName:     meta.AppName,
			ServiceName: meta.ServiceName,
			Version:     meta.Version,
			Env:         meta.Env,
			Tags:        meta.Tags,
		},
	}, nil
}

func (svc *FileService) FileInfo(ctx context.Context, req *v1.FileInfoRequest) (*v1.FileInfoResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "file id is empty")
	}

	res, err := svc.fileApplicationService.FileInfo(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &v1.FileInfoResponse{
		Code:    commonV1.Code_SUCCESS,
		Message: "success",
		File: &v1.File{
			Id:        res.ID.String(),
			MediaType: svc.mediaType(res.MediaType),
			Name:      res.Name,
			NameIndex: res.NameIndex,
			Path:      res.Path,
			Ext:       res.Ext,
			MimeType:  res.MimeType,
			Status:    uint32(res.Status),
			CreatedAt: res.CreatedAt.Format(time.DateTime),
		},
	}, nil
}

func (svc *FileService) FileReference(ctx context.Context, req *v1.FileReferenceRequest) (*v1.FileReferenceResponse, error) {
	return nil, nil
}

func (svc *FileService) RevokeFileReference(ctx context.Context, req *v1.RevokeFileReferenceRequest) (*v1.RevokeFileReferenceResponse, error) {
	return nil, nil
}

func (svc *FileService) mediaType(mediaType uint8) v1.MediaType {
	var t v1.MediaType
	switch mediaType {
	case 1:
		t = v1.MediaType_Image
	case 2:
		t = v1.MediaType_Video
	case 3:
		t = v1.MediaType_Audio
	case 4:
		t = v1.MediaType_Document
	}
	return t
}
