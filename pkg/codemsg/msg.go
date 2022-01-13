package codemsg

var msgFlags = map[int]string{
	Success:                         "成功",
	Error:                           "错误",

	ErrorNoToken:					 "token不能为空",
	ErrorToken:                      "错误的token",
	ErrorTokenCreate:                "token生成失败",
	ErrorTokenExpire:  			 	 "token已过期",
	ErrorTokenFormat:				 "token格式错误",

	ErrorInvalidParams:              "请求参数错误",

	ErrorUserNotExist:               "错误的用户名",
	ErrorUserWrongPass:              "错误的用户密码",
	ErrorUserStatus:                 "用户已禁用",
}

func GetMsg(code int) string {
	msg, ok := msgFlags[code]
	if ok {
		return msg
	}

	return msgFlags[Error]
}
