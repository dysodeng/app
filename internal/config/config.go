package config

import (
	"github.com/joho/godotenv"
)

// Env 应用环境
type Env string

func (e Env) Valid() bool {
	return e == Dev || e == Test || e == Prod
}

func (e Env) String() string {
	return string(e)
}

const (
	// Prod 生产环境
	Prod Env = "prod"
	// Test 测试环境
	Test Env = "test"
	// Dev 开发环境
	Dev Env = "dev"
)

const (
	VarPath  string = "var"
	LogPath         = VarPath + "/logs"
	TempPath        = VarPath + "/tmp"
)

func init() {
	App = &appConfig{}
	Server = &serverConfig{}
	Database = &database{}
	Cache = &cache{}
	Redis = &redis{}
	MessageQueue = &messageQueue{}
	Storage = &storage{}
	Monitor = &monitor{}
	ThirdParty = &thirdParty{}

	_ = godotenv.Load()

	loadAppConfig()
	loadDatabaseConfig()
	loadRedisConfig()
	loadCacheConfig()
	loadMessageQueueConfig()
	loadEtcdConfig()
	loadStorageConfig()
	loadMonitorConfig()
	loadThirdPartyConfig()

	env := App.Env
	if !env.Valid() {
		env = Dev
	}
}

func Load() {}
