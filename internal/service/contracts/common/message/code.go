package message

import (
	messageModel "github.com/dysodeng/app/internal/dal/model/common"
	"github.com/dysodeng/app/internal/service"
)

// CodeMessageServiceInterface 验证码消息服务
type CodeMessageServiceInterface interface {
	// SendValidCode 发送验证码消息
	// @param senderType sms-短信 email-邮件
	// @param bizType 自定义业务类型
	// @param userType 用户类型 user-终端用户
	// @param account 用户账号 senderType=sms时为手机号 senderType=email时为邮箱地址
	// @param userId 用户ID 对应于userType的用户ID
	SendValidCode(
		senderType messageModel.SenderType,
		bizType,
		userType,
		account string,
	) service.Error

	// VerifyValidCode 验证码验证
	// @param senderType sms-短信消息 email-邮件消息
	// @param bizType 自定义业务类型
	// @param account 用户账号 senderType=sms时为手机号 senderType=email时为邮箱地址
	// @param code 验证码
	VerifyValidCode(
		senderType messageModel.SenderType,
		bizType,
		account,
		code string,
	) service.Error
}
