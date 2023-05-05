package internal

import (
	"commonserver/base/chanrpc"
	"commonserver/base/conf"
	"commonserver/base/network/msgpool"
	"commonserver/base/network/protobuf"
	"commonserver/base/util"
	"commonserver/pb"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"reflect"
	"sync"
)

type MysqlDB struct {
	chanRpc     *chanrpc.Server                // RPC代理服务器
	SqlLink     *sql.DB                        //mysql 连接
	Processor   protobuf.Codec                 // 编码器
	router      *Router                        // router
	clientSteam pb.RedisDBService_StreamClient // 连接redis的grpc流
	sendMQ      chan interface{}               // 发送消息缓冲
	streamConn  *grpc.ClientConn               //与redisserver的grpc流连接
}

var mysqlDBMgr *MysqlDB

var mu sync.RWMutex

type UserInfo struct {
	name string
	age  int16
	home string
}

var UserInfos = make(map[int]*UserInfo)

// 连接mysql
func (db *MysqlDB) ConnectMysql() {
	url := conf.MysqlUser + ":" + conf.MysqlPassword + "@tcp(" + conf.MysqlAddr + ")/userdb?charset=utf8"
	dbLink, err := sql.Open("mysql", url)
	if err != nil {
		panic(err)
	}
	db.SqlLink = dbLink
	fmt.Println("mysql connect userdb sucess\n")
}

// 从mysql加载user表数据
func (db *MysqlDB) LoadUserTable() {
	rows, err := db.SqlLink.Query("select * from user")
	if err != nil {
		panic(err)
	}
	mu.Lock()
	defer mu.Unlock()

	UserInfos = make(map[int]*UserInfo)

	for rows.Next() {
		var id int
		info := &UserInfo{}
		err = rows.Scan(&id, &info.name, &info.age, &info.home)
		UserInfos[id] = info
	}
	logrus.Infof("load user table success, len(info):%d", len(UserInfos))

}

// 获取用户信息表数据
func GetUserTableData() map[int]*UserInfo {
	return UserInfos
}

// 创建DB模块
func NewMysqlDB() *MysqlDB {
	return &MysqlDB{
		sendMQ: make(chan interface{}, conf.SingleClientSendMqLen),
	}
}

// 设置DB模块
func SetMysqlDBMgr(dbmgr *MysqlDB) {
	mysqlDBMgr = dbmgr
}

// 获取DB模块
func GetMysqlDBMgr() *MysqlDB {
	return mysqlDBMgr
}
func (db *MysqlDB) OnInit() {
	if db.chanRpc != nil {
		return
	}
	db.chanRpc = chanrpc.NewServer(10240)
	db.Processor = pb.Processor
	//db.grpcserver = util.NewGRPCServer(5*time.Minute, 5*time.Second)
	// 连接mysql
	db.ConnectMysql()
	//加载user表
	db.LoadUserTable()
	//初始化grpc路由
	db.initRouter()
	//创建连接
	conn, err := grpc.Dial(conf.RedisServer, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("connect redisserver falid!,err=%v\n", err)
		return
	}
	db.streamConn = conn
	//声明客户端
	client := pb.NewRedisDBServiceClient(conn)
	//创建数据流
	stream, err := client.Stream(context.Background())
	if err != nil {
		fmt.Printf("create Stream falid!,err=%v\n", err)
		return
	}
	db.clientSteam = stream
	return
}

// OnTestMsgReq
func (db *MysqlDB) OnTestMsgReq() {
	useInfos := GetUserTableData()
	if len(useInfos) < 1 {
		return
	}
	msg := &pb.TestMsgReq{}
	for id, info := range useInfos {
		msg.Infos = append(msg.Infos, &pb.UserInfo{
			ID:   int32(id),
			Name: info.name,
			Age:  int32(info.age),
			Home: info.home,
		})
	}
	db.streamSend(msg)
}

// streamSend 将消息放入写缓冲
func (db *MysqlDB) streamSend(msg interface{}) {
	select {
	case db.sendMQ <- msg:
		fmt.Println("streamSend put msg to sendMQ")
	default:
		fmt.Println("streamSend recive error msg")
	}
}

func (db *MysqlDB) sendMsgToRedisServer(msg interface{}) {
	if db.clientSteam == nil {
		return
	}
	var buffer *msgpool.Buffer
	var data []byte
	if pbmsg, ok := msg.(protobuf.IMsg); ok {
		b, err := db.Processor.MarshalToBuffer(pbmsg, 1)
		if err != nil {
			logrus.Errorf(" sendMsgToRedisServer marshal %v err: %v", reflect.TypeOf(msg).Elem().Name(), err)
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
		logrus.Errorf("sendMsgToRedisServer send invalid msg type: %v", reflect.TypeOf(msg))
		return
	}
	err := db.clientSteam.Send(&pb.StreamData{Msg: data, GenTs: util.NowTs()})
	if err != nil {
		fmt.Errorf("send to redis server, [%v] stream.send err:%v", reflect.TypeOf(msg).Elem(), err)
		return
	} else {
		fmt.Printf("send msg to redis server,msg:%v\n", reflect.TypeOf(msg).Elem())
	}
}

// Run 启动
func (db *MysqlDB) Run(closeSig chan bool) {
	//发送消息
	go func() {
		for {
			select {
			case sendMsg := <-db.sendMQ:
				if sendMsg == nil {
					continue
				}
				db.sendMsgToRedisServer(sendMsg)
			}
		}
	}()
	//接收消息
	go func() {
		var recvMQ chan *pb.StreamData
		recvMQ = goClientStreamReader(db.clientSteam)
		for {
			select {
			case frame := <-recvMQ:
				if frame == nil {
					fmt.Printf("Stream closed")
					continue
				}
				//处理帧
				id, msg, err := db.Processor.Unmarshal(frame.Msg)
				if err != nil {
					fmt.Errorf("clientSteam Unmarshal err:%v", err)
					continue
				}
				// 路由处理消息
				handler := db.router.getHook(id)
				if handler != nil {
					handler(nil, msg.(protobuf.IMsg))
					continue
				}
			}
		}
	}()

}

func (db *MysqlDB) OnDestroy() {
	db.SqlLink.Close()
	db.chanRpc.Close()
	db.streamConn.Close()
	fmt.Println("mysql DB closed!")
}

// ChanRPC 消息通道
func (db *MysqlDB) ChanRPC() *chanrpc.Server {
	return db.chanRpc
}
