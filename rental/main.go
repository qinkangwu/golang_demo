package main

import (
	"context"
	"github.com/namsral/flag"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	rentalpb "server2/rental/api/gen/v1"
	"server2/rental/api/service"
	"server2/rental/trip/client"
	"server2/rental/trip/dao"
	"server2/shared/auth"
)

var addr = flag.String("addr", ":8089", "rental监听的端口")
var mongoUri = flag.String("mongo_uri", "mongodb://admin:123456@127.0.0.1:27017", "mongodb地址")
var publicKeyPath = flag.String("public_key_file_path", "shared/auth/public.key", "jwt签名publicKey文件地址")

func main() {
	flag.Parse()
	dLog, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalln("创建日志实例失败")
		return
	}
	listen, err2 := net.Listen("tcp", *addr)
	if err2 != nil {
		dLog.Fatal("tcp连接创建失败", zap.Error(err))
		return
	}
	c := context.Background()
	mongoConnect, err := mongo.Connect(c, options.Client().ApplyURI(*mongoUri))
	in, interceptorErr := auth.Interceptor(*publicKeyPath)
	if interceptorErr != nil {
		dLog.Fatal("auth.Interceptor - 错误", zap.Error(interceptorErr))
		return
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(in))
	rentalpb.RegisterTripServiceServer(server, &service.Service{
		Logger:         dLog,
		Mongo:          dao.NewMongo(mongoConnect.Database("serverDemo")),
		ProfileManager: &client.ProfileManager{},
		CarManager:     &client.CarManager{},
		PoiManager:     &client.PoiManager{},
	})

	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	err3 := server.Serve(listen)
	if err3 != nil {
		dLog.Fatal("创建服务失败", zap.Error(err3))
		return
	}
}
