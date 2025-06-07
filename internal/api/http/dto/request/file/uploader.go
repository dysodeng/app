package file

import (
	"github.com/dysodeng/app/internal/application/file/dto/command"
)

type Part struct {
	PartNumber int    `json:"part_number"`
	ETag       string `json:"etag"`
	Size       int64  `json:"size"`
}

type PartList []Part

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

type InitMultipartUploadReq struct {
	Filename string `json:"filename" form:"filename" binding:"required" msg:"请选择上传文件"`
	FileSize int64  `json:"file_size" form:"file_size" binding:"required" msg:"缺少文件大小"`
}

// CompleteMultipartUploadReq 完成分片上传请求
type CompleteMultipartUploadReq struct {
	UploadID string `json:"upload_id" binding:"required" msg:"缺少上传ID"`
	Parts    []Part `json:"parts" binding:"required" msg:"缺少分片信息"`
}

type MultipartUploadStatusReq struct {
	UploadID string `json:"upload_id" binding:"required" msg:"缺少上传ID"`
}
