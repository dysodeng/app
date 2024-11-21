package user

import (
	"context"

	"github.com/dysodeng/app/internal/pkg/telemetry/trace"

	userDao "github.com/dysodeng/app/internal/dal/dao/user"
	userModel "github.com/dysodeng/app/internal/dal/model/user"
	"github.com/dysodeng/app/internal/pkg/filesystem"
	"github.com/dysodeng/app/internal/pkg/helper"
	"github.com/dysodeng/app/internal/pkg/model"
	userDo "github.com/dysodeng/app/internal/service/do/user"
)

type DomainService struct {
	ctx               context.Context
	userDao           *userDao.Dao
	baseTraceSpanName string
}

func NewUserDomainService(ctx context.Context) *DomainService {
	return &DomainService{
		ctx:               ctx,
		userDao:           userDao.NewUserDao(ctx),
		baseTraceSpanName: "service.domain.user.UserDomainService",
	}
}

func (ds *DomainService) CreateUser(userInfo userDo.User) (*userDo.User, error) {
	spanCtx, span := trace.Tracer().Start(ds.ctx, ds.baseTraceSpanName+".CreateUser")
	defer span.End()

	uDao := userDao.NewUserDao(spanCtx)
	user, err := uDao.CreateUser(userModel.User{
		Telephone: userInfo.Telephone,
		Password:  helper.GeneratePassword([]byte(userInfo.Password)),
		RealName:  userInfo.RealName,
		Nickname:  userInfo.Nickname,
		Avatar:    filesystem.Instance().OriginalPath(userInfo.Avatar),
		Gender:    userInfo.Gender,
		Birthday:  userInfo.Birthday,
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
		Birthday:  user.Birthday,
		Status:    user.Status.Uint(),
	}, nil
}

func (ds *DomainService) Info(userId uint64) (*userDo.User, error) {
	uDao := userDao.NewUserDao(ds.ctx)
	user, err := uDao.Info(userId)
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
		Birthday:  user.Birthday,
		Status:    user.Status.Uint(),
	}, nil
}

func (ds *DomainService) ListUser(page, pageSize int, condition map[string]interface{}) ([]userDo.User, int64, error) {
	uDao := userDao.NewUserDao(ds.ctx)
	userList, count, err := uDao.ListUser(page, pageSize, condition)
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
			Birthday:  userList[i].Birthday,
			Status:    userList[i].Status.Uint(),
		})
	}

	return userDoList, count, nil
}

func (ds *DomainService) UpdateUser(userInfo userDo.User) (*userDo.User, error) {
	uDao := userDao.NewUserDao(ds.ctx)
	u := userModel.User{
		PrimaryKeyID: model.PrimaryKeyID{ID: userInfo.Id},
		Telephone:    userInfo.Telephone,
		RealName:     userInfo.RealName,
		Nickname:     userInfo.Nickname,
		Avatar:       filesystem.Instance().OriginalPath(userInfo.Avatar),
		Gender:       userInfo.Gender,
		Birthday:     userInfo.Birthday,
		Status:       model.BinaryStatusByUint(userInfo.Status),
	}
	if userInfo.Password != "" {
		u.Password = helper.GeneratePassword([]byte(userInfo.Password))
	}
	user, err := uDao.UpdateUser(u)
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
		Birthday:  user.Birthday,
		Status:    user.Status.Uint(),
	}, nil
}

func (ds *DomainService) DeleteUser(userId uint64) error {
	return userDao.NewUserDao(ds.ctx).DeleteUser(userId)
}
