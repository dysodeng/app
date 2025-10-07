package websocket

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dysodeng/app/internal/infrastructure/shared/logger"
	"github.com/dysodeng/app/internal/infrastructure/shared/token"
)

type Authenticator interface {
	// Authenticate 验证请求是否合法，第一个返回值为用户 id，第二个返回值为错误
	Authenticate(r *http.Request) (map[string]interface{}, error)
}

var _ Authenticator = &JWTAuthenticator{}

type JWTAuthenticator struct{}

func (J *JWTAuthenticator) Authenticate(r *http.Request) (map[string]interface{}, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		authHeader = r.URL.Query().Get("token")
	} else {
		if !strings.HasPrefix(authHeader, "Bearer ") || len(authHeader) < 8 {
			return nil, fmt.Errorf("invalid authorization format")
		}
	}

	if strings.Contains(authHeader, "DebugToken") { // DebugToken
		tokens := strings.Split(authHeader, "_") // DebugToken-user-1-1 // 第2项为用户类型，第3项为用户id，第4项为平台类型
		if len(tokens) != 4 {
			return nil, fmt.Errorf("token无效")
		}

		return map[string]interface{}{
			"user_type":     tokens[1],
			"user_id":       tokens[2],
			"platform_type": tokens[3],
		}, nil
	}

	claims, err := token.VerifyToken(authHeader)
	if err != nil {
		logger.Error(context.Background(), err.Error(), logger.ErrorField(err))
		return nil, fmt.Errorf("token is expired")
	}

	return claims, nil
}
