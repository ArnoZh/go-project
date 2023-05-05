package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"redisserver/pb"
	"redisserver/server/internal"
	"syscall"
)

var CloseSig chan bool //关闭信号

func main() {
	//注册网络协议
	pb.Init()
	// 运行redisDB服务器
	go startRedisDBServer()
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
			fmt.Printf("main end")
		}
	}
}

// 启动redisDB服务器
func startRedisDBServer() {
	//创建redisDB
	redisDb := internal.NewRedisDB()
	internal.SetRedisDBMgr(redisDb)
	defer redisDb.OnDestroy()
	//初始化
	redisDb.OnInit()
	// 运行redisDB
	redisDb.Run(CloseSig)
}
