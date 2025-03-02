package debug

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dysodeng/app/internal/pkg/cache"

	"github.com/dysodeng/app/internal/api/grpc/proto"
	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/dal/model/common"
	"github.com/dysodeng/app/internal/pkg/db"
	"github.com/dysodeng/app/internal/pkg/helper"
	"github.com/dysodeng/app/internal/pkg/logger"
	"github.com/dysodeng/app/internal/pkg/request"
	"github.com/dysodeng/app/internal/pkg/retry"
	"github.com/dysodeng/app/internal/pkg/telemetry/metrics"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/dysodeng/app/internal/pkg/token"
	"github.com/dysodeng/app/internal/service/reply/api"
	"github.com/dysodeng/app/internal/service/rpc"
	"github.com/dysodeng/app/internal/service/rpc/user"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
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

	ctx.JSON(http.StatusOK, api.Success(spanCtx, helper.RandomString(32, helper.ModeAlphanumeric)))
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
	flusher := request.NewFlusher(ctx.Writer, ctx.Writer)

	data, _ := json.Marshal(map[string]interface{}{
		"query":              "你好，现在几点了",
		"response_mode":      "streaming",
		"user":               "dds",
		"auto_generate_name": true,
		"inputs": map[string]interface{}{
			"call":     "王先生",
			"gender":   "男",
			"age":      "32",
			"patBedId": 8,
		},
	})

	var message string

	statusCode, err := request.StreamRequest(
		config.ThirdParty.Dify.Api+"/chat-messages", // api
		"POST",
		bytes.NewBuffer(data),
		func(chunk []byte) error {
			chunkString := string(chunk)
			if strings.HasPrefix(chunkString, "data: ") {
				chunkString = strings.Replace(chunkString, "data: ", "", 1)
			}

			if chunkString != "" && chunkString != "\n" && chunkString != "\n\n" {
				_, _ = flusher.WriterWithFlush("data: " + chunkString + "\n\n")
				var msg Message
				_ = json.Unmarshal([]byte(chunkString), &msg)
				if msg.Event == "agent_message" || msg.Event == "message" {
					message += msg.Answer
				}
			}
			return nil
		},
		request.WithTimeout(2*time.Minute),
		request.WithContext(spanCtx),
		request.WithStreamMaxBufferSize(1024*1024),
		request.WithHeader("Authorization", fmt.Sprintf("Bearer %s", config.ThirdParty.Dify.ChatAppKey)), // api key
		request.WithHeader("Content-Type", "application/json"),
		request.WithTracer("Trace-Id", "Span-Id"),
	)
	if err != nil {
		logger.Error(spanCtx, "请求错误", logger.Field{Key: "statusCode", Value: statusCode}, logger.ErrorField(err))
	} else {
		logger.Info(spanCtx, "done", logger.Field{Key: "statusCode", Value: statusCode})
		logger.Info(spanCtx, "完整消息", logger.Field{Key: "message", Value: message})
	}
}

type Message struct {
	Event  string `json:"event"`
	Answer string `json:"answer"`
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

// Retry 重试
func Retry(ctx *gin.Context) {
	spanCtx, span := trace.Tracer().Start(trace.Gin(ctx), "debug.Retry")
	defer span.End()

	i := 1
	waitTime := 3 * time.Second
	retry.Invoke(
		func() error {
			if i == 3 {
				return nil
			}
			i++
			return errors.New("发生错误了")
		},
		retry.WithRetryNum(5),
		retry.WithRetryWaitTime(waitTime), // 重试等待时间
		retry.WithRetryWaitTimeFunc(func(retryNum int) time.Duration { // 自定义重试等待时间，每次按重试次数递增
			return time.Duration(retryNum) * waitTime
		}),
	)

	ctx.JSON(http.StatusOK, api.Success(spanCtx, true))
}

// Cache 缓存
func Cache(ctx *gin.Context) {
	spanCtx, span := trace.Tracer().Start(trace.Gin(ctx), "debug.Cache")
	defer span.End()

	userId := ctx.Query("user_id")
	userID, _ := strconv.ParseUint(userId, 10, 64)
	if userID <= 0 {
		ctx.JSON(http.StatusOK, api.Fail(ctx, "缺少用户ID", api.CodeFail))
		return
	}

	cli, err := cache.NewCache()
	if err != nil {
		logger.Error(spanCtx, "cache error", logger.ErrorField(err))
		ctx.JSON(http.StatusOK, api.Fail(ctx, "内部错误", api.CodeFail))
		return
	}

	cacheKey := fmt.Sprintf("user_info_%d", userID)
	userCache, err := cli.Get(cacheKey)
	if err == nil {
		var userInfo *proto.UserResponse
		if err = json.Unmarshal(helper.StringToBytes(userCache), &userInfo); err != nil {
			logger.Error(spanCtx, "cache error", logger.ErrorField(err))
			ctx.JSON(http.StatusOK, api.Fail(ctx, "内部错误", api.CodeFail))
			return
		}
		ctx.JSON(http.StatusOK, api.Success(ctx, userInfo))
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
		return
	}

	userBytes, _ := json.Marshal(userInfo)
	_ = cli.Put(cacheKey, helper.BytesToString(userBytes), 1*time.Hour)

	ctx.JSON(http.StatusOK, api.Success(ctx, userInfo))
}
