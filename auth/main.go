package main

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"os"
	authpb "server2/auth/api/gen/v1"
	"server2/auth/dao"
	"server2/auth/token"
	"server2/auth/wechat"
	"time"
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
	file, osOpenErr := os.Open("auth/private.key")
	if osOpenErr != nil {
		dLog.Fatal("打开private.key文件失败", zap.Error(err))
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			dLog.Fatal("关闭private.key文件失败", zap.Error(err))
			return
		}
	}(file)
	readAll, readAllErr := ioutil.ReadAll(file)
	if readAllErr != nil {
		dLog.Fatal("读取private.key文件失败", zap.Error(err))
		return
	}
	pemKey, parseJwtErr := jwt.ParseRSAPrivateKeyFromPEM(readAll)
	if parseJwtErr != nil {
		dLog.Fatal("jwt.ParseRSAPrivateKeyFromPEM - 错误", zap.Error(err))
		return
	}
	server := grpc.NewServer()
	authpb.RegisterAuthServiceServer(server, &authpb.Service{
		Logger:         dLog,
		OpenIdResolver: &wechat.Service{},
		Mongo:          dao.NewMongo(mongoConnect.Database("serverDemo")),
		TokenExp:       time.Second * 7000,
		TokenGen:       token.NewJWTTokenGen("qkw123", pemKey),
	})

	err3 := server.Serve(listen)
	if err3 != nil {
		dLog.Fatal("创建服务失败", zap.Error(err3))
		return
	}
}
