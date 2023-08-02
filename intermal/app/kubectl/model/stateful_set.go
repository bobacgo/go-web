package model

type StatefulSet struct {
	NamespaceWithName
	TypeID int64 `json:"typeId"`
	//中间件的端口
	Ports []PodPort `json:"ports"`
	//默认生成的账号密码
	Config StatefulConfig `json:"config"`
	//环境变量
	PodEnv []PodEnv `json:"podEnv"`

	CPUMax    string `gorm:"NOT_NULL;COMMENT:'cpu最大值'" json:"cpuMax"`     // CPU, in cores. (500m = .5 cores)
	CPUMin    string `gorm:"NOT_NULL;COMMENT:'cpu最小'" json:"cpuMin"`       // CPU, in cores. (500m = .5 cores)
	MemoryMax string `gorm:"NOT_NULL;COMMENT:'内存最大值'" json:"memoryMax"` // Memory, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	MemoryMin string `gorm:"NOT_NULL;COMMENT:'内存最小'" json:"memoryMin"`   // Memory, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)

	//中间件存储
	Storage []Storage `json:"storage"`
	//中间件副本
	Replicas int32  `json:"replicas"`
	Image    string `json:"image"`
}