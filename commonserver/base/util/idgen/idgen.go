// Package idgen .
 
package idgen

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

/*
 * int64 ID生成器
 * 编码规则:
 * 32位时间戳 + 12位BaseID + 20位自增ID
 * 支持每秒生成 2^20 > 104万(1048575) 个ID
 * 每个节点需要用集群唯一的 BaseID(最大值4095) 来初始化该 ID 生成器
 */

// baseID 初始化器BaseID
var baseID int64

// idCounter ID计数器
var idCounter = int64(0)

// Lua51MaxInt Lua number精度范围
const Lua51MaxInt = int64(9007199254740992)

// Init 初始化BaseID
func Init(id int32) {
	if id > 4095 {
		panic(fmt.Sprintf("BaseID %d need < 4095", id))
	}
	baseID = int64(id)
}

// NewID 获取一个唯一ID
func NewID() int64 {
	counter := atomic.AddInt64(&idCounter, 1)
	now := time.Now().UTC().Unix()
	return (now << 32) | (baseID << 20) | counter
}

// RandomLuaID 随机一个Lua number精度范围内的ID 测试使用
func RandomLuaID() int64 {
	id := rand.New(rand.NewSource(time.Now().Unix())).Int63n(Lua51MaxInt)
	return id
}

// NewStrID 获取一个唯一ID 字符串
func NewStrID() string {
	id := NewID()
	return strconv.FormatInt(id, 16)
}

var idCounter32 = int32(0)

// NewID32 32位自增ID 每次重启清0 有重复风险 线程安全
func NewID32() int32 {
	counter := atomic.AddInt32(&idCounter32, 1)
	return counter
}

// ID2Str 数字ID转换为字符串
func ID2Str(id int64) string {
	return strconv.FormatInt(id, 16)
}

// Str2ID 字符串ID转换为数字
func Str2ID(s string) int64 {
	i, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		logrus.Errorf("Str2ID failed: %v", s)
	}
	return i
}
