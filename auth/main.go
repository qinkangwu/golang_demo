package main

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	authpb "server2/auth/api/gen/v1"
	"server2/auth/wechat"
)

func main() {
	dLog, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalln("创建日志实例失败")
		return
	}
	listen, err2 := net.Listen("tcp", ":8088")
	if err2 != nil {
		dLog.Fatal("tcp连接创建失败", zap.Error(err))
		return
	}
	server := grpc.NewServer()
	authpb.RegisterAuthServiceServer(server, &authpb.Service{
		Logger:         dLog,
		OpenIdResolver: &wechat.Service{},
	})

	err3 := server.Serve(listen)
	if err3 != nil {
		dLog.Fatal("创建服务失败", zap.Error(err3))
		return
	}
}
