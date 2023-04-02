package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gogoclouds/go-web/intermal/systme/model"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

func main() {
	ReadFile()
}

func listToTree(list []*model.RegionCode) {
	count := 0
	for _, item := range list {
		switch {
		case strings.HasSuffix(item.Code, "0000"): // 省
			count++
			continue
		case strings.HasSuffix(item.Code, "00") && !strings.HasSuffix(item.Code, "0000"): // 市
			for _, v := range list {
				if strings.HasPrefix(item.Code, v.Code[:2]) {
					item.ParentCode = v.Code
					v.Children = append(v.Children, item)
					count++
					break
				}
			}
		default: // 县
			for _, v := range list {
				if strings.HasPrefix(item.Code, v.Code[:4]) {
					item.ParentCode = v.Code
					v.Children = append(v.Children, item)
					count++
					break
				}
			}
		}
	}
	fmt.Println("----------------", len(list), count)
	var oneLevelTree []model.RegionCode
	for _, code := range list {
		if code.ParentCode == "" {
			oneLevelTree = append(oneLevelTree, *code)
		}
	}
	WriteFile("./2020_region_code_tree.json", oneLevelTree)
}

func ReadFile() {
	f, err := os.Open("./2020年12月中华人民共和国县以上行政区划代码.log")
	defer f.Close()
	if err != nil {
		log.Fatalln(err)
	}

	scanner := bufio.NewScanner(f)
	regexpMun, _ := regexp.Compile(`(\d+)`)
	var regionCodes []*model.RegionCode
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		code := regexpMun.FindString(line)
		value := GetStrCn(line)
		if code == "" || value == "" {
			fmt.Println(code, value)
			continue
		}
		rc := model.RegionCode{
			Code: code,
			Name: value,
		}
		regionCodes = append(regionCodes, &rc)
	}

	sort.Slice(regionCodes, func(i, j int) bool {
		if regionCodes[i].Code < regionCodes[j].Code {
			return true
		}
		return false
	})
	WriteFile("./2020_region_code_list.json", regionCodes)
	listToTree(regionCodes)
}

// GetStrCn 中文提取
func GetStrCn(str string) (cnStr string) {
	r := []rune(str)
	var strSlice []string
	for i := 0; i < len(r); i++ {
		if r[i] <= 40869 && r[i] >= 19968 {
			cnStr = cnStr + string(r[i])
			strSlice = append(strSlice, cnStr)
		}
	}
	return
}

func WriteFile(path string, data any) {
	tree, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	if err = os.WriteFile(path, tree, 0666); err != nil {
		log.Fatalln(err)
	}
}