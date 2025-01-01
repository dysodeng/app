package user

import (
	"context"

	userDao "github.com/dysodeng/app/internal/dal/dao/user"
	userModel "github.com/dysodeng/app/internal/dal/model/user"
	"github.com/dysodeng/app/internal/pkg/filesystem"
	"github.com/dysodeng/app/internal/pkg/helper"
	"github.com/dysodeng/app/internal/pkg/model"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	userDo "github.com/dysodeng/app/internal/service/do/user"
)

// DomainService 用户领域服务
type DomainService interface {
	Info(ctx context.Context, userId uint64) (*userDo.User, error)
	ListUser(ctx context.Context, page, pageSize int, condition map[string]interface{}) ([]userDo.User, int64, error)
	CreateUser(ctx context.Context, userInfo userDo.User) (*userDo.User, error)
	UpdateUser(ctx context.Context, userInfo userDo.User) (*userDo.User, error)
	DeleteUser(ctx context.Context, userId uint64) error
}

// domainService 用户领域服务
type domainService struct {
	userDao           userDao.Dao
	baseTraceSpanName string
}

var userDomainServiceInstance DomainService

func NewUserDomainService(userDao userDao.Dao) DomainService {
	if userDomainServiceInstance == nil {
		userDomainServiceInstance = &domainService{
			baseTraceSpanName: "service.domain.user.UserDomainService",
			userDao:           userDao,
		}
	}
	return userDomainServiceInstance
}

func (ds *domainService) CreateUser(ctx context.Context, userInfo userDo.User) (*userDo.User, error) {
	spanCtx, span := trace.Tracer().Start(ctx, ds.baseTraceSpanName+".CreateUser")
	defer span.End()

	password, err := helper.GeneratePassword(userInfo.Password)
	if err != nil {
		trace.Error(err, span)
		return nil, err
	}

	user, err := ds.userDao.CreateUser(spanCtx, userModel.User{
		Telephone: userInfo.Telephone,
		Password:  password,
		RealName:  userInfo.RealName,
		Nickname:  userInfo.Nickname,
		Avatar:    filesystem.Instance().OriginalPath(userInfo.Avatar),
		Gender:    userInfo.Gender,
		Birthday:  model.JSONDate{Time: userInfo.Birthday},
		Status:    model.BinaryStatusYes,
	})
	if err != nil {
		trace.Error(err, span)
		return nil, err
	}
	return &userDo.User{
		Id:        user.ID,
		Telephone: user.Telephone,
		RealName:  user.RealName,
		Nickname:  user.Nickname,
		Avatar:    filesystem.Instance().FullPath(user.Avatar),
		Gender:    user.Gender,
		Birthday:  user.Birthday.Time,
		Status:    user.Status.Uint(),
	}, nil
}

func (ds *domainService) Info(ctx context.Context, userId uint64) (*userDo.User, error) {
	spanCtx, span := trace.Tracer().Start(ctx, ds.baseTraceSpanName+".Info")
	defer span.End()
	user, err := ds.userDao.Info(spanCtx, userId)
	if err != nil {
		return nil, err
	}
	return &userDo.User{
		Id:        user.ID,
		Telephone: user.Telephone,
		RealName:  user.RealName,
		Nickname:  user.Nickname,
		Avatar:    filesystem.Instance().FullPath(user.Avatar),
		Gender:    user.Gender,
		Birthday:  user.Birthday.Time,
		Status:    user.Status.Uint(),
	}, nil
}

func (ds *domainService) ListUser(ctx context.Context, page, pageSize int, condition map[string]interface{}) ([]userDo.User, int64, error) {
	userList, count, err := ds.userDao.ListUser(ctx, page, pageSize, condition)
	if err != nil {
		return nil, 0, err
	}

	var userDoList []userDo.User
	for i := range userList {
		userDoList = append(userDoList, userDo.User{
			Id:        userList[i].ID,
			Telephone: userList[i].Telephone,
			RealName:  userList[i].RealName,
			Nickname:  userList[i].Nickname,
			Avatar:    filesystem.Instance().FullPath(userList[i].Avatar),
			Gender:    userList[i].Gender,
			Birthday:  userList[i].Birthday.Time,
			Status:    userList[i].Status.Uint(),
		})
	}

	return userDoList, count, nil
}

func (ds *domainService) UpdateUser(ctx context.Context, userInfo userDo.User) (*userDo.User, error) {
	u := userModel.User{
		PrimaryKeyID: model.PrimaryKeyID{ID: userInfo.Id},
		Telephone:    userInfo.Telephone,
		RealName:     userInfo.RealName,
		Nickname:     userInfo.Nickname,
		Avatar:       filesystem.Instance().OriginalPath(userInfo.Avatar),
		Gender:       userInfo.Gender,
		Birthday:     model.JSONDate{Time: userInfo.Birthday},
		Status:       model.BinaryStatusByUint(userInfo.Status),
	}
	if userInfo.Password != "" {
		password, err := helper.GeneratePassword(userInfo.Password)
		if err != nil {
			return nil, err
		}
		u.Password = password
	}
	user, err := ds.userDao.UpdateUser(ctx, u)
	if err != nil {
		return nil, err
	}
	return &userDo.User{
		Id:        user.ID,
		Telephone: user.Telephone,
		RealName:  user.RealName,
		Nickname:  user.Nickname,
		Avatar:    filesystem.Instance().FullPath(user.Avatar),
		Gender:    user.Gender,
		Birthday:  user.Birthday.Time,
		Status:    user.Status.Uint(),
	}, nil
}

func (ds *domainService) DeleteUser(ctx context.Context, userId uint64) error {
	return ds.userDao.DeleteUser(ctx, userId)
}
