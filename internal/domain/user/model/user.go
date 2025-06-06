package model

import (
	"context"
	"time"

	"github.com/dysodeng/app/internal/infrastructure/persistence/model/user"
	"github.com/dysodeng/app/internal/pkg/model"
	"github.com/dysodeng/app/internal/pkg/storage"
)

type User struct {
	ID        uint64             `json:"id"`
	Telephone string             `json:"telephone"`
	Password  string             `json:"password"`
	RealName  string             `json:"real_name"`
	Avatar    string             `json:"avatar"`
	Nickname  string             `json:"nickname"`
	Gender    uint8              `json:"gender"`
	Birthday  time.Time          `json:"birthday"`
	Status    model.BinaryStatus `json:"status"`
}

func (u *User) Validate() error {
	return nil
}

func UserFromModel(user *user.User) *User {
	avatar := user.Avatar
	if avatar != "" {
		avatar = storage.Instance().FullUrl(context.Background(), avatar)
	}
	return &User{
		ID:        user.ID,
		Telephone: user.Telephone,
		Password:  user.Password,
		RealName:  user.RealName,
		Avatar:    avatar,
		Nickname:  user.Nickname,
		Gender:    user.Gender,
		Birthday:  user.Birthday.Time,
		Status:    model.BinaryStatusByUint(user.Status),
	}
}

func UserListFromModel(users []user.User) []User {
	result := make([]User, len(users))
	for i, u := range users {
		avatar := u.Avatar
		if avatar != "" {
			avatar = storage.Instance().FullUrl(context.Background(), avatar)
		}
		result[i] = User{
			ID:        u.ID,
			Telephone: u.Telephone,
			Password:  u.Password,
			RealName:  u.RealName,
			Avatar:    avatar,
			Nickname:  u.Nickname,
			Gender:    u.Gender,
			Birthday:  u.Birthday.Time,
			Status:    model.BinaryStatusByUint(u.Status),
		}
	}
	return result
}

func (u *User) ToModel() *user.User {
	avatar := u.Avatar
	if avatar != "" {
		avatar = storage.Instance().FullUrl(context.Background(), avatar)
	}
	dataModel := &user.User{
		Telephone: u.Telephone,
		Password:  u.Password,
		RealName:  u.RealName,
		Avatar:    avatar,
		Nickname:  u.Nickname,
		Gender:    u.Gender,
		Birthday:  model.JSONDate{Time: u.Birthday},
	}
	if u.ID > 0 {
		dataModel.ID = u.ID
	}
	return dataModel
}
