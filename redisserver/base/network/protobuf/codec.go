package protobuf

import (
	"redisserver/base/network/msgpool"
)

// Codec 所有的方法必须 goroutine 安全
type Codec interface {
	Unmarshal(data []byte) (uint32, IMsg, error)
	// 普通Marshal 无内存复用
	Marshal(msg IMsg) ([]byte, error)
	// 编码到指定的buffer,外部通过buffer.DeRef或buffer.Release来回收复用
	MarshalToBuffer(msg IMsg, refCnt int32) (*msgpool.Buffer, error)
	// 只读 MsgId, 不 Unmarshal
	ReadMsgID(data []byte) (uint32, error)
	// 传入 msg pointer
	MsgID(msg IMsg) uint32
}
