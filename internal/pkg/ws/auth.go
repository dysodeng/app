package ws

import (
	"fmt"
	"net/http"

	"github.com/dysodeng/app/internal/pkg/api/token"
	"github.com/dysodeng/app/internal/pkg/helper"
)

type Authenticator interface {
	// Authenticate 验证请求是否合法，第一个返回值为用户 id，第二个返回值为错误
	Authenticate(r *http.Request) (string, error)
}

var _ Authenticator = &JWTAuthenticator{}

type JWTAuthenticator struct{}

func (J *JWTAuthenticator) Authenticate(r *http.Request) (string, error) {
	tokenString := r.Header.Get("Authorization")
	claims, err := token.VerifyToken(tokenString)
	if err != nil {
		return "", err
	}

	userID := helper.IfaceConvertInt64(claims["user_id"])
	if userID > 0 {
		return fmt.Sprintf("%d", userID), nil
	}
	return "", fmt.Errorf("userId should be string")
}
