package file

import (
	"net/http"
	"strconv"

	fileReq "github.com/dysodeng/app/internal/api/http/dto/request/file"
	"github.com/dysodeng/app/internal/api/http/dto/response/api"
	"github.com/dysodeng/app/internal/api/http/validator"
	"github.com/dysodeng/app/internal/application/file/service"
	"github.com/dysodeng/app/internal/pkg/logger"
	"github.com/dysodeng/app/internal/pkg/telemetry/trace"
	"github.com/gin-gonic/gin"
)

// UploaderController 文件上传
type UploaderController struct {
	baseTraceSpanName string
	uploaderService   service.UploaderApplicationService
}

func NewUploaderController(uploaderService service.UploaderApplicationService) *UploaderController {
	return &UploaderController{
		baseTraceSpanName: "api.http.controller.file.FileUploaderController",
		uploaderService:   uploaderService,
	}
}

// UploadFile 上传文件
func (c *UploaderController) UploadFile(ctx *gin.Context) {
	spanCtx, span := trace.Tracer().Start(trace.Gin(ctx), c.baseTraceSpanName+".UploadFile")
	defer span.End()

	fileForm, header, err := ctx.Request.FormFile("file")
	if err != nil {
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
func (c *UploaderController) InitMultipartUpload(ctx *gin.Context) {
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
func (c *UploaderController) UploadPart(ctx *gin.Context) {
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
func (c *UploaderController) CompleteMultipartUpload(ctx *gin.Context) {
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
func (c *UploaderController) MultipartUploadStatus(ctx *gin.Context) {
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
