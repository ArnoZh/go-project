// Package msgpool .
 
 
package msgpool

import (
	"sync"
	"sync/atomic"
)

// 网络消息长度由msgpool推导限制
const (
	MaxMsgLen = minMsgLen << maxLevel // 最大消息长度，256KB
	maxLevel  = 12
	minMsgLen = 64 // 最小梯度64字节
)

// MsgPool 网络消息池子
type MsgPool struct {
	pools [maxLevel + 1]sync.Pool
	sizes [maxLevel + 1]int
}

var msgPool MsgPool

func init() {
	for i := 0; i <= maxLevel; i++ {
		size := minMsgLen << uint(i)
		msgPool.pools[i] = sync.Pool{
			New: func() interface{} {
				return &Buffer{
					Bytes: make([]byte, 0, size),
				}
			},
		}
		msgPool.sizes[i] = size
	}
}

// Buffer 带引用计数的二进制Buffer
// Buffer 的所有方法都是goroutine safe的
type Buffer struct {
	Bytes  []byte
	refCnt int32
}

func getMsgLenLevel(msgLen int) int {
	for i := 0; i <= maxLevel; i++ {
		if msgLen <= msgPool.sizes[i] {
			return i
		}
	}
	return -1
}

// Get 获取 >= MsgLen 的一块 bytes
func Get(msgLen int, refCnt int32) *Buffer {
	i := getMsgLenLevel(msgLen)
	if i == -1 {
		return nil
	}

	buffer := msgPool.pools[i].Get().(*Buffer)
	buffer.Bytes = buffer.Bytes[:0]
	buffer.refCnt = refCnt
	return buffer
}

// DeRef 释放Buffer引用计数
func (b *Buffer) DeRef() {
	n := atomic.AddInt32(&b.refCnt, -1)
	if n == 0 {
		i := getMsgLenLevel(cap(b.Bytes))
		if i != -1 {
			msgPool.pools[i].Put(b)
		}
	}
}

// Release 强制回收Buffer，请确保之后不再引用
func (b *Buffer) Release() {
	atomic.StoreInt32(&b.refCnt, 0)
	i := getMsgLenLevel(cap(b.Bytes))
	if i != -1 {
		msgPool.pools[i].Put(b)
	}
}
