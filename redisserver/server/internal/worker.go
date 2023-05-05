// Package internal .

package internal

import (
	"redisserver/base/chanrpc"
	"redisserver/base/util"
	"time"

	"github.com/sirupsen/logrus"

	"redisserver/base/conf"
	"redisserver/base/redisc"
)

// workers ....

type worker struct {
	id       int
	client   *redisc.Client
	chanCall chan *chanrpc.CallInfo
	closeSig chan bool // 开关
	normal   bool      // 连接是否处于正常状态
	ops      int64     // 执行的消息数量
}

// NewWorker 创建worker
func newWorker(client *redisc.Client, chanLen, id int) *worker {
	w := &worker{
		id:       id,
		client:   client,
		normal:   true,
		closeSig: make(chan bool, 1),
		chanCall: make(chan *chanrpc.CallInfo, chanLen),
	}
	return w
}

// Run 启动worker
func (w *worker) Run() {
	// DB连接心跳
	heartbeat := time.NewTimer(conf.DBHeartbeat)
	for {
		select {
		case <-w.closeSig:
			heartbeat.Stop()
			w.flush()
			return
		case <-heartbeat.C:
			// 收到心跳信号 重置心跳
			heartbeat.Reset(conf.DBHeartbeat)
		case ci := <-w.chanCall:
			// 收到请求
			w.Op(ci, false)
		}
	}
}

// OnDestroy 准备停止
func (w *worker) OnDestroy() {
	w.closeSig <- true
	// worker 线程的收尾工作由worker自己完成
}

// flush 处理完所有任务
func (w *worker) flush() {
	logrus.Infof("db worker %d ops %v reqs remains %v", w.id, w.ops, len(w.chanCall))
	for {
		select {
		case ci := <-w.chanCall:
			w.Op(ci, true)
		default:
			logrus.Infof("db worker %d ops %v, worker chan call has not msg now", w.id, w.ops)
			return
		}
	}
}

// Op 执行消息
func (w *worker) Op(ci *chanrpc.CallInfo, Discard bool) {
	defer util.PrintPanicStack()
	msg := ci.Req.(*RedisDBOp)
	if Discard && msg.Discard {
		return
	}
	w.ops++
	msg.Handler(ci, msg.Agent)
	ci.Ret(nil)
}
