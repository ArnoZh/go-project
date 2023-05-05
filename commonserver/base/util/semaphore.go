// Package util .
 
package util

// Semaphore 模拟一个简单的信号量，用于控制并发
type Semaphore chan struct{}

// MakeSemaphore 创建一个容量为N的信号量
func MakeSemaphore(n int) Semaphore {
	return make(Semaphore, n)
}

// Acquire 获取信号量
func (s Semaphore) Acquire() {
	s <- struct{}{}
}

// Release 释放信号量
func (s Semaphore) Release() {
	<-s
}
