package command

import "github.com/dysodeng/app/internal/domain/file/model"

type Part struct {
	PartNumber int    `json:"part_number"`
	ETag       string `json:"etag"`
	Size       int64  `json:"size"`
}

type PartList []Part

func (parts PartList) ToDomainModel() []model.Part {
	domainParts := make([]model.Part, len(parts))
	for i, part := range parts {
		domainParts[i] = model.Part{
			PartNumber: part.PartNumber,
			ETag:       part.ETag,
			Size:       part.Size,
		}
	}
	return domainParts
}
