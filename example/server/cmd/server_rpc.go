package main

import (
	"log"

	"mrpc"
	"mrpc/example/server"
)

func main() {
	rpc := mrpc.NewRPC()
	err := rpc.RegisterService("rpcserver-instance-1", new(server.Arith)) // &Arith{} 不行？
	if err != nil {
		log.Fatalf("rpc.RegisterService error: %v", err)
	}
	select {}
}