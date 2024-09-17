package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type RPCServer struct {
	funcs map[string]interface{}
}

type RPCError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (r *RPCError) Error() string {
	return fmt.Sprintf("Error code: %s, %s", r.Code, r.Message)
}

type RPCResponse struct {
	Result any       `json:"result"`
	Error  *RPCError `json:"error,omitempty"`
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

func (rpc *RPCServer) Call(name string, args RPCArgs) *RPCResponse {
	if fn, ok := rpc.funcs[name]; ok {
		resp := fn.(func(RPCArgs) RPCResponse)(args)
		return &resp
	}
	return &RPCResponse{
		Error: &RPCError{
			Code:    "404",
			Message: "Method not found",
		},
	}
}

func (rpc *RPCServer) ServeRPC(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req struct {
		Method string  `json:"method"`
		Params RPCArgs `json:"params"`
	}
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the RPC method
	resp := rpc.Call(req.Method, req.Params)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	rpc := NewRPCServer()
	rpc.Register("test", func(args RPCArgs) RPCResponse {
		return RPCResponse{
			Result: args,
		}
	})

	http.HandleFunc("/rpc", rpc.ServeRPC)

	fmt.Println("RPC server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
