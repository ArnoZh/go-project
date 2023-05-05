package pb

import (
	"commonserver/base/conf"
	"commonserver/base/network/msgpool"
	"commonserver/base/network/protobuf"
)

var Processor = protobuf.NewProcessor(conf.LittleEndian, msgpool.MaxMsgLen, 4)

func Init() {
	Processor.Register((*TestMsgReq)(nil))
	Processor.Register((*TestMsgAck)(nil))
}
