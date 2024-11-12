package debug

import (
	"net/http"

	"github.com/dysodeng/app/internal/api/grpc/proto"
	"github.com/dysodeng/app/internal/dal/model/common"
	"github.com/dysodeng/app/internal/pkg/db"
	"github.com/dysodeng/app/internal/pkg/helper"
	"github.com/dysodeng/app/internal/pkg/logger"
	"github.com/dysodeng/app/internal/pkg/token"
	"github.com/dysodeng/app/internal/pkg/trace"
	"github.com/dysodeng/app/internal/service/reply/api"
	"github.com/dysodeng/app/internal/service/rpc"
	"github.com/dysodeng/app/internal/service/rpc/user"
	"github.com/gin-gonic/gin"
)

func Token(ctx *gin.Context) {
	t, _ := token.GenerateToken("user", map[string]interface{}{
		"user_id": 1,
	}, nil)
	ctx.JSON(200, api.Success(ctx, t))
}

// GenRandomString 生成随机字符串
// @route GET /debug/random_string
func GenRandomString(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, api.Success(ctx, helper.RandomStringBytesMask(24)))
}

func GormLogger(ctx *gin.Context) {
	traceCtx := trace.New().NewSpan(ctx, "debug.GormLogger")

	var mailConfig common.MailConfig
	db.DB().WithContext(traceCtx).First(&mailConfig)
	var smsConfig common.SmsConfig
	db.DB().WithContext(traceCtx).Where("a=?", "b").First(&smsConfig)

	go func() {
		childTraceCtx := trace.New().NewSpan(traceCtx, "debug.GormLogger.child")
		logger.Debug(childTraceCtx, "child logger")
		logger.Error(childTraceCtx, "child logger")
	}()

	ctx.JSON(200, api.Success(traceCtx, mailConfig))
}

func User(ctx *gin.Context) {
	userService, err := user.Service()
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(ctx, err.Error(), api.CodeFail))
		return
	}

	userInfo, err := userService.Info(rpc.Ctx(ctx), &proto.UserInfoRequest{
		Id: 2,
	})
	if err != nil {
		err, _ = rpc.Error(err)
		ctx.JSON(http.StatusOK, api.Fail(ctx, err.Error(), api.CodeFail))
		return
	}

	ctx.JSON(http.StatusOK, api.Success(ctx, userInfo))
}

func ListUser(ctx *gin.Context) {
	traceCtx := trace.New().NewSpan(ctx, "debug.ListUser")
	logger.Debug(traceCtx, "获取用户列表接口", logger.Field{Key: "params", Value: proto.UserListRequest{
		PageNum:  1,
		PageSize: 10,
	}})
	userService, err := user.Service()
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(traceCtx, err.Error(), api.CodeFail))
		return
	}
	res, err := userService.ListUser(rpc.Ctx(traceCtx), &proto.UserListRequest{
		PageNum:  1,
		PageSize: 10,
	})
	if err != nil {
		err, _ = rpc.Error(err)
		ctx.JSON(http.StatusOK, api.Fail(traceCtx, err.Error(), api.CodeFail))
		return
	}

	ctx.JSON(http.StatusOK, api.Success(traceCtx, res))
}

func CreateUser(ctx *gin.Context) {
	userService, err := user.Service()
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(ctx, err.Error(), api.CodeFail))
		return
	}

	_, err = userService.CreateUser(rpc.Ctx(ctx), &proto.UserRequest{
		Telephone: "13730825687",
		Password:  "dysodeng@112",
		RealName:  "dysodeng",
		Nickname:  "丹枫",
		Avatar:    "https://minio.dysodeng.com/user/avatar.png",
		Gender:    1,
		Birthday:  "1999-01-01",
	})
	if err != nil {
		err, _ = rpc.Error(err)
		ctx.JSON(http.StatusOK, api.Fail(ctx, err.Error(), api.CodeFail))
		return
	}
	ctx.JSON(http.StatusOK, api.Success(ctx, true))
}
