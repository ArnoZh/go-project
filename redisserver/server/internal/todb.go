package internal

import "redisserver/base/chanrpc"

// 操作函数类型
type OpHandler func(ci *chanrpc.CallInfo, agent *DBAgent)

// 消息填充操作
type RedisDBOp struct {
	Msg     interface{} // 消息
	HashKey int64       // hash key
	Discard bool        // 能否被抛弃
	Handler OpHandler   // 处理函数
	Agent   *DBAgent    //db代理
}
