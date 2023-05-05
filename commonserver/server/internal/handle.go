package internal

import (
	"commonserver/pb"
	"fmt"
	"reflect"
)

// 消息具体处理
func handleTestMsgAck(agent *DBAgent, req0 interface{}) {
	ack := req0.(*pb.TestMsgAck)
	fmt.Printf("receive  redis server msg:%v, retCode=%v\n", reflect.TypeOf(ack).Elem(), ack.RetCode)
}
