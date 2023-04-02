package model

const RegionCodeListResource = "https://www.mca.gov.cn/article/sj/xzqh/2020/20201201.html"

type RegionCode struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	ParentCode string `json:"parentCode"`
	Children   []*RegionCode
}