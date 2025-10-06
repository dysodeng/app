package errors

import (
	domainErrors "github.com/dysodeng/app/internal/domain/shared/errors"
)

// 文件领域错误码
const (
	CodeFileNotFound         = "FILE_NOT_FOUND"
	CodeFileNameEmpty        = "FILE_NAME_EMPTY"
	CodeFileIDEmpty          = "FILE_ID_EMPTY"
	CodeFilePathEmpty        = "FILE_PATH_EMPTY"
	CodeFileQueryFailed      = "FILE_QUERY_FAILED"
	CodeFileNameExists       = "FILE_NAME_EXISTS"
	CodeFileDeleteFailed     = "FILE_DELETE_FAILED"
	CodeFileUploadFailed     = "FILE_UPLOAD_FAILED"
	CodeFileCheckFailed      = "FILE_CHECK_FAILED"
	CodeFileInvalidType      = "FILE_INVALID_TYPE"
	CodeFileSizeExceeded     = "FILE_SIZE_EXCEEDED"
	CodeFileRecordSaveFailed = "FILE_RECORD_SAVE_FAILED"
)

// 分片上传错误码
const (
	CodeFileMultipartInitFailed     = "FILE_MULTIPART_INIT_FAILED"
	CodeFileMultipartUploadFailed   = "FILE_MULTIPART_UPLOAD_FAILED"
	CodeFileMultipartCompleteFailed = "FILE_MULTIPART_COMPLETE_FAILED"
	CodeFileMultipartStatusFailed   = "FILE_MULTIPART_STATUS_FAILED"
	CodeFileMultipartReadFailed     = "FILE_MULTIPART_READ_FAILED"
)

// 文件管理相关错误
var (
	ErrFileNotFound     = domainErrors.NewFileError(CodeFileNotFound, "文件不存在", nil)
	ErrFileNameEmpty    = domainErrors.NewFileError(CodeFileNameEmpty, "文件名不能为空", nil)
	ErrFileIDEmpty      = domainErrors.NewFileError(CodeFileIDEmpty, "文件ID不能为空", nil)
	ErrFilePathEmpty    = domainErrors.NewFileError(CodeFilePathEmpty, "文件路径不能为空", nil)
	ErrFileQueryFailed  = domainErrors.NewFileError(CodeFileQueryFailed, "文件查询失败", nil)
	ErrFileNameExists   = domainErrors.NewFileError(CodeFileNameExists, "已存在同名文件", nil)
	ErrFileDeleteFailed = domainErrors.NewFileError(CodeFileDeleteFailed, "文件删除失败", nil)
)

// 文件上传相关错误
var (
	// ErrFileUploadFailed 基础上传错误
	ErrFileUploadFailed = domainErrors.NewFileError(CodeFileUploadFailed, "文件上传失败", nil)
	ErrFileCheckFailed  = domainErrors.NewFileError(CodeFileCheckFailed, "文件检测失败", nil)

	// ErrMultipartInitFailed 分片上传相关错误
	ErrMultipartInitFailed     = domainErrors.NewFileError(CodeFileMultipartInitFailed, "分片上传初始化失败", nil)
	ErrMultipartUploadFailed   = domainErrors.NewFileError(CodeFileMultipartUploadFailed, "分片上传失败", nil)
	ErrMultipartCompleteFailed = domainErrors.NewFileError(CodeFileMultipartCompleteFailed, "分片上传完成失败", nil)
	ErrMultipartStatusFailed   = domainErrors.NewFileError(CodeFileMultipartStatusFailed, "分片上传状态查询失败", nil)
	ErrMultipartReadFailed     = domainErrors.NewFileError(CodeFileMultipartReadFailed, "文件分片读取失败", nil)

	// ErrFileInvalidType 文件类型和限制相关错误
	ErrFileInvalidType  = domainErrors.NewFileError(CodeFileInvalidType, "不支持的文件类型", nil)
	ErrFileSizeExceeded = domainErrors.NewFileError(CodeFileSizeExceeded, "文件大小超出限制", nil)

	// ErrFileRecordSaveFailed 存储相关错误
	ErrFileRecordSaveFailed = domainErrors.NewFileError(CodeFileRecordSaveFailed, "文件记录保存失败", nil)
)
