package aliyun

import (
	"encoding/json"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/pkg/errors"

	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/shared/cache"
)

type TokenResult struct {
	ErrMsg string `json:"ErrMsg"`
	Token  Token  `json:"Token"`
}

type Token struct {
	UserId     string `json:"UserId"`
	Id         string `json:"Id"`
	ExpireTime int64  `json:"ExpireTime"`
}

const IsiTokenCacheKey = "aliyun_isi_token"

// GetToken 获取语音服务token
func GetToken() (*TokenResult, error) {

	c, err := cache.NewCache()
	if err != nil {
		return nil, err
	}

	if c.IsExist(IsiTokenCacheKey) {
		data, err := c.Get(IsiTokenCacheKey)
		if err == nil {
			var token Token
			_ = json.Unmarshal([]byte(data), &token)
			if token.Id != "" {
				return &TokenResult{
					Token: token,
				}, nil
			}
		}
	}

	credentialsProvider := credentials.NewStaticAKCredentialsProvider(
		config.GlobalConfig.ThirdParty.TTS.Provider.Aliyun.AccessKeyId,
		config.GlobalConfig.ThirdParty.TTS.Provider.Aliyun.AccessKeySecret,
	)
	client, err := sdk.NewClientWithOptions("cn-shanghai", sdk.NewConfig(), credentialsProvider)
	if err != nil {
		return nil, err
	}

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Domain = "nls-meta.cn-shanghai.aliyuncs.com"
	request.ApiName = "CreateToken"
	request.Version = "2019-02-28"
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return nil, err
	}

	var tr TokenResult
	err = json.Unmarshal(response.GetHttpContentBytes(), &tr)
	if err != nil {
		return nil, err
	}
	if tr.ErrMsg != "" {
		return nil, errors.New(tr.ErrMsg)
	}

	expireTime := tr.Token.ExpireTime - time.Now().Unix()
	tokenBytes, _ := json.Marshal(tr.Token)
	_ = c.Put(IsiTokenCacheKey, string(tokenBytes), time.Second*time.Duration(expireTime-10))

	return &tr, nil
}
