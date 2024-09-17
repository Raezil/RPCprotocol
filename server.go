package main

import "fmt"

type RPCServer struct {
	funcs map[string]interface{}
}

type RPCError struct {
	Code    string
	Message string
}

func (r *RPCError) Error() {
	fmt.Errorf("Error code: %s, %s", r.Code, r.Message)
}

type RPCResponse struct {
	Result any
	Error  RPCError
}

func NewRPCServer() *RPCServer {
	return &RPCServer{
		funcs: make(map[string]interface{}),
	}
}

type RPCArgs map[string]any

func (rpc *RPCServer) Register(name string, f interface{}) {
	if rpc.funcs == nil {
		rpc.funcs = make(map[string]interface{})
	}
	rpc.funcs[name] = f

}

func (rpc *RPCServer) Call(name string, args RPCArgs) any {
	rpc.funcs[name].(func(RPCArgs) any)(args)
	return nil
}

func main() {
	rpc := NewRPCServer()
	rpc.Register("test", func(args RPCArgs) any {
		fmt.Println(args)
		return nil
	})
	rpc.Call("test", RPCArgs{"key": "value"})
}
