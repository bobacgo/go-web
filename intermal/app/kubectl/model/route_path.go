package model

type RoutePath struct {
	PathName           string `json:"pathName"`
	BackendService     string `json:"backendService"`
	BackendServicePort int32  `json:"backendServicePort"`
}