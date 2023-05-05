package internal

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"math/rand"
	"net"
	"redisserver/base/chanrpc"
	"redisserver/base/conf"
	"redisserver/base/network/protobuf"
	"redisserver/base/redisc"
	"redisserver/base/util"
	"redisserver/pb"
	"strconv"
	"sync"
	"time"
)

// redis操作模块
type RedisDB struct {
	chanRpc     *chanrpc.Server // RPC代理服务器
	workers     []*worker       // workers
	DBWorkerCnt int             // worker数量
	redisClient *redisc.Client
	grpcserver  *grpc.Server   //gRpc服务器
	agents      sync.Map       // 会话代理
	Processor   protobuf.Codec // 编码器
	router      *Router        // router
}

var redisDBMgr *RedisDB

// 创建DB模块
func NewRedisDB() *RedisDB {
	return &RedisDB{}
}

// 设置DB模块
func SetRedisDBMgr(redisdb *RedisDB) {
	redisDBMgr = redisdb
}

// 获取DB模块
func GetRedisDBMgr() *RedisDB {
	return redisDBMgr
}
func (db *RedisDB) OnInit() error {
	if db.chanRpc != nil {
		return nil
	}
	db.chanRpc = chanrpc.NewServer(10240)
	db.Processor = pb.Processor
	db.grpcserver = util.NewGRPCServer(5*time.Minute, 5*time.Second)
	//初始化grpc路由
	db.initRouter()
	//连接redis
	client, _ := redisc.NewClient(conf.RedisDBUrl, "123456", 1)
	db.redisClient = client
	db.DBWorkerCnt = conf.DBWorkerNum
	//初始化工作协程
	db.workerInit(client)
	//注册gRPC服务器
	pb.RegisterRedisDBServiceServer(db.grpcserver, db)
	lis, err := net.Listen("tcp", conf.TCPAddr)
	if err != nil {
		logrus.Fatalf("grpc net listen err: %v", err)
	}
	//启动gRPC监听
	go func() {
		err = db.grpcserver.Serve(lis)
		if err != nil {
			logrus.Fatalf("grpc fgate failed: %v", err)
		}
	}()
	return nil
}

// OnTestMsgReq
func (db *RedisDB) OnTestMsgReq(ci *chanrpc.CallInfo, agent *DBAgent) {
	req := ci.Req.(*pb.TestMsgReq)
	len := len(req.Infos)
	if len < 1 {
		return
	}
	db.ChanRPC().Call(
		&RedisDBOp{
			HashKey: int64(req.Infos[0].ID),
			Handler: db.OnTestMsgExc,
			Msg:     req,
			Agent:   agent,
		})
}

// 将接收到的mysql数据存入redis
func (db *RedisDB) OnTestMsgExc(ci *chanrpc.CallInfo, agent *DBAgent) {
	req := ci.Req.(*RedisDBOp).Msg.(*pb.TestMsgReq)
	len := len(req.Infos)
	if len < 1 {
		return
	}
	ack := &pb.TestMsgAck{}
	var ret error
	for _, info := range req.Infos {
		key := strconv.Itoa(int(info.ID))
		ret = db.redisClient.HMSetObject(key, info, 0)
		if ret != nil {
			ack.RetCode = 1
			agent.WriteMsg(ack)
			return
		}
	}
	agent.WriteMsg(ack)
	//key := redisc.NewKey(conf.UserBasePrefix, req.Infos[0].ID)
	//res, _ := db.redisClient.HMGetObject(key, &pb.TestMsgAck{})
	//post := res.(*pb.TestMsgAck)
	//fmt.Printf("get info from redis,retCode=%v", post.RetCode)

}

// Run 启动
func (db *RedisDB) Run(closeSig chan bool) {
	db.workerRun()
	for {
		select {
		case <-closeSig:
			fmt.Println("redisDB  receive close signal")
			return
		case ci := <-db.chanRpc.ChanCall:
			db.Dispatch(ci)
		}
	}
}

func (db *RedisDB) workerInit(client *redisc.Client) {
	for i := 0; i < db.DBWorkerCnt; i++ {
		worker := newWorker(client, conf.WorkChanLen, i+1)
		db.workers = append(db.workers, worker)
	}
}

func (db *RedisDB) workerRun() {
	for i := 0; i < db.DBWorkerCnt; i++ {
		worker := db.workers[i]
		go worker.Run()
	}
}

// Dispatch 根据hashkey分发到worker协程:使用hashkey为了负载均衡
func (db *RedisDB) Dispatch(ci *chanrpc.CallInfo) {
	msg := ci.Req.(*RedisDBOp)
	hashKey := msg.HashKey
	if hashKey == 0 {
		hashKey = rand.Int63()
	}
	workerIdx := hashKey % int64(db.DBWorkerCnt)
	db.workers[workerIdx].chanCall <- ci
}

func (db *RedisDB) OnDestroy() {
	for _, worker := range db.workers {
		worker.OnDestroy()
	}
	db.redisClient.Close()
	db.chanRpc.Close()
	fmt.Println("redisDB closed!")
}

// ChanRPC 消息通道
func (db *RedisDB) ChanRPC() *chanrpc.Server {
	return db.chanRpc
}
