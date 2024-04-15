package main

import (
	"context"
	"fmt"
	"log"

	serverproto "gin-grpc-demo/server/serverproto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 连接 server 端，此处禁用安全协议，没有进行加密验证
	conn, err := grpc.Dial("127.0.0.1:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// 关闭连接
	defer conn.Close()

	// 创建客户端 建立连接
	client := serverproto.NewDemoServiceClient(conn)

	// 执行 rpc 的调用（这个方法在服务器端来实现并返回结果）
	resp, _ := client.SayHello(context.Background(), &serverproto.HelloRequest{Name: "ziying"})

	// 打印结果
	fmt.Println(resp.GetMessage())
}
