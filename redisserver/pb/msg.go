package pb

import (
	"redisserver/base/conf"
	"redisserver/base/network/msgpool"
	"redisserver/base/network/protobuf"
)

var Processor = protobuf.NewProcessor(conf.LittleEndian, msgpool.MaxMsgLen, 4)

func Init() {
	Processor.Register((*TestMsgReq)(nil))
	Processor.Register((*TestMsgAck)(nil))
}
