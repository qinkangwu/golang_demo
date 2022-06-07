package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/namsral/flag"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net/http"
	authpb "server2/auth/api/gen/v1"
	rentalpb "server2/rental/api/gen/v1"
)

var addr = flag.String("addr", ":8090", "gateway监听的端口")
var authAddr = flag.String("auth_addr", ":8088", "auth监听的端口")
var rentalAddr = flag.String("rental_addr", ":8089", "rental监听的端口")

func main() {
	flag.Parse()
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
		*authAddr,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err2 != nil {
		dLog.Fatal("网关创建错误", zap.Error(err2))
		return
	}
	err2 = rentalpb.RegisterTripServiceHandlerFromEndpoint(
		c,
		mux,
		*rentalAddr,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err2 != nil {
		dLog.Fatal("网关创建错误", zap.Error(err2))
		return
	}
	err3 := http.ListenAndServe(*addr, mux)
	if err3 != nil {
		dLog.Fatal("网关监听错误", zap.Error(err3))
		return
	}
}
