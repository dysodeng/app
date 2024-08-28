package config

import (
	"github.com/joho/godotenv"
)

// Env 应用环境
type Env string

func (e Env) Valid() bool {
	return e == Dev || e == Test || e == Release
}

const (
	// Release 生产环境
	Release Env = "release"
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
}

func Load() {
	_ = godotenv.Load()
	appConfigLoad()
	databaseConfigLoad()
	cacheConfigLoad()

	env := App.Env
	if !env.Valid() {
		env = Dev
	}
}
