package provider

import (
	domainFilePort "github.com/dysodeng/app/internal/domain/file/port"
	domainSharedPort "github.com/dysodeng/app/internal/domain/shared/port"
	"github.com/dysodeng/app/internal/infrastructure/adapter/file"
	sharedAdapter "github.com/dysodeng/app/internal/infrastructure/adapter/shared"
	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/event"
	"github.com/dysodeng/app/internal/infrastructure/persistence/transactions"
	"github.com/dysodeng/app/internal/infrastructure/shared/storage"
)

// ProvideFileStoragePort 提供端口适配器：文件存储
func ProvideFileStoragePort(st *storage.Storage) domainFilePort.FileStorage {
	return file.NewFileStorageAdapter(st)
}

// ProvideFilePolicyPort 提供端口适配器：文件策略
func ProvideFilePolicyPort(cfg *config.Config) domainFilePort.FilePolicy {
	// 当前策略直接使用 AmsFileAllow 全局配置
	return file.NewFilePolicyAdapter()
}

// ProvideEventPublisherPort 提供端口适配器：事件发布
func ProvideEventPublisherPort(bus event.Bus) domainSharedPort.EventPublisher {
	return sharedAdapter.NewEventPublisherAdapter(bus)
}

// ProvideTransactionManagerPort 提供端口适配器：事务管理
func ProvideTransactionManagerPort(tx transactions.TransactionManager) domainSharedPort.TransactionManager {
	return sharedAdapter.NewTransactionManagerAdapter(tx)
}
