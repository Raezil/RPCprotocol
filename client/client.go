package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RPCArgs map[string]any

type RPCResponse struct {
	Result any       `json:"result"`
	Error  *RPCError `json:"error,omitempty"`
}

type RPCError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type RPCClient struct {
	url string
}

func NewRPCClient(url string) *RPCClient {
	return &RPCClient{
		url: url,
	}
}

func (r *RPCClient) Call(args map[string]interface{}) (*RPCResponse, error) {
	// Prepare the request payload

	// Prepare the request payload
	reqBody, err := json.Marshal(args)
	if err != nil {
		fmt.Println("Error encoding request:", err)
		return nil, err
	}

	// Send the request
	resp, err := http.Post(r.url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}

	var rpcResp RPCResponse
	err = json.Unmarshal(body, &rpcResp)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return nil, err
	}

	if rpcResp.Error != nil {
		fmt.Printf("RPC Error: %s\n", rpcResp.Error.Message)
	} else {
		fmt.Printf("RPC Result: %+v\n", rpcResp.Result)
	}
	return &rpcResp, nil
}

func main() {
	url := "http://localhost:8080/rpc"

	rpc := NewRPCClient(url)

	args := map[string]interface{}{
		"method": "add",
		"params": RPCArgs{"a": 1.22,
			"b": 2.33},
	}
	rpc.Call(args)

}
