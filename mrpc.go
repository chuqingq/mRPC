package mrpc

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"

	"github.com/hashicorp/mdns"
)

// RPC RPC结构
type RPC struct {
	ServiceEventChan    chan *ServiceEvent    // 通知service上线或下线
	AsyncCallResultChan chan *AsyncCallResult // 通知异步调用结果
	// other unexported fields
	clientsMap sync.Map // string -> serviceInfo
	serversMap sync.Map // string -> server
}

// NewRPC 创建RPC对象。服务端和客户端都使用这个对象做RPC。
func NewRPC() *RPC {
	return &RPC{
		ServiceEventChan:    make(chan *ServiceEvent),
		AsyncCallResultChan: make(chan *AsyncCallResult),
	}
}

// RegisterService 注册服务实例和receiver
// 其中rcvr实现了handler;
// 具体用什么ip、port，随机选取，对用户透明。需要rpc中维护连接池。
func (r *RPC) RegisterService(serviceName string, rcvrs ...interface{}) error {
	port := 8081 // TODO
	// rpc
	for _, rcvr := range rcvrs {
		rpc.Register(rcvr)
	}
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if e != nil {
		log.Fatalf("listen error: %v", e)
	}
	go http.Serve(l, nil) // TODO

	host, _ := os.Hostname()
	info := []string{"My awesome service"}
	service, err := mdns.NewMDNSService(host, serviceName, "", "", port, nil, info)
	if err != nil {
		log.Printf("NewMDNSService error: %v", err)
		return err
	}
	// Create the mDNS server, defer shutdown
	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		log.Printf("NewServer error: %v", err)
		return err
	}
	r.serversMap.Store(serviceName, server)
	return nil
}

func (r *RPC) UnRegisterService(service string) {
	s, ok := r.serversMap.Load(service)
	if ok {
		s.(*mdns.Server).Shutdown()
	}
}

// Call 同步调用。serviceMethod和rcvr名字相同。
func (r *RPC) Call(service string, method string, args interface{}, reply interface{}) error {
	c, err := r.getClient(service)
	if err != nil {
		return err
	}
	log.Printf("enter rpc.Call")
	return c.Call(method, args, reply)
}

type clientInfo struct {
	ip   string // 需要缓存IP端口，后续可能需要重连
	port int
	conn *rpc.Client
	lock sync.Mutex // TODO
}

func (r *RPC) getClient(service string) (*rpc.Client, error) {
	s, ok := r.clientsMap.Load(service)
	if ok {
		return s.(*clientInfo).conn, nil
	}
	// 如果不存在
	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	defer close(entriesCh)

	// Start the lookup
	mdns.Lookup(service, entriesCh)
	for entry := range entriesCh {
		fmt.Printf("Got new entry: %v\n", entry)
		// TODO 维护全局services
		addr := fmt.Sprintf("%s:%v", entry.AddrV4, entry.Port)
		log.Printf("addr: %v", addr)
		c, err := rpc.DialHTTP("tcp", addr)
		if err != nil {
			return nil, err
		}
		log.Printf("rpc.Dial success")
		r.clientsMap.Store(service, &clientInfo{
			ip:   entry.AddrV4.String(),
			port: entry.Port,
			conn: c,
		})
		return c, nil
	}
	return nil, errors.New("service not found")
}

// AsyncCall 异步调用。
// 通过RPC的AsyncCallResultChan接收响应。
func (r *RPC) AsyncCall(instanceName string, serviceMethod string, args interface{}) error {
	return nil
}

// Close 关闭RPC。
func (r *RPC) Close() error {
	close(r.ServiceEventChan)
	close(r.AsyncCallResultChan)
	return nil
}

// 相关数据结构

type ServiceEvent struct {
	ServiceName   string
	ServiceAction int
}

const (
	ServiceActionOnline = iota
	ServiceActionOffline
)

type AsyncCallResult struct {
	// TODO RequestID
	ServiceMethod string      // The name of the service and method to call.
	Args          interface{} // The argument to the function (*struct).
	Reply         interface{} // The reply from the function (*struct).
	Error         error       // After completion, the error status.
}
