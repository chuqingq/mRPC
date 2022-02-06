package main

import (
	"log"

	"mrpc"
	"mrpc/example/server"
)

func main() {
	rpc := mrpc.NewRPC()

	// Synchronous call
	args := &server.Args{
		A: 7,
		B: 8,
	}
	var reply int
	err := rpc.Call("rpcserver-instance-1", "Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	log.Printf("Arith: %d*%d=%d", args.A, args.B, reply)
}
