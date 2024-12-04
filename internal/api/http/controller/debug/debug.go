package debug

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dysodeng/app/internal/pkg/telemetry/metrics"
	"go.opentelemetry.io/otel/metric"

	"github.com/dysodeng/app/internal/api/grpc/proto"
	"github.com/dysodeng/app/internal/dal/model/common"
	"github.com/dysodeng/app/internal/pkg/db"
	"github.com/dysodeng/app/internal/pkg/helper"
	"github.com/dysodeng/app/internal/pkg/logger"
	"github.com/dysodeng/app/internal/pkg/request"
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

	userId := ctx.Query("user_id")
	userID, _ := strconv.ParseUint(userId, 10, 64)
	if userID <= 0 {
		ctx.JSON(http.StatusOK, api.Fail(ctx, "缺少用户ID", api.CodeFail))
		return
	}

	userService, err := user.Service(spanCtx)
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(ctx, err.Error(), api.CodeFail))
		return
	}

	userInfo, err := userService.Info(spanCtx, &proto.UserInfoRequest{
		Id: userID,
	})
	if err != nil {
		err, _ = rpc.Error(err)
		ctx.JSON(http.StatusOK, api.Fail(spanCtx, err.Error(), api.CodeFail))

		apiCounter, _ := metrics.Meter().Int64Counter(
			"user.fail",
			metric.WithDescription("获取用户信息失败数量"),
			metric.WithUnit("{call}"),
		)
		apiCounter.Add(spanCtx, 1)
		return
	}

	apiCounter, _ := metrics.Meter().Int64Counter(
		"user.success",
		metric.WithDescription("获取用户信息成功数量"),
		metric.WithUnit("{call}"),
	)
	apiCounter.Add(spanCtx, 1)

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

func ChatMessage(ctx *gin.Context) {
	spanCtx, span := trace.Tracer().Start(trace.Gin(ctx), "debug.Message")
	defer span.End()

	ctx.Writer.Header().Add("Content-Type", "text/event-stream; charset=utf-8")

	flusher, ok := ctx.Writer.(http.Flusher)
	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(map[string]interface{}{
		"query":              "你好，现在几点了",
		"response_mode":      "streaming",
		"user":               "dds",
		"auto_generate_name": true,
		"inputs": map[string]interface{}{
			"name":     "dds",
			"gender":   "男",
			"age":      "12",
			"patBedId": 8,
		},
	})
	statusCode, err := request.StreamRequest(
		"", // api
		"POST",
		bytes.NewBuffer(data),
		func(chunk []byte) error {
			chunkString := string(chunk)
			if strings.HasPrefix(chunkString, "data: ") {
				chunkString = strings.Replace(chunkString, "data: ", "", 1)
			}
			if chunkString != "" && chunkString != "\n" && chunkString != "\n\n" {
				_, _ = fmt.Fprintf(ctx.Writer, "data: "+chunkString+"\n\n")
				flusher.Flush()
			}
			return nil
		},
		request.WithTimeout(2*time.Minute),
		request.WithContext(spanCtx),
		request.WithStreamMaxBufferSize(1024*1024),
		request.WithHeader("Authorization", "Bearer "), // api key
		request.WithHeader("Content-Type", "application/json"),
		request.WithTracer("Trace-Id", "Span-Id"),
	)
	log.Printf("statusCode: %d, err: %v", statusCode, err)
	if err != nil {
		logger.Error(spanCtx, "done", logger.Field{Key: "statusCode", Value: statusCode}, logger.ErrorField(err))
	} else {
		logger.Info(spanCtx, "done", logger.Field{Key: "statusCode", Value: statusCode})
	}
}

func RemoteRequest(ctx *gin.Context) {
	spanCtx, span := trace.Tracer().Start(trace.Gin(ctx), "debug.RemoteRequest")
	defer span.End()

	userId := ctx.Query("user_id")

	body, statusCode, err := request.JsonRequest(
		"http://localhost:8080/api/v1/debug/grpc/user?user_id="+userId,
		"POST",
		nil,
		request.WithContext(spanCtx),
		request.WithTracer("Trace-Id", "Span-Id"),
	)
	if err != nil {
		logger.Error(spanCtx, "request error", logger.Field{Key: "error", Value: err})
		ctx.JSON(http.StatusOK, api.Fail(spanCtx, "接口请求失败", api.CodeFail))
		return
	}
	if statusCode != 200 {
		logger.Error(spanCtx, "request error", logger.Field{Key: "error", Value: string(body)}, logger.Field{Key: "status_code", Value: statusCode})
		ctx.JSON(http.StatusOK, api.Fail(spanCtx, "接口请求失败", api.CodeFail))
		return
	}

	var res map[string]interface{}
	_ = json.Unmarshal(body, &res)
	ctx.JSON(http.StatusOK, api.Success(spanCtx, res))
}
