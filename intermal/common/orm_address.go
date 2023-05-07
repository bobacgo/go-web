package common

// Location 需要更新零值
type Location struct {
	Province *string `json:"province" gorm:"type:varchar(100)"`
	City     *string `json:"city" gorm:"type:varchar(100)"`
	District *string `json:"district" gorm:"type:varchar(100)"`
	Specific *string `json:"specific" gorm:"type:varchar(500)"`
}
