package decorator

import (
	filePort "github.com/dysodeng/app/internal/domain/file/port"
	fileRepo "github.com/dysodeng/app/internal/domain/file/repository"
	fileDomainSvc "github.com/dysodeng/app/internal/domain/file/service"
)

// NewFileDomainServiceWithTracing 文件领域服务链路追踪装饰器
func NewFileDomainServiceWithTracing(
	fileRepository fileRepo.FileRepository,
) fileDomainSvc.FileDomainService {
	base := fileDomainSvc.NewFileDomainService(fileRepository)
	return NewTracedFileDomainService(base)
}

// NewUploaderDomainServiceWithTracing 文件上传领域服务链路追踪装饰器
func NewUploaderDomainServiceWithTracing(
	fileRepository fileRepo.FileRepository,
	uploaderRepository fileRepo.UploaderRepository,
	storage filePort.FileStorage,
	policy filePort.FilePolicy,
) fileDomainSvc.UploaderDomainService {
	base := fileDomainSvc.NewUploaderDomainService(
		fileRepository,
		uploaderRepository,
		storage,
		policy,
	)
	return NewTracedUploaderDomainService(base)
}
