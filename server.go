package main

import "fmt"

type RPCServer struct {
	funcs map[string]interface{}
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

func (rpc *RPCServer) Call(name string, args RPCArgs) error {
	rpc.funcs[name].(func(RPCArgs) error)(args)
	return nil
}

func main() {
	rpc := NewRPCServer()
	rpc.Register("test", func(args RPCArgs) error {
		fmt.Println(args)
		return nil
	})
	rpc.Call("test", RPCArgs{"key": "value"})
}
