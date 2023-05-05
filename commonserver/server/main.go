package main

import (
	"commonserver/pb"
	"commonserver/server/internal"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var CloseSig chan bool //关闭信号

func main() {
	//注册网络协议
	pb.Init()
	// 运行mysql服务器
	go startMysqlDBServer()
	// 退出信号监听
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	for {
		select {
		case sig := <-c:
			logrus.Infof("redisDB server closing down (signal: %v)", sig)
			if sig == syscall.SIGHUP {
				continue
			}
			CloseSig <- true
			internal.GetMysqlDBMgr().OnDestroy()
			fmt.Printf("main end")
		}
	}
}

// 启动mysqlDB服务器
func startMysqlDBServer() {
	//创建DB
	mysqlDB := internal.NewMysqlDB()
	//初始化
	mysqlDB.OnInit()
	internal.SetMysqlDBMgr(mysqlDB)
	// 运行ysqlDB
	mysqlDB.Run(CloseSig)
	//给redisdb发送请求
	mysqlDB.OnTestMsgReq()
}
