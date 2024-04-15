package main

import (
	"context"
	serverproto "gin-grpc-demo/server/serverproto"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	serverproto.UnimplementedDemoServiceServer
}

// 重写 SayHello 方法
func (s *Server) SayHello(ctx context.Context, req *serverproto.HelloRequest) (*serverproto.HelloResponse, error) {
	return &serverproto.HelloResponse{Message: "hello" + req.Name}, nil
}

func main() {
	// 开启端口
	listen, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("did not listen: %v", err)
		return // 或者执行其他适当的操作
	}

	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer()
	// 在 gRPC 服务器上注册服务
	serverproto.RegisterDemoServiceServer(grpcServer, &Server{})

	// 启动服务器
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return
	}
}
