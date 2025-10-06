package service

import (
	"context"
	"time"

	"github.com/dysodeng/wx/mini_program/auth"
	"github.com/google/uuid"

	"github.com/dysodeng/app/internal/application/passport/dto/command"
	"github.com/dysodeng/app/internal/application/passport/dto/response"
	passportErrors "github.com/dysodeng/app/internal/domain/passport/errors"
	"github.com/dysodeng/app/internal/domain/passport/model"
	"github.com/dysodeng/app/internal/domain/passport/valueobject"
	sharedErrors "github.com/dysodeng/app/internal/domain/shared/errors"
	userErrors "github.com/dysodeng/app/internal/domain/user/errors"
	userModel "github.com/dysodeng/app/internal/domain/user/model"
	"github.com/dysodeng/app/internal/domain/user/repository"
	"github.com/dysodeng/app/internal/domain/user/service"
	userVO "github.com/dysodeng/app/internal/domain/user/valueobject"
	"github.com/dysodeng/app/internal/infrastructure/shared/helper"
	"github.com/dysodeng/app/internal/infrastructure/shared/logger"
	"github.com/dysodeng/app/internal/infrastructure/shared/redis"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/trace"
	"github.com/dysodeng/app/internal/infrastructure/shared/token"
	"github.com/dysodeng/app/internal/infrastructure/shared/wx"
)

// PassportApplicationService 认证应用服务
type PassportApplicationService interface {
	Login(ctx context.Context, cmd *command.LoginCommand) (*response.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*response.LoginResponse, error)
	VerifyToken(ctx context.Context, cmd *command.VerifyTokenCommand) (map[string]interface{}, error)
}

type passportApplicationService struct {
	baseTraceSpanName string
	userRepository    repository.UserRepository
	userDomainService service.UserDomainService
}

func NewPassportApplicationService(userDomainService service.UserDomainService) PassportApplicationService {
	return &passportApplicationService{
		baseTraceSpanName: "application.passport.service.PassportApplicationService",
		userDomainService: userDomainService,
	}
}

func (svc *passportApplicationService) Login(ctx context.Context, cmd *command.LoginCommand) (*response.LoginResponse, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".Login")
	defer span.End()

	var data map[string]interface{}
	var attach map[string]interface{}

	switch cmd.UserType {
	case "user":
		info, err := svc.userLogin(spanCtx, cmd)
		if err != nil {
			return nil, err
		}
		if info == nil {
			return nil, passportErrors.ErrLoginFailed
		}
		if !info.Registered {
			return &response.LoginResponse{Registered: false}, nil
		}

		data = map[string]interface{}{
			"user_id":       info.UserId.String(),
			"platform_type": info.PlatformType.String(),
		}
		attach = map[string]interface{}{
			"nickname": info.Telephone,
			"avatar":   info.Avatar,
		}

	case "ams":
	default:
		return nil, passportErrors.ErrLoginUserTypeInvalid
	}

	tokenClaims, err := token.GenerateToken(cmd.UserType, data, attach)
	if err != nil {
		return nil, err
	}

	return &response.LoginResponse{
		Registered:         tokenClaims.Registered,
		Token:              tokenClaims.Token,
		Expire:             tokenClaims.Expire,
		RefreshToken:       tokenClaims.RefreshToken,
		RefreshTokenExpire: tokenClaims.RefreshTokenExpire,
		Attach:             tokenClaims.Attach,
	}, nil
}

func (svc *passportApplicationService) RefreshToken(ctx context.Context, refreshToken string) (*response.LoginResponse, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".RefreshToken")
	defer span.End()

	claims, err := token.VerifyToken(refreshToken)
	if err != nil {
		trace.Error(err, span)
		return nil, passportErrors.ErrTokenInvalid.Wrap(err)
	}

	if claims["is_refresh_token"] == false {
		return nil, passportErrors.ErrBizTokenCannotUsedForRefreshToken
	}

	var data map[string]interface{}
	var attach map[string]interface{}

	userType := claims["user_type"].(string)
	switch userType {
	case "user": // 用户
		platformType := claims["platform_type"].(string)
		userId := helper.IfaceConvertString(claims["user_id"])
		if userId == "" {
			return nil, userErrors.ErrUserInvalidInfo
		}
		uid, err := uuid.Parse(userId)
		if err != nil {
			logger.Error(spanCtx, userErrors.ErrUserInvalidInfo.Message, logger.ErrorField(err))
			return nil, userErrors.ErrUserInvalidInfo.Wrap(err)
		}

		user, err := svc.userDomainService.UserInfo(spanCtx, uid)
		if err != nil {
			return nil, err
		}
		if !user.Status.Bool() {
			return nil, userErrors.ErrUserDisabled
		}

		data = map[string]interface{}{
			"platform_type": platformType,
			"user_id":       userId,
		}

	case "ams": // 管理员

	default:
		return nil, passportErrors.ErrLoginUserTypeInvalid
	}

	tokenClaims, err := token.GenerateToken(userType, data, attach)
	if err != nil {
		logger.Error(spanCtx, "token生成失败", logger.ErrorField(err))
		return nil, err
	}

	return &response.LoginResponse{
		Registered:         tokenClaims.Registered,
		Token:              tokenClaims.Token,
		Expire:             tokenClaims.Expire,
		RefreshToken:       tokenClaims.RefreshToken,
		RefreshTokenExpire: tokenClaims.RefreshTokenExpire,
		Attach:             tokenClaims.Attach,
	}, nil
}

func (svc *passportApplicationService) VerifyToken(ctx context.Context, cmd *command.VerifyTokenCommand) (map[string]interface{}, error) {
	spanCtx, span := trace.Tracer().Start(ctx, svc.baseTraceSpanName+".VerifyToken")
	defer span.End()

	claims, err := token.VerifyToken(cmd.Token)
	if err != nil {
		trace.Error(err, span)
		return nil, passportErrors.ErrTokenInvalid.Wrap(err)
	}

	if claims["is_refresh_token"] == true {
		return nil, passportErrors.ErrRefreshTokenCannotUsedForBizToken
	}

	var data map[string]interface{}

	userType := claims["user_type"].(string)
	if userType != cmd.UserType {
		return nil, sharedErrors.ErrCommonUnauthorized
	}

	switch userType {
	case "user":
		platformType := claims["platform_type"].(string)
		userId := helper.IfaceConvertString(claims["user_id"])
		if userId == "" {
			return nil, userErrors.ErrUserInvalidInfo
		}
		uid, err := uuid.Parse(userId)
		if err != nil {
			return nil, userErrors.ErrUserInvalidInfo
		}

		user, err := svc.userDomainService.UserInfo(spanCtx, uid)
		if err != nil {
			return nil, err
		}
		if !user.Status.Bool() {
			return nil, userErrors.ErrUserDisabled
		}

		data = map[string]interface{}{
			"platform_type": platformType,
			"user_id":       userId,
		}

	case "ams":

	default:
		return nil, sharedErrors.ErrCommonUnauthorized
	}

	return data, nil
}

func (svc *passportApplicationService) userLogin(ctx context.Context, cmd *command.LoginCommand) (*model.UserLoginInfo, error) {
	var user *userModel.User
	var platformType valueobject.PlatformType

	switch cmd.GrantType {
	case "wx_code": // 微信小程序code静默登录
		_, openId, unionId, err := svc.getSessionKeyByCode(ctx, cmd.WxCode)
		if err != nil {
			return nil, err
		}
		if openId == "" {
			return nil, passportErrors.ErrPassportGetWxUserFailed
		}

		var userInfo *userModel.User
		if unionId != "" {
			userInfo, err = svc.userDomainService.FindByWxUnionId(ctx, unionId)
			if err != nil {
				return nil, passportErrors.ErrPassportGetWxUserFailed.Wrap(err)
			}
		}
		if userInfo == nil {
			userInfo, err = svc.userDomainService.FindByOpenId(ctx, valueobject.PlatformWxMinioProgram.String(), openId)
			if err != nil {
				return nil, passportErrors.ErrPassportGetWxUserFailed.Wrap(err)
			}
		}

		if userInfo == nil || userInfo.ID == uuid.Nil {
			return &model.UserLoginInfo{Registered: false}, nil
		}
		if !userInfo.Status.Bool() {
			return nil, userErrors.ErrUserDisabled
		}

		user = userInfo
		platformType = valueobject.PlatformWxMinioProgram

	case "wx_telephone": // 小程序授权手机号
		_, openId, unionId, err := svc.getSessionKeyByCode(ctx, cmd.WxCode)
		if err != nil {
			return nil, err
		}
		if openId == "" {
			return nil, passportErrors.ErrPassportGetWxUserFailed
		}

		phone, err := wx.MiniProgram().User().GetPhoneNumber(cmd.Code, openId)
		if err != nil {
			return nil, passportErrors.ErrPassportGetWxUserFailed.Wrap(err)
		}
		if phone == nil || phone.PurePhoneNumber == "" {
			return nil, userErrors.ErrUserWxTelephoneParsingFailed
		}

		userInfo, err := svc.userDomainService.FindByTelephone(ctx, phone.PurePhoneNumber)
		if err != nil {
			return nil, passportErrors.ErrPassportGetWxUserFailed.Wrap(err)
		}
		if userInfo == nil || userInfo.ID == uuid.Nil {
			// 注册
			userInfo, err = svc.userDomainService.Create(ctx, phone.PurePhoneNumber, unionId, openId, "")
			if err != nil {
				return nil, err
			}

			err = svc.userRepository.Save(ctx, userInfo)
			if err != nil {
				return nil, userErrors.ErrUserRegisterFailed.Wrap(err)
			}
		} else {
			if userInfo.WxMiniProgramOpenID.Value() != openId {
				return nil, userErrors.ErrUserTelephoneBound
			}
			// 更新微信绑定
			openIdVo, _ := userVO.NewWxMiniProgramOpenID(openId)
			unionidVo, _ := userVO.NewWxUnionID(unionId)
			userInfo.WxMiniProgramOpenID = openIdVo
			userInfo.WxUnionID = unionidVo
			_ = svc.userRepository.Save(ctx, userInfo)
		}

		user = userInfo
		platformType = valueobject.PlatformWxMinioProgram

	case "openid": // openid直接登录(测试使用)
		userInfo, err := svc.userDomainService.FindByOpenId(ctx, valueobject.PlatformWxMinioProgram.String(), cmd.OpenId)
		if err != nil {
			return nil, err
		}
		if userInfo == nil || userInfo.ID == uuid.Nil {
			return nil, userErrors.ErrUserNotFound
		}

		user = userInfo
		platformType = valueobject.PlatformWxMinioProgram

	default:
		return nil, passportErrors.ErrPassportUserGrantTypeInvalid
	}

	return &model.UserLoginInfo{
		Registered:   true,
		PlatformType: platformType,
		UserId:       user.ID,
		Telephone:    user.Telephone.Value(),
		Avatar:       user.Avatar.Value(),
	}, nil
}

// getSessionKeyByCode 根据wx.login的code获取session_key
func (svc *passportApplicationService) getSessionKeyByCode(ctx context.Context, code string) (sessionKey, openId, unionId string, err error) {
	cacheKey := redis.CacheKey("user:wx:login:code:" + code)
	cacheClient := redis.CacheClient()

	if cacheClient.Exists(ctx, cacheKey).Val() > 0 {
		cacheSessionKey := cacheClient.Get(ctx, cacheKey).Val()
		sessionCacheKey := "wx:mini_program_session_key:" + cacheSessionKey
		session := cacheClient.HGetAll(ctx, sessionCacheKey).Val()
		if o, ok := session["openid"]; ok {
			openId = o
		}
		if u, ok := session["union_id"]; ok {
			unionId = u
		}
		if s, ok := session["session_key"]; ok {
			sessionKey = s
		}
	} else {
		var session auth.Session
		session, err = wx.MiniProgram().Auth().Session(code)
		if err != nil {
			logger.Error(ctx, passportErrors.ErrPassportGetWxUserFailed.Message, logger.ErrorField(err))
			return "", "", "", passportErrors.ErrPassportGetWxUserFailed.Wrap(err)
		}

		openId = session.Openid
		unionId = session.UnionId
		sessionKey = session.SessionKey

		cacheSessionKey := uuid.NewString()
		sessionCacheKey := "wx:mini_program_session_key:" + cacheSessionKey

		cacheClient.HMSet(ctx, sessionCacheKey, map[string]string{
			"session_key": sessionKey,
			"openid":      openId,
			"union_id":    unionId,
		})
		cacheClient.Expire(ctx, sessionCacheKey, 30*time.Minute)
		cacheClient.Set(ctx, cacheKey, cacheSessionKey, 30*time.Minute)
	}

	return
}
