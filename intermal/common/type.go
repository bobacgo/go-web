package common

import "time"

// 关联关系表命名规范
// otm 一对多
// mto 多对一
// mtm 多对多

type OrmModel struct {
	ID        string     `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}