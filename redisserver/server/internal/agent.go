package internal

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"redisserver/base/conf"
	"redisserver/base/network/msgpool"
	"redisserver/base/network/protobuf"
	"redisserver/base/util"
	"redisserver/base/util/idgen"
	"redisserver/pb"
	"reflect"
)

type DBAgent struct {
	sessionID int64                          // agent session id
	stream    pb.RedisDBService_StreamServer // gRPC连接流
	recvMQ    chan *pb.StreamData            // 正向消息处理
	sendMQ    chan interface{}               // 发送消息缓冲
	closeSig  chan bool                      // 是否关闭完毕
	codec     protobuf.Codec                 // 编码器
	router    *Router                        // router

}

// NewDBAgent 创建代理 dispatcher必须线程安全
func NewDBAgent(stream pb.RedisDBService_StreamServer, processor protobuf.Codec, router *Router) *DBAgent {
	agent := &DBAgent{
		sessionID: idgen.NewID(),
		stream:    stream,                                             // gRPC流
		sendMQ:    make(chan interface{}, conf.SingleClientSendMqLen), // 发送消息缓冲
		closeSig:  make(chan bool, 1),                                 // 关闭就绪消息
		codec:     processor,                                          // 编码器
		router:    router,
	}
	return agent
}

// WriteMsg 向agent发送消息
func (agent *DBAgent) WriteMsg(msg interface{}) {
	if msg == nil {
		return
	}
	logrus.Debugf("DBAgent.WriteMsg msg:%v", msg)
	agent.streamSend(msg)
}

// sendToFrame 向客户端发送消息
func (agent *DBAgent) sendToFrame(msg interface{}) {
	defer util.PrintPanicStack()
	defer logrus.Debugf("session%v agent.stream.Send msg: %v", agent.sessionID, msg.(protobuf.IMsg))
	//if err := agent.reqStat.GotAck(msg); err != nil {
	//	logrus.Warnf("%v goAck failed:%v", agent, err)
	//}

	var buffer *msgpool.Buffer
	var data []byte
	if pbmsg, ok := msg.(protobuf.IMsg); ok {
		b, err := agent.codec.MarshalToBuffer(pbmsg, 1)
		if err != nil {
			logrus.Errorf("%v marshal %v err: %v", agent, reflect.TypeOf(msg).Elem().Name(), err)
			if b != nil {
				b.Release()
			}
			return
		}
		buffer = b
		data = buffer.Bytes
		//agent.frameStat.Add(reflect.TypeOf(msg).Elem().Name(), int64(len(data)))
	} else if b, ok := msg.(*msgpool.Buffer); ok {
		buffer = b
		data = buffer.Bytes
	} else if b, ok := msg.([]byte); ok {
		data = b
	} else {
		logrus.Errorf("%v send invalid msg type: %v", agent, reflect.TypeOf(msg))
		return
	}

	// 发送到client
	err := agent.stream.Send(&pb.StreamData{
		Msg:   data,
		GenTs: util.NowTs(),
	})
	fmt.Printf("send msg to stream client,msg:%v\n", reflect.TypeOf(msg).Elem())
	if buffer != nil {
		buffer.DeRef()
	}
	if err != nil {
		logrus.Errorf("%v agent.stream.Send err: %v", agent, err)
	}
}

// Run 启动代理
func (agent *DBAgent) Run() {
	defer util.PrintPanicStack()

	agent.recvMQ = goStreamReader(agent.stream)

	for {
		select {
		case frame := <-agent.recvMQ: //接收消息
			if frame == nil {
				fmt.Println("close by agent")
				return
			}
			// 处理帧
			if ok := agent.recvFrame(frame); !ok {
				fmt.Println("%v agent.recvFrame", agent.sessionID)
			}
		case msg := <-agent.sendMQ: //发送消息
			// 空消息强制踢掉
			if msg == nil {
				logrus.Debugf("%v kickout", agent.sessionID)
				err := agent.stream.Send(&pb.StreamData{})
				if err != nil {
					logrus.Debugf("%v agent.stream.Send err: %v", agent.sessionID, err)
				}
				continue
			}
			agent.sendToFrame(msg)
		}
	}
}

// streamSend 将消息放入写缓冲
func (agent *DBAgent) streamSend(msg interface{}) {
	select {
	case agent.sendMQ <- msg:
	default:
		if b, ok := msg.(*msgpool.Buffer); ok {
			b.DeRef() //失败的 引用计数要减
		}
	}
}

// recvFrame 处理来自其他服的帧消息
func (agent *DBAgent) recvFrame(frame *pb.StreamData) bool {
	msgID, err := agent.codec.ReadMsgID(frame.Msg)
	if err != nil {
		logrus.Errorf("recvFrame read msg id error: %v", err)
		return false
	}

	// hook
	if hook := agent.router.getHook(msgID); hook != nil {
		_, msg, err := agent.codec.Unmarshal(frame.Msg)
		if err != nil {
			logrus.Errorf("recvFrame unmarshal msg error: %v", err)
		}
		fmt.Printf("recive client stream ,msg:%v\n", reflect.TypeOf(msg).Elem())
		hook(agent, msg.(protobuf.IMsg))
		return true
	}

	//if err = agent.reqStat.PendReq(frame.Msg); err != nil {
	//	logrus.Errorf("recvFrame pendReq failed: %v", err)
	//}

	return false
}

// goStreamReader 从gRPC连接中读取消息帧放入缓冲池 并返回该缓冲池
func goStreamReader(stream pb.RedisDBService_StreamServer) chan *pb.StreamData {
	recvMQ := make(chan *pb.StreamData, conf.SingleClientRecvStreamLen)
	go func() {
		for {
			frame, err := stream.Recv()
			if err != nil {
				recvMQ <- nil
				break
			}
			recvMQ <- frame
		}
	}()
	return recvMQ
}

// Close 逻辑主动关闭网络连接
func (agent *DBAgent) Close() {
	agent.streamSend(nil)
}

// OnClose 关闭时处理
func (agent *DBAgent) OnClose() {
	close(agent.closeSig)
}
