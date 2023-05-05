// Package internal .
 
// Copyright 2012 Gary Burd
//
// Licensed under the Apache License, version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package internal

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// Error represents an error returned in a command reply.
type Error string

func (err Error) Error() string { return string(err) }

func ensureLen(d reflect.Value, n int) {
	if n > d.Cap() {
		d.Set(reflect.MakeSlice(d.Type(), n, n))
	} else {
		d.SetLen(n)
	}
}

func cannotConvert(d reflect.Value, s interface{}) error {
	var sname string
	switch s.(type) {
	case string:
		sname = "Client simple string"
	case Error:
		sname = "Client error"
	case int64:
		sname = "Client integer"
	case []byte:
		sname = "Client bulk string"
	case []interface{}:
		sname = "Client array"
	default:
		sname = reflect.TypeOf(s).String()
	}
	return fmt.Errorf("cannot convert from %s to %s", sname, d.Type())
}

func convertAssignBulkString(d reflect.Value, s string) (err error) {
	switch d.Type().Kind() {
	case reflect.Float32, reflect.Float64:
		var x float64
		x, err = strconv.ParseFloat(s, d.Type().Bits())
		d.SetFloat(x)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var x int64
		x, err = strconv.ParseInt(s, 10, d.Type().Bits())
		d.SetInt(x)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var x uint64
		x, err = strconv.ParseUint(s, 10, d.Type().Bits())
		d.SetUint(x)
	case reflect.Bool:
		var x bool
		x, err = strconv.ParseBool(s)
		d.SetBool(x)
	case reflect.String:
		d.SetString(s)
	case reflect.Slice:
		if d.Type().Elem().Kind() != reflect.Uint8 {
			err = cannotConvert(d, s)
		} else {
			d.SetBytes([]byte(s))
		}
	default:
		err = cannotConvert(d, s)
	}
	return
}

func convertAssignInt(d reflect.Value, s int64) (err error) {
	switch d.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		d.SetInt(s)
		if d.Int() != s {
			err = strconv.ErrRange
			d.SetInt(0)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if s < 0 {
			err = strconv.ErrRange
		} else {
			x := uint64(s)
			d.SetUint(x)
			if d.Uint() != x {
				err = strconv.ErrRange
				d.SetUint(0)
			}
		}
	case reflect.Bool:
		d.SetBool(s != 0)
	default:
		err = cannotConvert(d, s)
	}
	return
}

func convertAssignValue(d reflect.Value, s interface{}) (err error) {
	switch s := s.(type) {
	case string:
		err = convertAssignBulkString(d, s)
	case []byte:
		err = convertAssignBulkString(d, string(s))
	case int64:
		err = convertAssignInt(d, s)
	default:
		err = cannotConvert(d, s)
	}
	return err
}

func convertAssignArray(d reflect.Value, s []interface{}) error {
	if d.Type().Kind() != reflect.Slice {
		return cannotConvert(d, s)
	}
	ensureLen(d, len(s))
	for i := 0; i < len(s); i++ {
		if err := convertAssignValue(d.Index(i), s[i]); err != nil {
			return err
		}
	}
	return nil
}

func convertAssign(d, s interface{}) (err error) {
	// Handle the most common destination types using type switches and
	// fall back to reflection for all other types.
	switch s := s.(type) {
	case nil:
		// ingore
	case []byte:
		err = convertAssignByteSlice(d, s)
	case int64:
		err = convertAssignInt64(d, s)
	case string:
		switch d := d.(type) {
		case *string:
			*d = s
		default:
			err = cannotConvert(reflect.ValueOf(d), s)
		}
	case []interface{}:
		switch d := d.(type) {
		case *[]interface{}:
			*d = s
		case *interface{}:
			*d = s
		case nil:
			// skip value
		default:
			if d := reflect.ValueOf(d); d.Type().Kind() != reflect.Ptr {
				err = cannotConvert(d, s)
			} else {
				err = convertAssignArray(d.Elem(), s)
			}
		}
	case Error:
		err = s
	default:
		err = cannotConvert(reflect.ValueOf(d), s)
	}
	return err
}

func convertAssignByteSlice(d interface{}, s []byte) (err error) {
	switch d := d.(type) {
	case *string:
		*d = string(s)
	case *int:
		*d, err = strconv.Atoi(string(s))
	case *bool:
		*d, err = strconv.ParseBool(string(s))
	case *[]byte:
		*d = s
	case *interface{}:
		*d = s
	case nil:
		// skip value
	default:
		if d := reflect.ValueOf(d); d.Type().Kind() != reflect.Ptr {
			err = cannotConvert(d, s)
		} else {
			err = convertAssignBulkString(d.Elem(), string(s))
		}
	}

	return err
}

func convertAssignInt64(d interface{}, s int64) (err error) {
	switch d := d.(type) {
	case *int:
		x := int(s)
		if int64(x) != s {
			err = strconv.ErrRange
			x = 0
		}
		*d = x
	case *bool:
		*d = s != 0
	case *interface{}:
		*d = s
	case nil:
		// skip value
	default:
		if d := reflect.ValueOf(d); d.Type().Kind() != reflect.Ptr {
			err = cannotConvert(d, s)
		} else {
			err = convertAssignInt(d.Elem(), s)
		}
	}
	return err
}

// Scan copies from src to the values pointed at by dest.
//
// The values pointed at by dest must be an integer, float, boolean, string,
// []byte, interface{} or slices of these types. Scan uses the standard strconv
// package to convert bulk strings to numeric and boolean types.
//
// If a dest value is nil, then the corresponding src value is skipped.
//
// If a src element is nil, then the corresponding dest value is not modified.
//
// To enable easy use of Scan in a loop, Scan returns the slice of src
// following the copied values.
func Scan(src []interface{}, dest ...interface{}) ([]interface{}, error) {
	if len(src) < len(dest) {
		return nil, errors.New("san: array short")
	}
	var err error
	for i, d := range dest {
		err = convertAssign(d, src[i])
		if err != nil {
			err = fmt.Errorf("scan: cannot assign to dest %d: %v", i, err)
			break
		}
	}
	return src[len(dest):], err
}

type fieldSpec struct {
	name      string
	index     []int
	omitEmpty bool
}

type structSpec struct {
	m map[string]*fieldSpec
	l []*fieldSpec
}

func (ss *structSpec) fieldSpec(name []byte) *fieldSpec {
	return ss.m[string(name)]
}

func compileStructSpec(t reflect.Type, depth map[string]int, index []int, ss *structSpec) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		switch {
		case f.PkgPath != "" && !f.Anonymous:
			// Ignore unexported fields.
		case f.Anonymous:
			// TODO: Handle pointers. Requires change to decoder and
			// protection against infinite recursion.
			if f.Type.Kind() == reflect.Struct {
				compileStructSpec(f.Type, depth, append(index, i), ss)
			}
		default:
			fs := &fieldSpec{name: f.Name}
			tag := f.Tag.Get("redis")
			p := strings.Split(tag, ",")
			if len(p) > 0 {
				if p[0] == "-" {
					continue
				}
				if len(p[0]) > 0 {
					fs.name = p[0]
				}
				for _, s := range p[1:] {
					switch s {
					case "omitempty":
						fs.omitEmpty = true
					default:
						panic(fmt.Errorf("unknown field tag %s for type %s", s, t.Name()))
					}
				}
			}
			d, found := depth[fs.name]
			if !found {
				d = 1 << 30
			}
			switch {
			case len(index) == d:
				// At same depth, remove from result.
				delete(ss.m, fs.name)
				j := 0
				for k := 0; k < len(ss.l); k++ {
					if fs.name != ss.l[k].name {
						ss.l[j] = ss.l[k]
						j++
					}
				}
				ss.l = ss.l[:j]
			case len(index) < d:
				fs.index = make([]int, len(index)+1)
				copy(fs.index, index)
				fs.index[len(index)] = i
				depth[fs.name] = len(index)
				ss.m[fs.name] = fs
				ss.l = append(ss.l, fs)
			}
		}
	}
}

var (
	structSpecMutex sync.RWMutex
	structSpecCache = make(map[reflect.Type]*structSpec)
	// defaultFieldSpec = &fieldSpec{}
)

func structSpecForType(t reflect.Type) *structSpec {

	structSpecMutex.RLock()
	ss, found := structSpecCache[t]
	structSpecMutex.RUnlock()
	if found {
		return ss
	}

	structSpecMutex.Lock()
	defer structSpecMutex.Unlock()
	ss, found = structSpecCache[t]
	if found {
		return ss
	}

	ss = &structSpec{m: make(map[string]*fieldSpec)}
	compileStructSpec(t, make(map[string]int), nil, ss)
	structSpecCache[t] = ss
	return ss
}

var errscanStructValue = errors.New("scanStruct: value must be non-nil pointer to a struct")

// ScanStruct scans alternating names and values from src to a struct. The
// HGETALL and CONFIG GET commands return replies in this format.
//
// scanStruct uses exported field names to match values in the response. Use
// 'redis' field tag to override the name:
//
//      Field int `redis:"myName"`
//
// Fields with the tag redis:"-" are ignored.
//
// Integer, float, boolean, string and []byte fields are supported. Scan uses the
// standard strconv package to convert bulk string values to numeric and
// boolean types.
//
// If a src element is nil, then the corresponding field is not modified.
func ScanStruct(src []interface{}, dest interface{}) error {
	d := reflect.ValueOf(dest)
	if d.Kind() != reflect.Ptr || d.IsNil() {
		return errscanStructValue
	}
	d = d.Elem()
	if d.Kind() != reflect.Struct {
		return errscanStructValue
	}
	ss := structSpecForType(d.Type())

	if len(src)%2 != 0 {
		return errors.New("scanSlice: number of values not a multiple of 2")
	}

	for i := 0; i < len(src); i += 2 {
		s := src[i+1]
		if s == nil {
			continue
		}
		name, ok := src[i].([]byte)
		if !ok {
			return fmt.Errorf("scanStruct: key %d not a bulk string value", i)
		}
		fs := ss.fieldSpec(name)
		if fs == nil {
			continue
		}
		if err := convertAssignValue(d.FieldByIndex(fs.index), s); err != nil {
			return fmt.Errorf("scanStruct: cannot assign field %s: %v", fs.name, err)
		}
	}
	return nil
}

// ScanStructFromMap ..
func ScanStructFromMap(src map[string]string, dest interface{}) error {
	d := reflect.ValueOf(dest)
	if d.Kind() != reflect.Ptr || d.IsNil() {
		return errscanStructValue
	}
	d = d.Elem()
	if d.Kind() != reflect.Struct {
		return errscanStructValue
	}
	ss := structSpecForType(d.Type())

	for name, s := range src {
		fs := ss.fieldSpec([]byte(name))
		if fs == nil {
			continue
		}

		if err := convertAssignValue(d.FieldByIndex(fs.index), s); err != nil {
			return fmt.Errorf("scanSlice: cannot assign field %s: %v", fs.name, err)
		}
	}
	return nil
}

// Args is a helper for constructing command arguments from structured values.
type Args []interface{}

// Add returns the result of appending value to args.
func (args Args) Add(value ...interface{}) Args {
	return append(args, value...)
}

// AddFlat returns the result of appending the flattened value of v to args.
//
// Maps are flattened by appending the alternating keys and map values to args.
//
// Slices are flattened by appending the slice elements to args.
//
// Structs are flattened by appending the alternating names and values of
// exported fields to args. If v is a nil struct pointer, then nothing is
// appended. The 'redis' field tag overrides struct field names. See scanStruct
// for more information on the use of the 'redis' field tag.
//
// Other types are appended to args as is.
func (args Args) AddFlat(v interface{}) Args {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Struct:
		args = flattenStruct(args, rv)
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			args = append(args, rv.Index(i).Interface())
		}
	case reflect.Map:
		for _, k := range rv.MapKeys() {
			args = append(args, k.Interface(), rv.MapIndex(k).Interface())
		}
	case reflect.Ptr:
		if rv.Type().Elem().Kind() == reflect.Struct {
			if !rv.IsNil() {
				args = flattenStruct(args, rv.Elem())
			}
		} else {
			args = append(args, v)
		}
	default:
		args = append(args, v)
	}
	return args
}

func flattenStruct(args Args, v reflect.Value) Args {
	ss := structSpecForType(v.Type())
	for _, fs := range ss.l {
		fv := v.FieldByIndex(fs.index)
		if fs.omitEmpty {
			var empty = false
			switch fv.Kind() {
			case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
				empty = fv.Len() == 0
			case reflect.Bool:
				empty = !fv.Bool()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				empty = fv.Int() == 0
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				empty = fv.Uint() == 0
			case reflect.Float32, reflect.Float64:
				empty = fv.Float() == 0
			case reflect.Interface, reflect.Ptr:
				empty = fv.IsNil()
			}
			if empty {
				continue
			}
		}
		args = append(args, fs.name, fv.Interface())
	}
	return args
}
