// Package util .

package util

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// NewGRPCServer 创建grpc服务
// 当 Unary RPC 或 Stream RPC 服务panic时，可以避免节点挂掉，参考：
// https://godoc.org/github.com/grpc-ecosystem/go-grpc-middleware/recovery
// 注意服务实现方法中，不能写类似 defer util.PrintPanicStack() 等panic恢复逻辑
func NewGRPCServer(keepAliveTime, keepAliveTimeout time.Duration) *grpc.Server {
	//opts := []grpc_recovery.Option{
	//	grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
	//		return fmt.Errorf("panic: %v", p)
	//	}),
	//}
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    keepAliveTime,
			Timeout: keepAliveTimeout,
		}))
	//grpc_middleware.WithUnaryServerChain(
	//		grpc_recovery.UnaryServerInterceptor(opts...),
	//	),
	//	grpc_middleware.WithStreamServerChain(
	//		grpc_recovery.StreamServerInterceptor(opts...),
	//	))
	return grpcServer
}
