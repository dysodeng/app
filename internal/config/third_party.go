package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// thirdParty 第三方服务
type thirdParty struct {
	Wx   wx   `mapstructure:"wx"`
	Dify dify `mapstructure:"dify"`
}

type wx struct {
	MiniProgram wxAccount `mapstructure:"mini_program"`
}

// wxAccount 微信公众账户
type wxAccount struct {
	AppId  string `mapstructure:"app_id"`
	Secret string `mapstructure:"secret"`
}

// Dify.AI
type dify struct {
	Api        string `mapstructure:"api"`
	ChatAppKey string `mapstructure:"chat_app_key"`
	ToolSecret string `mapstructure:"tool_secret"`
}

var ThirdParty *thirdParty

func loadThirdPartyConfig() {
	v := viper.New()
	v.AutomaticEnv()

	v.SetConfigName("third_party")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	d := v.Sub("third_party")
	_ = d.BindEnv("wx.mini_program.app_id", "WX_MINI_PROGRAM_APP_ID")
	_ = d.BindEnv("wx.mini_program.secret", "WX_MINI_PROGRAM_SECRET")
	_ = d.BindEnv("dify.api", "DIFY_API")
	_ = d.BindEnv("dify.chat_app_key", "DIFY_CHAT_APP_KEY")
	_ = d.BindEnv("dify.tool_secret", "DIFY_TOOL_SECRET")

	if err := d.Unmarshal(&ThirdParty); err != nil {
		panic(err)
	}

	log.Println("配置文件`configs/third_party.yaml`加载完成")

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		log.Println("配置文件`configs/third_party.yaml`已变更")
		d = v.Sub("third_party")
		_ = d.Unmarshal(&ThirdParty)
	})
}
