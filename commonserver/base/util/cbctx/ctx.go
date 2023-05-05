// Package cbctx .
 
package cbctx

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

const (
	int32Tag  = "i32"
	int64Tag  = "i64"
	stringTag = "str"
	boolTag   = "bol"
	objectTag = "obj"
)

// SValue 所有简单值都会编码为字符串
type SValue string

// ---------- Decode ----------

// get tag
func (s SValue) tag() string {
	return string(s)[0:3]
}

// get value
func (s SValue) value() string {
	return string(s)[3:]
}

// check tag
func (s SValue) checkTag(tag string) {
	if s.tag() != tag {
		panic(fmt.Sprintf("svalue %s tag mismatch %v", string(s), tag))
	}
}

// Int32 Decode int32
func (s SValue) Int32() int32 {
	s.checkTag(int32Tag)
	i, err := strconv.Atoi(s.value())
	if err != nil {
		panic(fmt.Sprintf("svalue %v atoi failed", string(s)))
	}
	return int32(i)
}

// Int64 Decode int64
func (s SValue) Int64() int64 {
	s.checkTag(int64Tag)
	i, err := strconv.Atoi(s.value())
	if err != nil {
		panic(fmt.Sprintf("svalue %v atoi failed", string(s)))
	}
	return int64(i)
}

// Str Decode string
func (s SValue) Str() string {
	s.checkTag(stringTag)
	return s.value()
}

// Bool Decode bool
func (s SValue) Bool() bool {
	s.checkTag(boolTag)
	return s.value() == "t"
}

// Obj Decode bool
func (s SValue) Obj(v interface{}) {
	s.checkTag(objectTag)
	err := json.Unmarshal([]byte(s.value()), v)
	if err != nil {
		panic(fmt.Sprintf("svalue %v unmarshal %v failed", string(s), reflect.TypeOf(v)))
	}
}

// ---------- Encode ----------

// Int32 Encode int32
func Int32(i int32) SValue {
	return SValue(int32Tag + strconv.Itoa(int(i)))
}

// Int64 Encode int64
func Int64(i int64) SValue {
	return SValue(int64Tag + strconv.Itoa(int(i)))
}

// Str Encode string
func Str(s string) SValue {
	return SValue(stringTag + s)
}

// Bool Encode bool
func Bool(b bool) SValue {
	var v string
	if b {
		v = "t"
	} else {
		v = "f"
	}
	return SValue(boolTag + v)
}

// Obj Encode
func Obj(v interface{}) SValue {
	jsnbin, _ := json.Marshal(v)
	return SValue(objectTag + string(jsnbin))
}

// ---------- M ----------

// M Key为string, Value为简单值类型的Map
type M map[string]SValue

// Int32 取出一个int32值
func (m M) Int32(k string) int32 {
	return m[k].Int32()
}

// Int64 取出一个int64值
func (m M) Int64(k string) int64 {
	return m[k].Int64()
}

// Str 取出一个str值
func (m M) Str(k string) string {
	return m[k].Str()
}

// Bool 取出一个bool值
func (m M) Bool(k string) bool {
	return m[k].Bool()
}

// Clone 克隆
func (m M) Clone() (clone M) {
	if m == nil {
		return
	}

	clone = make(M)
	for key, value := range m {
		clone[key] = value
	}

	return clone
}
