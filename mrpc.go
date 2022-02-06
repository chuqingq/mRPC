package mrpc

// NewRPC 创建RPC对象。服务端和客户端都使用这个对象做RPC。
func NewRPC() *RPC {
	return &RPC{
		ServiceEventChan:    make(chan *ServiceEvent),
		AsyncCallResultChan: make(chan *AsyncCallResult),
	}
}

// RPC RPC结构
type RPC struct {
	ServiceEventChan    chan *ServiceEvent    // 通知service上线或下线
	AsyncCallResultChan chan *AsyncCallResult // 通知异步调用结果
	// other unexported fields
}

// RegisterService 注册服务实例和receiver
// 其中rcvr实现了handler;
// 具体用什么ip、port，随机选取，对用户透明。需要rpc中维护连接池。
func (*RPC) RegisterService(instanceName string, rcvr ...interface{}) error {
	return nil
}

// Call 同步调用。serviceMethod和rcvr名字相同。
func (*RPC) Call(instanceName string, serviceMethod string, args interface{}, reply interface{}) error {
	return nil
}

// AsyncCall 异步调用。
// 通过RPC的AsyncCallResultChan接收响应。
func (*RPC) AsyncCall(instanceName string, serviceMethod string, args interface{}) error {
	return nil
}

// Close 关闭RPC。
func (*RPC) Close() error {
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
