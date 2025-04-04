package ws

import (
	"fmt"
	"net/http"

	"github.com/dysodeng/app/internal/pkg/token"
)

type Authenticator interface {
	// Authenticate 验证请求是否合法，第一个返回值为用户 id，第二个返回值为错误
	Authenticate(r *http.Request) (map[string]interface{}, error)
}

var _ Authenticator = &JWTAuthenticator{}

type JWTAuthenticator struct{}

func (J *JWTAuthenticator) Authenticate(r *http.Request) (map[string]interface{}, error) {
	tokenString := r.Header.Get("Authorization")
	claims, err := token.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	var userId string
	switch claims["user_type"] {
	case "user":
		userId = fmt.Sprintf("%d", claims["user_id"])
	case "ams":
		userId = fmt.Sprintf("%d", claims["admin_id"])
	}

	return map[string]interface{}{
		"user_id":   userId,
		"user_type": claims["user_type"],
	}, nil
}
