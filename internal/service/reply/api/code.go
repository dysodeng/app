package api

// Code 响应码
type Code int

func (c Code) ToInt() int {
	return int(c)
}

func (c Code) ToInt64() int64 {
	return int64(c)
}

// 业务接口响应码
const (
	CodeOk             Code = 200 // 成功
	CodeFail           Code = 0   // 业务错误
	CodeAccountDisable Code = 1   // 账号已被禁止
)

// 系统异常接口响应码
const (
	CodeBadRequest          Code = 400 // 请求出错
	CodeUnauthorized        Code = 401 // 未授权
	CodeForbidden           Code = 403 // 无权限
	CodeNotFound            Code = 404 // 未找到
	CodeMethodNotAllowed    Code = 405 // 请求方法不允许
	CodeInternalServerError Code = 500 // 服务器内部错误
)

// 错误消息内容
const (
	ErrorBusy        string = "系统繁忙，请稍后再试"
	ErrorMissUid     string = "缺少用户ID"
	ErrorMissAdminId string = "缺少管理员ID"
	ErrForbidden     string = "当前账号没有操作权限"
)
