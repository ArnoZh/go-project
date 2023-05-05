// Package util .
 
package util

import (
	"bytes"
	"strconv"
)

// Buffer 内嵌bytes.Buffer，支持连写
type Buffer struct {
	*bytes.Buffer
}

// NewBuffer 创建Buffer
func NewBuffer() *Buffer {
	return &Buffer{Buffer: new(bytes.Buffer)}
}

// Append 追加数据对象
func (b *Buffer) Append(i interface{}) *Buffer {
	switch val := i.(type) {
	case int:
		b.append(strconv.Itoa(val))
	case int64:
		b.append(strconv.FormatInt(val, 10))
	case uint:
		b.append(strconv.FormatUint(uint64(val), 10))
	case uint64:
		b.append(strconv.FormatUint(val, 10))
	case string:
		b.append(val)
	case []byte:
		b.Write(val)
	case rune:
		b.WriteRune(val)
	}

	return b
}

func (b *Buffer) append(s string) *Buffer {
	b.WriteString(s)
	return b
}
