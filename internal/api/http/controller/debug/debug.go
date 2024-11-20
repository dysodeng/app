package debug

import (
	"net/http"

	"github.com/dysodeng/app/internal/api/grpc/proto"
	"github.com/dysodeng/app/internal/dal/model/common"
	"github.com/dysodeng/app/internal/pkg/db"
	"github.com/dysodeng/app/internal/pkg/helper"
	"github.com/dysodeng/app/internal/pkg/logger"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/dysodeng/app/internal/pkg/token"
	"github.com/dysodeng/app/internal/service/reply/api"
	"github.com/dysodeng/app/internal/service/rpc"
	"github.com/dysodeng/app/internal/service/rpc/user"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
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
	spanCtx, span := trace.Tracer().Start(trace.Gin(ctx), "debug.GenRandomString")
	span.End()
	ctx.JSON(http.StatusOK, api.Success(spanCtx, helper.RandomStringBytesMask(24)))
}

func GormLogger(ctx *gin.Context) {
	spanCtx, span := trace.Tracer().Start(trace.Gin(ctx), "debug.GormLogger")
	defer span.End()

	span.SetStatus(codes.Ok, "ok")
	logger.Debug(spanCtx, "trace logger")

	var mailConfig common.MailConfig
	db.DB().WithContext(spanCtx).First(&mailConfig)
	var smsConfig common.SmsConfig
	db.DB().WithContext(spanCtx).Where("a=?", "b").First(&smsConfig)

	go func() {
		childSpanCtx, childSpan := trace.Tracer().Start(spanCtx, "debug.GormLogger.child")
		defer childSpan.End()
		logger.Debug(childSpanCtx, "child logger")
		logger.Error(childSpanCtx, "child logger")
	}()

	ctx.JSON(200, api.Success(ctx, mailConfig))
}

func User(ctx *gin.Context) {
	spanCtx, span := trace.Tracer().Start(trace.Gin(ctx), "debug.User")
	defer span.End()

	userService, err := user.Service(spanCtx)
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(ctx, err.Error(), api.CodeFail))
		return
	}

	userInfo, err := userService.Info(spanCtx, &proto.UserInfoRequest{
		Id: 1,
	})
	if err != nil {
		err, _ = rpc.Error(err)
		ctx.JSON(http.StatusOK, api.Fail(spanCtx, err.Error(), api.CodeFail))
		return
	}

	ctx.JSON(http.StatusOK, api.Success(ctx, userInfo))
}

func ListUser(ctx *gin.Context) {
	spanCtx, span := trace.Tracer().Start(trace.Gin(ctx), "debug.ListUser")
	defer span.End()

	logger.Debug(spanCtx, "获取用户列表接口", logger.Field{Key: "params", Value: proto.UserListRequest{
		PageNum:  1,
		PageSize: 10,
	}})
	userService, err := user.Service(spanCtx)
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(ctx, err.Error(), api.CodeFail))
		return
	}
	res, err := userService.ListUser(spanCtx, &proto.UserListRequest{
		PageNum:  1,
		PageSize: 10,
	})
	if err != nil {
		err, _ = rpc.Error(err)
		ctx.JSON(http.StatusOK, api.Fail(ctx, err.Error(), api.CodeFail))
		return
	}

	ctx.JSON(http.StatusOK, api.Success(ctx, res))
}

func CreateUser(ctx *gin.Context) {
	spanCtx, span := trace.Tracer().Start(trace.Gin(ctx), "debug.CreateUser")
	defer span.End()

	userService, err := user.Service(spanCtx)
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(ctx, err.Error(), api.CodeFail))
		return
	}

	_, err = userService.CreateUser(spanCtx, &proto.UserRequest{
		Telephone: "13011223344",
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
