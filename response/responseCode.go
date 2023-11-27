package response

type RspCode int64

const (
	CodeSuccess         RspCode = 1000
	CodeInvalidParams   RspCode = 1001
	CodeUserExist       RspCode = 1002
	CodeUserNotExist    RspCode = 1003
	CodeInvalidPassword RspCode = 1004
	CodeServerBusy      RspCode = 1005

	CodeInvalidToken      RspCode = 1006
	CodeInvalidAuthFormat RspCode = 1007
	CodeNotLogin          RspCode = 1008
	CodeTokenGenerateFail RspCode = 1009

	CodeEmailSendFail RspCode = 1010
	CodeRegisterFail  RspCode = 1011

	CodeBlogNotExist RspCode = 1012
	CodeInvalidJson  RspCode = 1013
)

var msgFlags = map[RspCode]string{
	CodeSuccess:         "success",
	CodeInvalidParams:   "请求参数错误",
	CodeUserExist:       "用户名重复",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",

	CodeInvalidToken:      "无效的Token",
	CodeInvalidAuthFormat: "认证格式有误",
	CodeNotLogin:          "未登录",
	CodeTokenGenerateFail: "Token生成失败",

	CodeEmailSendFail: "邮件发送失败",
	CodeRegisterFail:  "注册失败",

	CodeBlogNotExist: "博客不存在",

	CodeInvalidJson: "json格式有误",
}

func (c RspCode) Msg() string {
	msg, ok := msgFlags[c]
	if ok {
		return msg
	}
	return msgFlags[CodeServerBusy]
}
