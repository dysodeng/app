package response

import (
	"time"

	"github.com/dysodeng/app/internal/domain/user/model"
)

type UserResponse struct {
	ID        uint64 `json:"id"`
	Telephone string `json:"telephone"`
	Avatar    string `json:"avatar"`
	Nickname  string `json:"nickname"`
	RealName  string `json:"real_name"`
	Gender    uint8  `json:"gender"`
	Birthday  string `json:"birthday"`
	Status    uint8  `json:"status"`
}

func FromDomainUser(user *model.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Telephone: user.Telephone,
		Avatar:    user.Avatar,
		Nickname:  user.Nickname,
		RealName:  user.RealName,
		Gender:    user.Gender,
		Birthday:  user.Birthday.Format(time.DateOnly),
		Status:    user.Status.Uint(),
	}
}

func ListFromDomainUser(users []model.User) []UserResponse {
	result := make([]UserResponse, len(users))
	for i, user := range users {
		result[i] = UserResponse{
			ID:        user.ID,
			Telephone: user.Telephone,
			Avatar:    user.Avatar,
			Nickname:  user.Nickname,
			RealName:  user.RealName,
			Gender:    user.Gender,
			Birthday:  user.Birthday.Format(time.DateOnly),
			Status:    user.Status.Uint(),
		}
	}
	return result
}

type UserListResponse struct {
	List  []UserResponse `json:"list"`
	Total int64          `json:"total"`
}
