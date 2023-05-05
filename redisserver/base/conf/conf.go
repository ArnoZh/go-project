// Package conf .

package conf

var (
	// ServerID .
	ServerID int32
	// TCPAddr TCP地址
	TCPAddr = "127.0.0.1:2380"
	// ServerName 名字
	ServerName string
	// DBWorkerNum 数据库工作者线程数量
	DBWorkerNum = 3
	// DBUrl 数据库地址(默认连接内网地址服务器)
	RedisDBUrl    = "127.0.0.1:6379"
	MysqlUser     = "root"           // mysql用户名
	MysqlPassword = "123456"         // mysql密码
	MysqlAddr     = "127.0.0.1:3306" // mysql地址
)
