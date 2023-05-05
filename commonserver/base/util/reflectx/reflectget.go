// Package reflectx .
 
package reflectx

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// GetV 获取对象的字段interface
func GetV(obj interface{}, keys string) (interface{}, error) {
	if keys == "" {
		return obj, nil
	}
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
			k, err := convertStr(rvalue.Type().Key().Kind(), field)
			if err != nil {
				return nil, err
			}
			rvalue = rvalue.MapIndex(k)
			findex++
		case reflect.Array, reflect.Slice:
			i, err := strconv.Atoi(field)
			if err != nil {
				return nil, fmt.Errorf("invalid arr index: %v", field)
			}
			rvalue = rvalue.Index(i)
			findex++
		case reflect.Ptr:
			rvalue = rvalue.Elem()
		default:
			return nil, fmt.Errorf("unsuport type: %s", rvalue.Type().Kind().String())
		}
	}
	if !rvalue.IsValid() {
		return nil, fmt.Errorf("invalid field: %s", fields)
	}
	return rvalue.Interface(), nil
}
