package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	rentalpb "server2/rental/api/gen/v1"
	"server2/rental/trip/client"
	"server2/rental/trip/dao"
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
	c := context.Background()
	mongoConnect, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://admin:123456@127.0.0.1:27017"))
	in, interceptorErr := auth.Interceptor("shared/auth/public.key")
	if interceptorErr != nil {
		dLog.Fatal("auth.Interceptor - 错误", zap.Error(interceptorErr))
		return
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(in))
	rentalpb.RegisterTripServiceServer(server, &rentalpb.Service{
		Logger:         dLog,
		Mongo:          dao.NewMongo(mongoConnect.Database("serverDemo")),
		ProfileManager: &client.ProfileManager{},
		CarManager:     &client.CarManager{},
		PoiManager:     &client.PoiManager{},
	})

	err3 := server.Serve(listen)
	if err3 != nil {
		dLog.Fatal("创建服务失败", zap.Error(err3))
		return
	}
}
