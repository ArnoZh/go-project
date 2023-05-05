// Package reflectx .
 
package reflectx

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

/*
 * 用于Debug的reflect小工具
 */

// SetV 设置对象的字段值
func SetV(obj interface{}, keys, value string) error {
	fields := strings.Split(keys, ".")
	rvalue := reflect.ValueOf(obj)
	findex := 0
	for findex < len(fields) {
		field := fields[findex]
		switch rvalue.Type().Kind() {
		case reflect.Struct:
			rvalue = rvalue.FieldByName(field)
			findex++
		case reflect.Map:
			if findex == len(fields)-1 {
				return setMap(rvalue, field, value)
			}
			k, err := convertStr(rvalue.Type().Key().Kind(), field)
			if err != nil {
				return err
			}
			rvalue = rvalue.MapIndex(k)
			findex++
		case reflect.Array, reflect.Slice:
			i, err := strconv.Atoi(field)
			if err != nil {
				return fmt.Errorf("invalid arr index: %v", field)
			}
			rvalue = rvalue.Index(i)
			findex++
		case reflect.Ptr:
			rvalue = rvalue.Elem()
		default:
			return fmt.Errorf("unsuport type: %s", rvalue.Type().Kind().String())
		}
	}
	if !rvalue.IsValid() || !rvalue.CanSet() {
		return fmt.Errorf("invalid field: %s", fields)
	}
	return setField(rvalue, value)
}

func setField(field reflect.Value, value string) error {
	switch field.Type().Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int32, reflect.Int64:
		i, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(i))
	}
	return nil
}

func setMap(m reflect.Value, key, value string) error {
	k, err := convertStr(m.Type().Key().Kind(), key)
	if err != nil {
		return err
	}
	v, err := convertStr(m.Type().Elem().Kind(), value)
	if err != nil {
		return err
	}
	m.SetMapIndex(k, v)
	return nil
}

/*
func setArr(m reflect.Value, i int, value string) error {
	v, err := convertStr(m.Type().Elem().Kind(), value)
	if err != nil {
		return err
	}
	return nil
}
*/

func convertStr(kind reflect.Kind, key string) (reflect.Value, error) {
	var k reflect.Value
	switch kind {
	case reflect.String:
		k = reflect.ValueOf(key)
	case reflect.Int32:
		i, err := strconv.Atoi(key)
		if err != nil {
			return k, err
		}
		k = reflect.ValueOf(int32(i))
	case reflect.Int:
		i, err := strconv.Atoi(key)
		if err != nil {
			return k, err
		}
		k = reflect.ValueOf(i)
	case reflect.Int64:
		i, err := strconv.Atoi(key)
		if err != nil {
			return k, err
		}
		k = reflect.ValueOf(int64(i))
	default:
		return k, fmt.Errorf("convertStr unsuport type: %v", kind.String())
	}
	return k, nil
}
