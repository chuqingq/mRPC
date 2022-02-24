package mrpc

import (
	"errors"
	"log"
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
	log.Printf("Multiply: %v", *t)
	*reply = args.A * args.B
	*t = Arith(*reply)
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
	service := "server-instance-1"
	// server
	rpcs := NewRPC()
	defer rpcs.Close()
	a := new(Arith)
	err := rpcs.RegisterService(service, a) // &Arith{} 不行？
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
		if int(*a) != reply {
			t.Fatalf("arith error2: %d*%d!=%d", args.A, args.B, reply)
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
