// Package conf .

package conf

import "time"

var (
	LenStackBuf        = 4096
	LittleEndian       = true
	TimerDispatcherLen = 20000
	DBHeartbeat        = 5 * time.Second // 心跳检查间隔，做断线重连
	DBDialTimeout      = 3 * time.Second
	DBQueueLen         = 5000
	WorkChanLen        = 5000
)

// 网络通信消息channel length
const (
	SingleClientRecvStreamLen = 10
	SingleClientSendMqLen     = 256
	AgentChanRPCLen           = 10
)

// redis key前缀设定
const (
	Sep            = ":"
	UserBasePrefix = "user"
)
