package response

import "github.com/dysodeng/app/internal/domain/file/model"

type InitMultipartUploadResponse struct {
	UploadId string `json:"upload_id"`
	Path     string `json:"path"`
}

type Part struct {
	PartNumber int    `json:"part_number"`
	ETag       string `json:"etag"`
	Size       int64  `json:"size"`
}

type PartList []Part

func PartListFormDomainModel(parts []model.Part) []Part {
	result := make([]Part, len(parts))
	for i, part := range parts {
		result[i] = Part{
			PartNumber: part.PartNumber,
			ETag:       part.ETag,
			Size:       part.Size,
		}
	}
	return result
}

type MultipartUploadStatusResponse struct {
	Parts []Part `json:"parts"`
	Path  string `json:"path"`
}
