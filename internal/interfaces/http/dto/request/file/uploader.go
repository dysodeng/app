package file

import (
	"github.com/dysodeng/app/internal/application/file/dto/command"
)

// Part 文件分片
type Part struct {
	PartNumber int    `json:"part_number"`
	ETag       string `json:"etag"`
	Size       int64  `json:"size"`
}

// PartList 文件分片列表
type PartList []Part

// ToAppDTO 将request dto 转换为 应用层 dto
func (parts PartList) ToAppDTO() []command.Part {
	domainParts := make([]command.Part, len(parts))
	for i, part := range parts {
		domainParts[i] = command.Part{
			PartNumber: part.PartNumber,
			ETag:       part.ETag,
			Size:       part.Size,
		}
	}
	return domainParts
}

// InitMultipartUploadReq 初始化分片上传请求体
type InitMultipartUploadReq struct {
	Filename string `json:"filename" form:"filename" binding:"required" msg:"请选择上传文件"`
	FileSize int64  `json:"file_size" form:"file_size" binding:"required" msg:"缺少文件大小"`
}

// CompleteMultipartUploadReq 完成分片上传请求
type CompleteMultipartUploadReq struct {
	UploadID string `json:"upload_id" binding:"required" msg:"缺少上传ID"`
	Parts    []Part `json:"parts" binding:"required" msg:"缺少分片信息"`
}

// MultipartUploadStatusReq 查询分片上传状态请求体
type MultipartUploadStatusReq struct {
	UploadID string `json:"upload_id" binding:"required" msg:"缺少上传ID"`
}
