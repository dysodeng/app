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
	// Dev 开发环境
	Dev Env = "dev"
	// Test 测试环境
	Test Env = "test"
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
	MQ = &mq{}
	Filesystem = &filesystem{}
	Monitor = &monitor{}

	_ = godotenv.Load()

	appConfigLoad()
	databaseConfigLoad()
	redisConfigLoad()
	cacheConfigLoad()
	mqConfigLoad()
	etcdConfigLoad()
	filesystemConfigLoad()
	monitorConfigLoad()

	env := App.Env
	if !env.Valid() {
		env = Dev
	}
}

func Load() {}
