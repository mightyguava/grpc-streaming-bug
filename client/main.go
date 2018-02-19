package main

import (
	"context"
	"time"

	"github.com/mightyguava/grpc-streaming-bug/helloworld"
	"google.golang.org/grpc"
)

func doSend(recv bool, waitBeforeClose time.Duration) {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := helloworld.NewGreeterClient(conn)
	stream, err := client.SayHello(context.Background())
	if err != nil {
		panic(err)
	}
	var name string
	if recv {
		name = "CloseAndRecv()"
		defer stream.CloseAndRecv()
	} else {
		name = "CloseSend()"
		defer stream.CloseSend()
	}
	err = stream.SendMsg(&helloworld.HelloRequest{
		Name: name,
	})
	if err != nil {
		panic(err)
	}
	if waitBeforeClose > 0 {
		time.Sleep(waitBeforeClose)
	}
}

func main() {
	// grpclog.SetLoggerV2(grpclog.NewLoggerV2WithVerbosity(os.Stderr, os.Stderr, os.Stderr, 9))
	for i := 0; i < 10; i++ {
		// This succeeds
		doSend(true, 0)
		// This fails
		doSend(false, 0)
		// This succeeds
		doSend(false, 1*time.Millisecond)
	}
}
