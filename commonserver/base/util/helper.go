// Package util .

package util

import (
	"crypto/md5"
	"fmt"
	"strings"
	"time"
)

// NowSecTs 秒级时间戳
func NowSecTs() int64 {
	return time.Now().UTC().UnixNano() / 1e9
}

// NowTs 毫秒级时间戳
func NowTs() int64 {
	return time.Now().UTC().UnixNano() / 1e6
}

// NowUsTs 微秒级时间戳
func NowUsTs() int64 {
	return time.Now().UTC().UnixNano() / 1e3
}

// SysTime2Ts 系统时间转化为 ms时间戳
func SysTime2Ts(t time.Time) int64 {
	return t.UTC().UnixNano() / 1e6
}

// TodayStartTs 取得今天零点的时间戳（毫秒）
func TodayStartTs() int64 {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.UTC().Location())
	return todayStart.UnixNano() / 1e6
}

// TodayStartTs 取得UTC今天零点的时间戳（毫秒）
func TodayStartTsUTC() int64 {
	now := time.Now().UTC()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return todayStart.UnixNano() / 1e6
}

// NextDayTs 获取相对于传入时间的下一个UTC零点的时间戳
func NextDayTs(nowMs int64) int64 {
	now := time.Unix(0, nowMs*1e6).UTC()
	nowDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	nextDay := nowDay.AddDate(0, 0, 1)
	return nextDay.UTC().UnixNano() / 1e6
}

// Ms2Day 将得到的时间戳转化到天
func Ms2Day(t int64) int64 {
	return t / 1000 / 3600 / 24
}

// Ms2Sec 毫秒 -> 秒
func Ms2Sec(t int64) int64 {
	return t / 1000
}

// Ms2Min 毫秒 -> 分
func Ms2Min(t int64) int64 {
	return t / 1000 / 60
}

// Ms2Hour 毫秒 -> 小时
func Ms2Hour(dura int64) float32 {
	return float32(dura) / 1000 / 3600
}

// Sec2Ms 秒 -> 毫秒
func Sec2Ms(t int32) int32 {
	return t * 1000
}

// Min2Ms 分钟 -> 毫秒
func Min2Ms(t int32) int32 {
	return t * 1000 * 60
}

// Hour2Ms 小时 -> 毫秒
func Hour2Ms(t int32) int32 {
	return t * 3600 * 1000
}

// GetNowWeekDay 当前星期几
func GetNowWeekDay() time.Weekday {
	return time.Now().UTC().Weekday()
}

// Assert 强制断言
func Assert(cond bool) {
	if !cond {
		panic("assert failed")
	}
}

// IDirty 脏标记接口
type IDirty interface {
	MakeDirty()
	IsDirty() bool
	CleanDirty()
}

// Dirty 脏标记
// +k8s:deepcopy-gen=true
type Dirty struct {
	dirty bool
}

// MakeDirty 标记
func (d *Dirty) MakeDirty() {
	d.dirty = true
}

// IsDirty 判断标记
func (d *Dirty) IsDirty() bool {
	return d.dirty
}

// CleanDirty 清除标记
func (d *Dirty) CleanDirty() {
	d.dirty = false
}

// Md5 md5算法
func Md5(str string) string {
	data := []byte(str)
	s := fmt.Sprintf("%x", md5.Sum(data))
	return strings.ToUpper(s)
}
