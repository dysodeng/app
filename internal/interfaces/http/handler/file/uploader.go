package file

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/dysodeng/app/internal/application/file/service"
	"github.com/dysodeng/app/internal/infrastructure/shared/logger"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/trace"
	fileReq "github.com/dysodeng/app/internal/interfaces/http/dto/request/file"
	"github.com/dysodeng/app/internal/interfaces/http/dto/response/api"
	"github.com/dysodeng/app/internal/interfaces/http/validator"
)

// UploaderHandler 文件上传
type UploaderHandler struct {
	baseTraceSpanName string
	uploaderService   service.UploaderApplicationService
}

// NewUploaderHandler 创建文件上传控制器
func NewUploaderHandler(uploaderService service.UploaderApplicationService) *UploaderHandler {
	return &UploaderHandler{
		baseTraceSpanName: "api.http.controller.file.UploaderHandler",
		uploaderService:   uploaderService,
	}
}

// UploadFile 上传文件
func (c *UploaderHandler) UploadFile(ctx *gin.Context) {
	spanCtx, span := trace.Tracer().Start(trace.Gin(ctx), c.baseTraceSpanName+".UploadFile")
	defer span.End()

	fileForm, header, err := ctx.Request.FormFile("file")
	if err != nil {
		logger.Warn(ctx, "文件上传失败", logger.ErrorField(err))
		ctx.JSON(http.StatusOK, api.Fail(spanCtx, "文件上传失败", api.CodeFail))
		return
	}
	defer func() {
		_ = fileForm.Close()
	}()

	file, err := c.uploaderService.UploadFile(spanCtx, header)
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(spanCtx, err.Error(), api.CodeFail))
		return
	}

	ctx.JSON(http.StatusOK, api.Success(spanCtx, file))
}

// InitMultipartUpload 初始化分片上传
func (c *UploaderHandler) InitMultipartUpload(ctx *gin.Context) {
	spanCtx, span := trace.Tracer().Start(trace.Gin(ctx), c.baseTraceSpanName+".InitMultipartUpload")
	defer span.End()

	var req fileReq.InitMultipartUploadReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, api.Fail(spanCtx, validator.TransError(err), api.CodeFail))
		return
	}

	res, err := c.uploaderService.InitMultipartUpload(spanCtx, req.Filename, req.FileSize)
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(spanCtx, err.Error(), api.CodeFail))
		return
	}

	ctx.JSON(http.StatusOK, api.Success(spanCtx, res))
}

// UploadPart 上传分片
func (c *UploaderHandler) UploadPart(ctx *gin.Context) {
	spanCtx, span := trace.Tracer().Start(trace.Gin(ctx), c.baseTraceSpanName+".UploadPart")
	defer span.End()

	uploadID := ctx.PostForm("upload_id")
	filePath := ctx.PostForm("path")
	partNumberStr := ctx.PostForm("part_number")
	if uploadID == "" || filePath == "" || partNumberStr == "" {
		ctx.JSON(http.StatusOK, api.Fail(spanCtx, "参数不完整", api.CodeFail))
		return
	}

	partNumber, err := strconv.Atoi(partNumberStr)
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(spanCtx, "分片编号格式错误", api.CodeFail))
		return
	}

	fileForm, header, err := ctx.Request.FormFile("file")
	if err != nil {
		logger.Error(spanCtx, "文件分片读取失败", logger.ErrorField(err))
		ctx.JSON(http.StatusOK, api.Fail(spanCtx, "文件分片读取失败", api.CodeFail))
		return
	}
	defer func() {
		_ = fileForm.Close()
	}()

	res, err := c.uploaderService.UploadPart(spanCtx, filePath, uploadID, partNumber, header)
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(spanCtx, err.Error(), api.CodeFail))
		return
	}

	ctx.JSON(http.StatusOK, api.Success(spanCtx, res))
}

// CompleteMultipartUpload 完成分片上传
func (c *UploaderHandler) CompleteMultipartUpload(ctx *gin.Context) {
	spanCtx, span := trace.Tracer().Start(trace.Gin(ctx), c.baseTraceSpanName+".CompleteMultipartUpload")
	defer span.End()

	var req fileReq.CompleteMultipartUploadReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, api.Fail(spanCtx, validator.TransError(err), api.CodeFail))
		return
	}

	file, err := c.uploaderService.CompleteMultipartUpload(spanCtx, req.UploadID, fileReq.PartList(req.Parts).ToAppDTO())
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(spanCtx, err.Error(), api.CodeFail))
		return
	}

	ctx.JSON(http.StatusOK, api.Success(spanCtx, file))
}

// MultipartUploadStatus 查询分片上传状态
func (c *UploaderHandler) MultipartUploadStatus(ctx *gin.Context) {
	spanCtx, span := trace.Tracer().Start(trace.Gin(ctx), c.baseTraceSpanName+".MultipartUploadStatus")
	defer span.End()

	var req fileReq.MultipartUploadStatusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, api.Fail(spanCtx, validator.TransError(err), api.CodeFail))
		return
	}

	res, err := c.uploaderService.MultipartUploadStatus(spanCtx, req.UploadID)
	if err != nil {
		ctx.JSON(http.StatusOK, api.Fail(spanCtx, err.Error(), api.CodeFail))
		return
	}

	ctx.JSON(http.StatusOK, api.Success(spanCtx, res))
}
