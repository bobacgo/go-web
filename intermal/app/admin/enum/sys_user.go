package enum

// 说明
// 不能以零值开始设计枚举值，当枚举值为零值时，gorm 默认不更新该值
// 奇偶作为正反

type UserStatus uint8

const (
	UserStatusEnable UserStatus = iota + 1
	UserStatusDisable
	UserStatusOnline
	UserStatusOffline
)

// Name 返回枚举值的名称
func (u UserStatus) Name() string {
	if u < UserStatusEnable || u > UserStatusOffline {
		return "未知"
	}
	return [...]string{"启用", "禁用", "在线", "离线"}[u]
}

// 将枚举值转成字符串，便于输出
func (u UserStatus) String() string {
	return u.Name()
}

type UserGenders uint8

const (
	UserGendersFemale UserGenders = iota + 1
	UserGendersMale
	UserGendersUnknown
)

// Name 返回枚举值的名称
func (u UserGenders) Name() string {
	if u < UserGendersFemale || u >= UserGendersUnknown {
		return "保密"
	}
	return [...]string{"女", "男", "保密"}[u]
}

// 将枚举值转成字符串，便于输出
func (u UserGenders) String() string {
	return u.Name()
}
