// Package stat .

package stat

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// MsgStatSortTyp 消息统计排序类型
type MsgStatSortTyp int

// 消息统计排序类型
const (
	MsgStatSortCnt = 0
	MsgStatSortAvg = 1
	MsgStatSortMax = 2
	MsgStatSortPct = 3
)

// MsgStat 统计消息执行时间 groutine safe todo Ryan 增加开关
type MsgStat struct {
	mu    sync.Mutex
	msgs  map[string]*msgstat
	total int64 // 总时间
}

// NewStat 创建Stat管理
func NewStat() *MsgStat {
	return &MsgStat{
		msgs: make(map[string]*msgstat),
	}
}

type msgstat struct {
	ID interface{}
	// 调用次数
	Cnt int64
	// 有些是合并发送的算多次Cnt
	SingleCnt int64
	// 消耗统计 单位微妙
	Avg int64
	// 有些是合并发送的算多次Avg
	SingleAvg int64
	Max       int64
	MaxTime   string
	Total     int64
	SingleMax int64

	// 占总时间百分比 dump时再计算
	Pct float32
}

func (s *MsgStat) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.total = 0
	s.msgs = make(map[string]*msgstat)
}

// Add 不统计id为nil的消息
func (s *MsgStat) Add(id string, cost int64) {
	s.MulAdd(id, cost, 1, cost)
}

// Add 不统计id为nil的消息
func (s *MsgStat) MulAdd(id string, cost int64, totalCnt int64, singleMax int64) {
	// 统计
	s.mu.Lock()
	defer s.mu.Unlock()

	v, found := s.msgs[id]
	s.total += cost
	if found {
		v.Total += cost
		v.Cnt++
		v.SingleCnt += totalCnt
		if v.SingleCnt <= 0 {
			logrus.Errorf("MsgStat.Add id %v Totalcnt %v <= 0", id, v.SingleCnt)
			v.SingleCnt = 1
		}
		if v.Cnt <= 0 {
			logrus.Errorf("MsgStat.Add id %v cnt %v <= 0", id, v.Cnt)
			v.Cnt = 1
		}
		v.Avg = v.Total / v.Cnt
		v.SingleAvg = v.Total / v.SingleCnt
		if cost > v.Max {
			v.Max = cost
			v.MaxTime = time.Now().Format("2006-01-02 15:04:05.999999")
		}
		if singleMax > v.SingleMax {
			v.SingleMax = singleMax
		}
	} else {
		s.msgs[id] = &msgstat{
			ID:        id,
			Cnt:       1,
			SingleCnt: totalCnt,
			Total:     cost,
			Avg:       cost,
			SingleAvg: cost / totalCnt,
			Max:       cost,
			MaxTime:   time.Now().Format("2006-01-02 15:04:05.999999"),
		}
	}
}

func (s *MsgStat) Total() int64 {
	return s.total
}

// Dump 将消息统计信息编码为字符串
func (s *MsgStat) Dump(n int, typ MsgStatSortTyp) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	msgs := make([]*msgstat, 0, len(s.msgs))
	for _, v := range s.msgs {
		v.Pct = float32(v.Total) / float32(s.total)
		msgs = append(msgs, v)
	}

	sort.Sort(MsgStatList{
		msgs: msgs,
		typ:  typ,
	})
	if n > len(msgs) || n <= 0 {
		n = len(msgs)
	}
	//b, _ := json.Marshal(msgs[:n])
	//return string(b)
	res := s.DumpString(msgs[:n])
	return res
}

// DumpDefault 获取前n个处理耗时最长的消息
func (s *MsgStat) DumpDefault(n int) string {
	var dump strings.Builder
	dump.WriteString("stat (avg):\n")
	dump.WriteString(s.Dump(n, MsgStatSortAvg))
	dump.WriteString("stat (pct):\n")
	dump.WriteString(s.Dump(n, MsgStatSortPct))
	return dump.String()
}

func (s *MsgStat) DumpString(msgs []*msgstat) string {
	var dump strings.Builder
	for _, msg := range msgs {
		dump.WriteString(fmt.Sprintf("\nID:%v, Cnt:%v, Avg:%v, Max:%v, Total:%v, Pct:%v, "+
			"MaxTime%v, SingleCnt:%v, SingleAvg:%v, SingleMax:%v", msg.ID, msg.Cnt, msg.Avg,
			msg.Max, msg.Total, msg.Pct, msg.MaxTime, msg.SingleCnt, msg.SingleAvg, msg.SingleMax))
	}
	return dump.String()
}

// MsgStatList 用于对消息按照平均处理时间排序
type MsgStatList struct {
	msgs []*msgstat
	typ  MsgStatSortTyp
}

// Len 长度
func (s MsgStatList) Len() int {
	return len(s.msgs)
}

// Swap 交换
func (s MsgStatList) Swap(i, j int) {
	s.msgs[i], s.msgs[j] = s.msgs[j], s.msgs[i]
}

// Less 比较
func (s MsgStatList) Less(i, j int) bool {
	switch s.typ {
	case MsgStatSortAvg:
		return s.msgs[i].Avg > s.msgs[j].Avg
	case MsgStatSortMax:
		return s.msgs[i].Max > s.msgs[j].Max
	case MsgStatSortPct:
		return s.msgs[i].Pct > s.msgs[j].Pct
	default: // MsgStatSortCnt
		return s.msgs[i].Cnt > s.msgs[j].Cnt
	}
}
