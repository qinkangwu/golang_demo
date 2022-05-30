package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	authpb "server2/auth/api/gen/v1"
	"server2/auth/dao"
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
	c := context.Background()
	mongoConnect, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://admin:123456@127.0.0.1:27017"))
	if err != nil {
		dLog.Fatal("mongo连接失败", zap.Error(err))
		return
	}
	server := grpc.NewServer()
	authpb.RegisterAuthServiceServer(server, &authpb.Service{
		Logger:         dLog,
		OpenIdResolver: &wechat.Service{},
		Mongo:          dao.NewMongo(mongoConnect.Database("serverDemo")),
	})

	err3 := server.Serve(listen)
	if err3 != nil {
		dLog.Fatal("创建服务失败", zap.Error(err3))
		return
	}
}
