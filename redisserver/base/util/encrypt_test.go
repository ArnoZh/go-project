package util

import (
	"fmt"
	"testing"
)

// 生成加密
func TestEncrypt(t *testing.T) {
	// "123456"- 7F5EB9374D5A7F796F4F4E730B58334B
	// "Mena2021()" 10B7C46AB971B096EDEAD4F4788F42F9
	pwd := "Mena2021()"
	epwd := EncryptPassword(pwd)
	fmt.Printf("pwd:%s after encrypt:%s\n", pwd, epwd)
}
