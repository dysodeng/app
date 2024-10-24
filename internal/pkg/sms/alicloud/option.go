package alicloud

import "encoding/json"

// option 短信配置
type option struct {
	phoneNumber   string
	signName      string
	templateId    string
	templateParam string
	accessKey     string
	accessSecret  string
}

// Config 短信主要配置
type Config func(*option)

// WithConfig 设置主要配置
// @param accessKey string 短信AppKey
// @param accessSecret string 短信SecretKey
// @param signName string 短信签名
func WithConfig(accessKey, accessSecret, signName string) Config {
	return func(o *option) {
		o.accessKey = accessKey
		o.accessSecret = accessSecret
		o.signName = signName
	}
}

// Option 短信额外配置
type Option func(*option)

// WithParams 设置短信参数
// @param params map[string]string 参数列表
func WithParams(params map[string]string) Option {
	return func(o *option) {
		if params != nil {
			p, err := json.Marshal(params)
			if err == nil {
				o.templateParam = string(p)
			}
		}
	}
}
