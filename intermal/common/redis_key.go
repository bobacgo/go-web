package common

// redis key
// 模块:主题:xxx:数据结构类型

const (
	RedisKeyCaptchaFmt       = "admin:captcha:%s:string"
	RedisKeyPasswdErrIncrFmt = "admin:passwdErrIncr:%s:string"
	RedisKeyJwt              = "admin:token:%s:string"
)
