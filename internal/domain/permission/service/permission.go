package service

import (
	"context"

	"github.com/dysodeng/app/internal/domain/permission/model"
	"github.com/dysodeng/app/internal/domain/permission/repository"
	sharedErrors "github.com/dysodeng/app/internal/domain/shared/errors"
	sharedVO "github.com/dysodeng/app/internal/domain/shared/valueobject"
	sharedModel "github.com/dysodeng/app/internal/infrastructure/shared/model"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/trace"
)

// PermissionDomainService 管理权限领域服务
type PermissionDomainService interface{}

type permissionDomainService struct {
	baseTraceSpanName string
	adminRepository   repository.AdminRepository
}

func NewPermissionDomainService(adminRepository repository.AdminRepository) PermissionDomainService {
	return &permissionDomainService{
		baseTraceSpanName: "domain.permission.service.PermissionDomainService",
		adminRepository:   adminRepository,
	}
}

func (svc *permissionDomainService) CreateAdmin(
	ctx context.Context,
	username,
	password,
	realName,
	telephone,
	remark string,
	status uint8,
) (*model.Admin, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".CreateAdmin")
	defer span.End()

	usernameVO, err := sharedVO.NewUsername(username)
	if err != nil {
		return nil, err
	}
	passwordVO, err := sharedVO.NewPassword(password)
	if err != nil {
		return nil, err
	}

	if ok, err := svc.adminRepository.ExistsByUsername(spanCtx, usernameVO); err != nil {
		return nil, err
	} else if ok {
		return nil, sharedErrors.ErrSharedUsernameAlreadyExists
	}

	telephoneVO, err := sharedVO.NewTelephone(telephone)
	if err != nil {
		return nil, err
	}

	admin, err := model.NewAdmin(usernameVO, passwordVO, realName, telephoneVO, remark, sharedModel.BinaryStatusByUint(status))
	if err != nil {
		return nil, err
	}

	return admin, nil
}
