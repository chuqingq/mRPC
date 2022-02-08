package mrpc

import (
	"errors"
	"testing"
)

// rpc protocol

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

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
	service := "_foobar._tcp"
	// server
	rpcs := NewRPC()
	defer rpcs.Close()
	err := rpcs.RegisterService(service, new(Arith)) // &Arith{} 不行？
	if err != nil {
		t.Fatalf("rpc.RegisterService error: %v", err)
	}
	// client
	{
		rpc := NewRPC()
		// Synchronous call
		args := &Args{
			A: 7,
			B: 8,
		}
		var reply int
		err = rpc.Call(service, "Arith.Multiply", args, &reply)
		if err != nil {
			t.Fatalf("rpc call Arith.Multiply error: %v", err)
		}
		if reply != 56 {
			t.Fatalf("arith error: %d*%d!=%d", args.A, args.B, reply)
		}
		rpc.Close()
	}
	// client again
	{
		rpc := NewRPC()
		// Synchronous call
		args := &Args{
			A: 6,
			B: 9,
		}
		var reply int
		err = rpc.Call(service, "Arith.Multiply", args, &reply)
		if err != nil {
			t.Fatalf("rpc call Arith.Multiply error: %v", err)
		}
		if reply != 54 {
			t.Fatalf("arith error: %d*%d!=%d", args.A, args.B, reply)
		}
		rpc.Close()
	}
}
