// Package leakybuf .
 
package leakybuf

// LeakyBuf 固定大小的，不会被GC的[]byte pool
type LeakyBuf struct {
	bufSize  int // size of each buffer
	freeList chan []byte
}

// New 创建leakybuf
func New(n, bufSize int) *LeakyBuf {
	return &LeakyBuf{
		bufSize:  bufSize,
		freeList: make(chan []byte, n),
	}
}

// Get 取出一个[]byte
func (lb *LeakyBuf) Get() (b []byte) {
	select {
	case b = <-lb.freeList:
	default:
		b = make([]byte, 0, lb.bufSize)
	}
	return
}

// Put 归还一个[]byte
func (lb *LeakyBuf) Put(b []byte) {
	b = b[:0]
	select {
	case lb.freeList <- b:
	default:
	}
}
