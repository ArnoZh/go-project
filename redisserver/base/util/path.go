// Package util .
 
package util

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// GetCurDir 获取当前目录
func GetCurDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logrus.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
