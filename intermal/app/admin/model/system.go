package model

// LoginReq 登录请求
// 1.username + password
// 2.phone + password
// 3.email + password
// 4.phone + smsCode
// 5.第三方授权认证 TODO
type LoginReq struct {
	Username string `json:"username"`
	Phone    string `json:"phone"` // 手机号
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRsp struct {
	SysUser
	RoleInfo SimpleRole `json:"roleInfo"`
	Menus    []*SysMenu
}