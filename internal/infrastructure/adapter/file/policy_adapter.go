package file

import (
	domainPort "github.com/dysodeng/app/internal/domain/file/port"
	"github.com/dysodeng/app/internal/domain/file/valueobject"
	infraConfig "github.com/dysodeng/app/internal/infrastructure/config"
)

// PolicyAdapter 文件上传策略端口适配器
type PolicyAdapter struct{}

func NewFilePolicyAdapter() domainPort.FilePolicy {
	return &PolicyAdapter{}
}

func (a *PolicyAdapter) Allow(mediaType valueobject.MediaType) ([]string, int64) {
	switch mediaType {
	case valueobject.MediaTypeImage:
		return infraConfig.AmsFileAllow.Image.AllowMimeType, infraConfig.AmsFileAllow.Image.AllowCapacitySize.ToInt()
	case valueobject.MediaTypeAudio:
		return infraConfig.AmsFileAllow.Audio.AllowMimeType, infraConfig.AmsFileAllow.Audio.AllowCapacitySize.ToInt()
	case valueobject.MediaTypeVideo:
		return infraConfig.AmsFileAllow.Video.AllowMimeType, infraConfig.AmsFileAllow.Video.AllowCapacitySize.ToInt()
	case valueobject.MediaTypeDocument:
		return infraConfig.AmsFileAllow.Document.AllowMimeType, infraConfig.AmsFileAllow.Document.AllowCapacitySize.ToInt()
	case valueobject.MediaTypeCompressed:
		return infraConfig.AmsFileAllow.Compressed.AllowMimeType, infraConfig.AmsFileAllow.Compressed.AllowCapacitySize.ToInt()
	default:
		return nil, 0
	}
}
