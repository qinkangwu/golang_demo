package main

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	rentalpb "server2/rental/api/gen/v1"
	"server2/shared/auth"
)

func main() {
	dLog, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalln("创建日志实例失败")
		return
	}
	listen, err2 := net.Listen("tcp", ":8089")
	if err2 != nil {
		dLog.Fatal("tcp连接创建失败", zap.Error(err))
		return
	}
	in, interceptorErr := auth.Interceptor("shared/auth/public.key")
	if interceptorErr != nil {
		dLog.Fatal("auth.Interceptor - 错误", zap.Error(interceptorErr))
		return
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(in))
	rentalpb.RegisterTripServiceServer(server, &rentalpb.Service{
		Logger: dLog,
	})

	err3 := server.Serve(listen)
	if err3 != nil {
		dLog.Fatal("创建服务失败", zap.Error(err3))
		return
	}
}
