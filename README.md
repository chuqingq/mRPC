# mRPC
mRPC = mDNS + RPC.

## Features

1. mDNS的特性，即局域网服务发现。
2. Go的RPC，即无需像gRPC那样生成代码，可以直接使用。

## Example

RPC protocol and server implementation:

```go
type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

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
```

Server RPC service:

```go
rpc := mrpc.NewRPC()
err := rpc.RegisterService("rpcserver-instance-1", new(server.Arith)) // &Arith{} 不行？
if err != nil {
    log.Fatalf("rpc.RegisterService error: %v", err)
}
select {}
```

Client RPC:

```go
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
```

## Document

TODO
