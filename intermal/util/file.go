package util

import (
	"encoding/json"
	"os"
)

var File = fileUtil{}

type fileUtil struct{}

func (*fileUtil) WriteJson(path string, jsonData any) error {
	bytes, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	if err = os.WriteFile(path, bytes, 0666); err != nil {
		return err
	}
	return nil
}
