package idcard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

const apiUrl = "https://jmidcardv1.market.alicloudapi.com/idcard/validate"

// 身份证二要素认证
type valid struct {
	appCode string
}

func NewValid(appCode string) *valid {
	return &valid{
		appCode: appCode,
	}
}

// Check 检查姓名与身份证号是否一致
func (v *valid) Check(realName, idCardNo string) (*Body, error) {
	body := bytes.NewBuffer([]byte(fmt.Sprintf("idCardNo=%s&name=%s", idCardNo, realName)))
	req, err := http.NewRequest("POST", apiUrl, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("APPCODE %s", v.appCode))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(response.Body)
		if b != nil {
			var failBody map[string]string
			_ = json.Unmarshal(b, &failBody)
			if msg, ok := failBody["msg"]; ok {
				return nil, errors.New(msg)
			}
		}
		log.Printf("http get error : uri=%v , statusCode=%v", apiUrl, response.StatusCode)
		return nil, errors.New("实名验证失败")
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var res responseBody
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}

	if res.Code != 200 {
		return nil, errors.New(res.Msg)
	}

	return &res.Data, nil
}

type Body struct {
	Result   int    `json:"result"`
	Desc     string `json:"desc"`
	Sex      string `json:"sex"`
	Birthday string `json:"birthday"`
	Address  string `json:"address"`
}

type responseBody struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	TaskNo string `json:"taskNo"`
	Data   Body   `json:"data"`
}
