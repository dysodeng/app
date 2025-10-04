package model

import "github.com/dysodeng/app/internal/infrastructure/shared/errors"

// 文件管理相关错误
var (
	ErrFileNotFound     = &errors.DomainError{Code: "FILE_NOT_FOUND", Message: "文件不存在"}
	ErrFileNameEmpty    = &errors.DomainError{Code: "FILE_NAME_EMPTY", Message: "文件名不能为空"}
	ErrFileIDEmpty      = &errors.DomainError{Code: "FILE_NAME_EMPTY", Message: "文件ID不能为空"}
	ErrFilePathEmpty    = &errors.DomainError{Code: "FILE_PATH_EMPTY", Message: "文件路径不能为空"}
	ErrFileQueryFailed  = &errors.DomainError{Code: "FILE_QUERY_FAILED", Message: "文件查询失败"}
	ErrFileNameExists   = &errors.DomainError{Code: "FILE_NAME_EXISTS", Message: "已存在同名文件"}
	ErrFileDeleteFailed = &errors.DomainError{Code: "FILE_DELETE_FAILED", Message: "文件删除失败"}
)

// 文件上传相关错误
var (
	// ErrFileUploadFailed 基础上传错误
	ErrFileUploadFailed = &errors.DomainError{Code: "FILE_UPLOAD_FAILED", Message: "文件上传失败"}
	ErrFileCheckFailed  = &errors.DomainError{Code: "FILE_CHECK_FAILED", Message: "文件检测失败"}

	// ErrMultipartInitFailed 分片上传相关错误
	ErrMultipartInitFailed     = &errors.DomainError{Code: "MULTIPART_INIT_FAILED", Message: "分片上传初始化失败"}
	ErrMultipartUploadFailed   = &errors.DomainError{Code: "MULTIPART_UPLOAD_FAILED", Message: "分片上传失败"}
	ErrMultipartCompleteFailed = &errors.DomainError{Code: "MULTIPART_COMPLETE_FAILED", Message: "分片上传完成失败"}
	ErrMultipartStatusFailed   = &errors.DomainError{Code: "MULTIPART_STATUS_FAILED", Message: "分片上传状态查询失败"}
	ErrMultipartReadFailed     = &errors.DomainError{Code: "MULTIPART_READ_FAILED", Message: "文件分片读取失败"}

	// ErrFileInvalidType 文件类型和限制相关错误
	ErrFileInvalidType  = &errors.DomainError{Code: "FILE_INVALID_TYPE", Message: "不支持的文件类型"}
	ErrFileSizeExceeded = &errors.DomainError{Code: "FILE_SIZE_EXCEEDED", Message: "文件大小超出限制"}

	// ErrFileRecordSaveFailed 存储相关错误
	ErrFileRecordSaveFailed = &errors.DomainError{Code: "FILE_RECORD_SAVE_FAILED", Message: "文件记录保存失败"}
)
