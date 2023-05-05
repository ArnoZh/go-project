package internal

import (
	"redisserver/base/chanrpc"
	"redisserver/base/util"
	"redisserver/pb"
)

// 消息具体处理
func handleTestMsgReq(agent *DBAgent, req0 interface{}) {
	defer util.PrintPanicStack()
	ci := &chanrpc.CallInfo{}
	ci.Req = req0.(*pb.TestMsgReq)
	GetRedisDBMgr().OnTestMsgReq(ci, agent)
}
