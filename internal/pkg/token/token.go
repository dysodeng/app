package token

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dysodeng/app/internal/config"

	"github.com/golang-jwt/jwt/v4"

	"github.com/pkg/errors"
)

const JwtAuthIdentifier = "github.com/dysodeng/app/auth"

// Token token 数据结构
type Token struct {
	Exists             uint8       `json:"exists"`
	Token              json.Token  `json:"token"`
	Expire             int64       `json:"expire"`
	RefreshToken       json.Token  `json:"refresh_token"`
	RefreshTokenExpire int64       `json:"refresh_token_expire"`
	Attach             interface{} `json:"attach,omitempty"`
}

// AuthCodeToken 核验码token数据结构
type AuthCodeToken struct {
	Token  json.Token
	Expire int64 `json:"expire"`
}

// GenerateToken 构建用户token
func GenerateToken(userType string, data map[string]interface{}, attach map[string]interface{}) (Token, error) {
	currentTime := time.Now().Unix()
	var tokenMethod *jwt.Token
	var refreshTokenMethod *jwt.Token
	var expire int64
	var refreshTokenExpire int64

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	switch userType {
	case "user": // 终端用户
		expire = 30 * 24 * 3600
		refreshTokenExpire = 2 * 30 * 24 * 3600
		// Token
		tokenMethod = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":          data["user_id"],
			"platform_type":    data["platform_type"],
			"user_type":        "user",
			"is_refresh_token": false,
			"token_type":       "biz", // 业务token
			"iss":              JwtAuthIdentifier,
			"aud":              JwtAuthIdentifier,
			"iat":              currentTime,
			"exp":              currentTime + int64(expire),
		})

		// RefreshToken
		refreshTokenMethod = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":          data["user_id"],
			"platform_type":    data["platform_type"],
			"user_type":        "user",
			"is_refresh_token": true,
			"token_type":       "biz", // 业务token
			"iss":              JwtAuthIdentifier,
			"aud":              JwtAuthIdentifier,
			"iat":              currentTime,
			"exp":              currentTime + int64(expire),
		})
		break

	case "ams": // 运营平台
		expire = 12 * 3600
		refreshTokenExpire = 24 * 3600
		// Token
		tokenMethod = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"admin_id":         data["admin_id"],
			"is_super":         data["is_super"],
			"user_type":        "ams",
			"is_refresh_token": false,
			"token_type":       "biz", // 业务token
			"iss":              JwtAuthIdentifier,
			"aud":              JwtAuthIdentifier,
			"iat":              currentTime,
			"exp":              currentTime + int64(expire),
		})

		refreshTokenMethod = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"admin_id":         data["admin_id"],
			"is_super":         data["is_super"],
			"user_type":        "ams",
			"is_refresh_token": true,
			"token_type":       "biz", // 业务token
			"iss":              JwtAuthIdentifier,
			"aud":              JwtAuthIdentifier,
			"iat":              currentTime,
			"exp":              currentTime + int64(expire),
		})
	default:
		return Token{}, errors.New("用户类型错误")
	}

	if tokenMethod == nil {
		log.Println("tokenMethod nil")
		return Token{}, errors.New("token生成错误")
	}
	if refreshTokenMethod == nil {
		log.Println("refreshTokenMethod nil")
		return Token{}, errors.New("token生成错误")
	}

	// token
	var tokenSecret = []byte(config.App.Jwt.Secret)
	token, err := tokenMethod.SignedString(tokenSecret)
	if err != nil {
		return Token{}, errors.New("TOKEN生成错误")
	}

	// refreshToken
	refreshToken, err := refreshTokenMethod.SignedString(tokenSecret)
	if err != nil {
		return Token{}, errors.New("TOKEN生成错误")
	}

	t := Token{
		Exists:             1,
		Token:              token,
		Expire:             expire,
		RefreshToken:       refreshToken,
		RefreshTokenExpire: refreshTokenExpire,
	}

	if len(attach) > 0 {
		t.Attach = attach
	}

	return t, nil
}

// VerifyToken 验证用户token
func VerifyToken(token string) (map[string]interface{}, error) {
	jwtToken, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", jwtToken.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(config.App.Jwt.Secret), nil
	})
	if err != nil {
		log.Printf("%+v", err)
		errMsg := "token错误"
		if errors.Is(err, jwt.ErrTokenExpired) {
			errMsg = "token已过期"
		}
		return nil, errors.New(errMsg)
	}

	var ok bool
	var claims jwt.MapClaims
	if claims, ok = jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		if claims["aud"] != JwtAuthIdentifier || claims["iss"] != JwtAuthIdentifier {
			return nil, errors.New("illegal token")
		}
		if claims["token_type"] != "biz" {
			return nil, errors.New("token类型错误")
		}
	} else {
		return nil, errors.New("illegal token")
	}

	return claims, nil
}
