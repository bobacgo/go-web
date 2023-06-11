package test

import (
	"fmt"
	"github.com/gogoclouds/gogo/pkg/util"
	"strings"
	"testing"
)

func TestPwd(t *testing.T) {
	pwd := "$2a$10$teXMlpbSN5lLLDtrWAEPpeSgW5kR981N5adoJSbNO1TVN57QyZ7dK$0HCsBXam"
	idx := strings.LastIndex(pwd, "$")
	fmt.Println(pwd[idx+1:], pwd[:idx])
}

func TestBcrypt(t *testing.T) {
	hash, salt := util.BcryptHash("abc@123")
	fmt.Println(salt, hash)
	verify := util.BcryptVerify(salt, hash, "abc@123")
	t.Logf("verify: %v\n", verify)
}
