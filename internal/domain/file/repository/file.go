package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/dysodeng/app/internal/domain/file/model"
	"github.com/dysodeng/app/internal/domain/file/valueobject"
)

// FileQuery 文件查询参数
type FileQuery struct {
	MediaType valueobject.MediaType // 媒体类型列表，可选
	Keyword   string                // 关键词搜索，可选
	StartTime *time.Time            // 开始时间，可选
	EndTime   *time.Time            // 结束时间，可选
	OrderBy   string                // 排序字段
	OrderType string                // 排序方式：asc/desc
	Page      int                   // 页码
	PageSize  int                   // 每页数量
	FileIDs   []uint64
}

// FileRepository 文件仓储接口
type FileRepository interface {
	// FindList 查询文件列表
	FindList(ctx context.Context, query FileQuery) ([]model.File, int64, error)
	// FindByID 根据ID获取文件信息
	FindByID(ctx context.Context, id uuid.UUID) (*model.File, error)
	// FindListByIds 根据文件id列表获取文件列表
	FindListByIds(ctx context.Context, ids []uuid.UUID) ([]model.File, error)
	// Save 保存文件记录
	Save(ctx context.Context, file *model.File) error
	// Delete 删除文件
	Delete(ctx context.Context, id uuid.UUID) error
	// BatchDelete 批量删除文件
	BatchDelete(ctx context.Context, ids []uuid.UUID) error
	// CheckFileNameExists 检查文件名是否已存在
	// name: 文件名
	// excludeId: 排除的文件ID（用于文件重命名时排除自身）
	CheckFileNameExists(ctx context.Context, name string, excludeId uuid.UUID) (bool, error)
}
