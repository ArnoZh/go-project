// Package util .
 
package util

import (
	"bytes"
	"strings"
	"unicode"
)

// UnderscoreName 驼峰式写法转为下划线写法
func UnderscoreName(name string) string {
	buffer := NewBuffer()
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.Append('_')
			}
			buffer.Append(unicode.ToLower(r))
		} else {
			buffer.Append(r)
		}
	}

	return buffer.String()
}

// CamelName 下划线写法转为驼峰写法
func CamelName(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

// FindString 查找切片中是否存在某个字符串并返回所在的索引值，若不存在返回-1
func FindString(aim []string, name string) int {
	for i, v := range aim {
		if v == name {
			return i
		}
	}
	return -1
}

// IsRepeated 是否有重复的字符串
func IsRepeated(a []string) bool {
	m := make(map[string]bool, len(a))
	for _, s := range a {
		m[s] = true
	}
	return len(a) != len(m)
}

// BKDRBytesHash Hash字节序列
func BKDRBytesHash(b []byte) uint32 {
	seed := uint32(131)
	hash := uint32(0)

	for _, v := range b {
		hash = hash*seed + uint32(v)
	}
	return hash
}

// BKDRHash Hash一个字符串
func BKDRHash(s string) uint32 {
	b := bytes.NewBufferString(s).Bytes()
	return BKDRBytesHash(b)
}
