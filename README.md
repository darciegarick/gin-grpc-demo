#### 使用 Golang 简单的使用 gRPC

1. 使用以下命令安装 Go 的协议编译器插件
```shell
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

2. 更新你的PATH以便protoc编译器可以找到插件
```shell
$ export PATH="$PATH:$(go env GOPATH)/bin"
```

3. 初始化项目
- 首先，我们需要初始化一个新的 Go 项目，并安装必要的依赖。
```shell
mkdir gin-grpc-demo
cd gin-grpc-demo
go mod init gin-grpc-demo
```

- 接下来，安装 Gin 和 gRPC 的 Go 库：
```shell
go get -u github.com/gin-gonic/gin
go get -u google.golang.org/grpc
```

4. gin-grpc-demo 项目下创建 `server` 和 `client` 文件夹,并且俩文件夹下分别创建 `serverproto` 和 `clientproto` 文件夹 (proto生成的go文件放在该文件夹下)
```shell
mkdir server client
cd server
mkdir serverproto
cd ../client
mkdir clientproto
```

5. 创建 `gin-grpc-demo/server/serverproto/demo.proto` 文件
```protobuf
// gin-grpc-demo/server/serverproto/demo.proto


// 使用proto3 语法
syntax = "proto3";

// 这部分内容是关于最后生成的go文件是处在哪个目录哪个包下， . 表示当前目录， proto 表示生成的go文件放在当前目录的 proto 文件夹下
option go_package = "./;serverproto";

// 定义一个服务 服务中有方法，该方法可以接收客户端的参数，再返回服务端的响应
// DemoService 表示服务名 该服务中有一个rpc方法，名为 SayHello
// SayHello方法会发送一个 HelloRequest，然后返回一个 HelloResponse
service DemoService {
  // 定义一个方法
  rpc SayHello (HelloRequest) returns (HelloResponse) {}
}

// 定义一个请求消息
// message 关键字，可以理解为 Golang 中的 struct
// 这里面比较特别的是变量后面的“赋值”。注意，在这里并不是赋值，而是定义这个变量在这个message中的位置
message HelloRequest {
  string name = 1;
  // int64 age = 2;
}

// 定义一个响应消息
message HelloResponse {
  string message = 1;
}
```
6. 编写完上述内容后，在`gin-grpc-demo/server/serverproto`目录下执行如下命令
```shell
protoc --go_out=. demo.proto
protoc --go-grpc_out=. demo.proto
```

7. 创建 `gin-grpc-demo/client/clientproto/demo.proto` 文件

```
// gin-grpc-demo/client/clientproto/demo.proto


syntax = "proto3";
option go_package = "./;clientproto";

// 定义一个服务 
service DemoService {
  // 定义一个方法
  rpc SayHello (HelloRequest) returns (HelloResponse) {}
}

// 定义一个请求消息
message HelloRequest {
  string name = 1;
  // int64 age = 2;
}

// 定义一个响应消息
message HelloResponse {
  string message = 1;
}
```

8. 编写完上述内容后，在`gin-grpc-demo/client/clientproto`目录下执行如下命令
```shell
protoc --go_out=. demo.proto
protoc --go-grpc_out=. demo.proto
```

9. 服务端编写
- 创建 `gRPC Server` 对象
- 将 `server` （其包含需要被调用的接口）注册到 `gRPC Server` 的内部注册中心。
  这样可以在接受到请求时，通过内部的服务发现，发现该服务端接口并转接进行逻辑处理
- 创建`Listen`，监听 TCP 端口
- `gRPC Server` 开始 `lis.Accept`，直到 Stop

**代码实现**
编辑 `gin-grpc-demo/server/main.go` 重写上述 `SayHello` 方法
```go
//gin-grpc-demo/server/main.go
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
```

10. 客户端编写
- 创建与给定目标（服务端）的连接交互
- 创建 `server` 的客户端对象
- 发送 RPC 请求，等待同步响应，得到回调后返回响应结果
- 输出响应结果


**代码实现**
编辑 `gin-grpc-demo/client/main.go`
```go
// gin-grpc-demo/client/main.go

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
```

11. 启动服务端和客户端
```shell
go run gin-grpc-demo/server/main.go
go run gin-grpc-demo/client/main.go
```


















