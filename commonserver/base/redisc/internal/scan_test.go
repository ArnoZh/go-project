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
	"fmt"
	"reflect"
	"testing"

	"github.com/go-redis/redis"
)

// dial wraps DialDefaultServer() with a more suitable function name for examples.
func dial() *redis.Client {
	return redis.NewClient(&redis.Options{Addr: ":6379"})
}

func ExampleScan() {
	var c = dial()
	defer c.Close()

	_, err := c.Do("HMSET", "album:1", "title", "Red", "rating", 5).Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = c.Do("HMSET", "album:2", "title", "Earthbound", "rating", 1).Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = c.Do("HMSET", "album:3", "title", "Beat").Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = c.Do("LPUSH", "albums", "1").Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = c.Do("LPUSH", "albums", "2").Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = c.Do("LPUSH", "albums", "3").Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	values, err := c.Do("SORT", "albums",
		"BY", "album:*->rating",
		"GET", "album:*->title",
		"GET", "album:*->rating").Result()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("sorts %#v\n", values)
}

type s0 struct {
	X  int
	Y  int `redis:"y"`
	Bt bool
}

type s1 struct {
	X  int    `redis:"-"`
	I  int    `redis:"i"`
	U  uint   `redis:"u"`
	S  string `redis:"s"`
	P  []byte `redis:"p"`
	B  bool   `redis:"b"`
	Bt bool
	Bf bool
	s0
}

var scanStructTests = []struct {
	title string
	reply []string
	value interface{}
}{
	{"basic",
		[]string{"i", "-1234", "u", "5678", "s", "hello", "p", "world", "b", "t", "Bt", "1", "Bf", "0", "X", "123", "y", "456"},
		&s1{I: -1234, U: 5678, S: "hello", P: []byte("world"), B: true, Bt: true, Bf: false, s0: s0{X: 123, Y: 456}},
	},
}

func TestScanStruct(t *testing.T) {
	for _, tt := range scanStructTests {

		var reply []interface{}
		for _, v := range tt.reply {
			reply = append(reply, []byte(v))
		}

		value := reflect.New(reflect.ValueOf(tt.value).Type().Elem())

		if err := ScanStruct(reply, value.Interface()); err != nil {
			t.Fatalf("ScanStruct(%s) returned error %v", tt.title, err)
		}

		if !reflect.DeepEqual(value.Interface(), tt.value) {
			t.Fatalf("ScanStruct(%s) returned %v, want %v", tt.title, value.Interface(), tt.value)
		}
	}
}

func TestBadScanStructArgs(t *testing.T) {
	x := []interface{}{"A", "b"}
	test := func(v interface{}) {
		if err := ScanStruct(x, v); err == nil {
			t.Errorf("Expect error for ScanStruct(%T, %T)", x, v)
		}
	}

	test(nil)

	var v0 *struct{}
	test(v0)

	var v1 int
	test(&v1)

	x = x[:1]
	v2 := struct{ A string }{}
	test(&v2)
}

var argsTests = []struct {
	title    string
	actual   Args
	expected Args
}{
	{"struct ptr",
		Args{}.AddFlat(&struct {
			I  int               `redis:"i"`
			U  uint              `redis:"u"`
			S  string            `redis:"s"`
			P  []byte            `redis:"p"`
			M  map[string]string `redis:"m"`
			Bt bool
			Bf bool
		}{
			-1234, 5678, "hello", []byte("world"), map[string]string{"hello": "world"}, true, false,
		}),
		Args{"i", int(-1234), "u", uint(5678), "s", "hello", "p", []byte("world"), "m",
			map[string]string{"hello": "world"}, "Bt", true, "Bf", false},
	},
	{"struct",
		Args{}.AddFlat(struct{ I int }{123}),
		Args{"I", 123},
	},
	{"slice",
		Args{}.Add(1).AddFlat([]string{"a", "b", "c"}).Add(2),
		Args{1, "a", "b", "c", 2},
	},
	{"struct omitempty",
		Args{}.AddFlat(&struct {
			I  int               `redis:"i,omitempty"`
			U  uint              `redis:"u,omitempty"`
			S  string            `redis:"s,omitempty"`
			P  []byte            `redis:"p,omitempty"`
			M  map[string]string `redis:"m,omitempty"`
			Bt bool              `redis:"Bt,omitempty"`
			Bf bool              `redis:"Bf,omitempty"`
		}{
			0, 0, "", []byte{}, map[string]string{}, true, false,
		}),
		Args{"Bt", true},
	},
}

func TestArgs(t *testing.T) {
	for _, tt := range argsTests {
		if !reflect.DeepEqual(tt.actual, tt.expected) {
			t.Fatalf("%s is %v, want %v", tt.title, tt.actual, tt.expected)
		}
	}
}
