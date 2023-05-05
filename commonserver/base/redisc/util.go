// Package redisc .

package redisc

import (
	"errors"
	"fmt"
	"redisserver/base/conf"
	"redisserver/base/redisc/internal"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

// Check indicates that a reply value is nil.
func Check(err error) error {
	if err != nil && err != redis.Nil {
		return err
	}
	return nil
}

// StringMap is a helper that converts an array of strings (alternating key, value)
// into a map[string]string. The HGetAll and CONFIG GET commands return replies in this format.
// Requires an even number of values in result.
func StringMap(values []string, err error) (map[string]string, error) {
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("redis: StringMap expects even number of values result")
	}
	m := make(map[string]string, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, value := values[i], values[i+1]
		m[key] = value
	}
	return m, nil
}

// Int64Map is a helper that converts an array of strings (alternating key, value)
// into a map[string]int64. The HGetAll commands return replies in this format.
// Requires an even number of values in result.
func Int64Map(values []string, err error) (map[string]int64, error) {
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("redis: Int64Map expects even number of values result")
	}
	m := make(map[string]int64, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key := values[i]
		value, err := Int64(values[i+1], nil)
		if err != nil {
			return nil, err
		}
		m[key] = value
	}
	return m, nil
}

// Int64SliceWithScore is a helper that converts an array of strings (alternating key, value)
// into a []*ScorePair.
// Requires an even number of values in result.
// func Int64SliceWithScore(values []string, err error) ([]*ScorePair, error) {
//	if err != nil {
//		return nil, err
//	}
//	if len(values)%2 != 0 {
//		return nil, errors.New("redis: Int64Slice expects even number of values result")
//	}
//	m := make([]*ScorePair, 0, len(values)/2)
//	curIndex := 0
//	for i := 0; i < len(values); i += 2 {
//		key := values[i]
//		value, err := Int64(values[i+1], nil)
//		if err != nil {
//			return nil, err
//		}
//		m[curIndex] = &ScorePair{element: key, score: value}
//		curIndex++
//	}
//	return m, nil
// }
//

// Int64Slice ..
func Int64Slice(values []string, err error) ([]int64, error) {
	if err != nil {
		return nil, err
	}

	m := make([]int64, 0, len(values))
	for i := 0; i < len(values); i++ {
		value, err := Int64(values[i], nil)
		if err != nil {
			continue
		}
		m = append(m, value)
	}
	return m, nil
}

// Int64 is a helper that converts a command reply to 64 bit integer. If err is
// not equal to nil, then Int returns 0, err. Otherwise, Int64 converts the
// reply to an int64 as follows:
//
//	Reply type    Result
//	integer       reply, nil
//	bulk string   parsed reply, nil
//	nil           0, ErrNil
//	other         0, error
func Int64(reply interface{}, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	switch reply := reply.(type) {
	case int64:
		return reply, nil
	case string:
		n, err := strconv.ParseInt(reply, 10, 64)
		return n, err
	case []byte:
		n, err := strconv.ParseInt(string(reply), 10, 64)
		return n, err
	case nil:
		return 0, redis.Nil
	}
	return 0, fmt.Errorf("redigo: unexpected type for Int64, got type %T", reply)
}

// ConvertTo  ConvertTo
// other solution:
// "github.com/mitchellh/mapstructure"
// "github.com/fatih/structs"
func ConvertTo(fields map[string]string, val interface{}) (interface{}, error) {
	value := reflect.New(reflect.ValueOf(val).Type().Elem())

	err := internal.ScanStructFromMap(fields, value.Interface())
	return value.Interface(), err
}

// Flat  Structs are flattened by appending the alternating names and values of
// exported fields to args.
func Flat(v interface{}, args ...interface{}) internal.Args {
	argList := make(internal.Args, 0, 16)
	return argList.Add(args...).AddFlat(v)
}

// RankKey 生成对应服务器的Key
func RankKey(key string, serverID int32) string {
	return I64toa(int64(serverID)) + conf.Sep + key
}

// I64toa 转成string
func I64toa(i64 int64) string {
	return strconv.FormatInt(i64, 10)
}

// AtoI64 转成int64
func AtoI64(str string) int64 {
	i64, _ := strconv.ParseInt(str, 10, 64)
	return i64
}

// U64toa 转成string
func U64toa(u64 uint64) string {
	return strconv.FormatUint(u64, 10)
}

// NewKey 生成redis key
func NewKey(prefix string, keys ...interface{}) string {
	var builder strings.Builder
	builder.WriteString(prefix)
	for _, key := range keys {
		builder.WriteString(conf.Sep)
		switch key := key.(type) {
		case int:
			builder.WriteString(I64toa(int64(key)))
		case int32:
			builder.WriteString(I64toa(int64(key)))
		case int64:
			builder.WriteString(I64toa(key))
		case uint:
			builder.WriteString(U64toa(uint64(key)))
		case uint32:
			builder.WriteString(U64toa(uint64(key)))
		case uint64:
			builder.WriteString(U64toa(key))
		case []byte:
			builder.Write(key)
		case string:
			builder.WriteString(key)
		default:
			builder.WriteString(fmt.Sprint(key))
		}
	}

	return builder.String()
}
