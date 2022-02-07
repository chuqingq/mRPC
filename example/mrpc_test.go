package main

import (
	"errors"
	"testing"

	"github.com/chuqingq/mrpc"
)

// server implementation

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func TestCall(t *testing.T) {
	// server
	rpcs := mrpc.NewRPC()
	defer rpcs.Close()
	err := rpcs.RegisterService("rpcserver-instance-1", new(Arith)) // &Arith{} 不行？
	if err != nil {
		t.Fatalf("rpc.RegisterService error: %v", err)
	}
	// client
	rpc := mrpc.NewRPC()
	defer rpc.Close()
	// Synchronous call
	args := &Args{
		A: 7,
		B: 8,
	}
	var reply int
	err = rpc.Call("rpcserver-instance-1", "Arith.Multiply", args, &reply)
	if err != nil {
		t.Fatalf("rpc call Arith.Multiply error: %v", err)
	}
	if reply != 56 {
		t.Fatalf("arith error: %d*%d!=%d", args.A, args.B, reply)
	}
}
