package model

type StatefulConfig struct {
	//可能存在的root 用户
	RootUser string `json:"root_user"`
	//可能存在的root 密码
	RootPwd string `json:"root_pwd"`
	//可能存在的普通用户
	User string `json:"user"`
	//普通用户的密码
	Pwd string `json:"pwd"`
	//预置数据库名称
	Database string `json:"Database"`
}