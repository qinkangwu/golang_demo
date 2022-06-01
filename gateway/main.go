package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net/http"
	authpb "server2/auth/api/gen/v1"
	rentalpb "server2/rental/api/gen/v1"
)

func main() {
	dLog, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalln("创建日志实例失败")
		return
	}
	c := context.Background()
	c, cancel := context.WithCancel(c)
	defer cancel()
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard,
		&runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseEnumNumbers: true,
				UseProtoNames:  false,
			},
		},
	))
	err2 := authpb.RegisterAuthServiceHandlerFromEndpoint(
		c,
		mux,
		":8088",
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err2 != nil {
		dLog.Fatal("网关创建错误", zap.Error(err2))
		return
	}
	err2 = rentalpb.RegisterTripServiceHandlerFromEndpoint(
		c,
		mux,
		":8089",
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err2 != nil {
		dLog.Fatal("网关创建错误", zap.Error(err2))
		return
	}
	err3 := http.ListenAndServe(":8090", mux)
	if err3 != nil {
		dLog.Fatal("网关监听错误", zap.Error(err3))
		return
	}
}
