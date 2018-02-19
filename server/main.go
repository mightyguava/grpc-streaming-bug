package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/mightyguava/grpc-streaming-bug/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type GreeterServer struct {
	counter int
}

func (g *GreeterServer) SayHello(stream helloworld.Greeter_SayHelloServer) error {
	for {
		_, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Println("server saw client stream closed, closing server stream")
				return stream.SendAndClose(&helloworld.HelloReply{})
			}
			fmt.Println("server got error: ", err)
			return err
		}
		// The client makes requests sequentially, so no locks needed here
		g.counter++
		fmt.Println("server got request: ", g.counter)
	}
}

func main() {
	grpclog.SetLoggerV2(grpclog.NewLoggerV2WithVerbosity(os.Stderr, os.Stderr, os.Stderr, 4))
	server := grpc.NewServer()
	helloworld.RegisterGreeterServer(server, &GreeterServer{})
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	server.Serve(l)
}
