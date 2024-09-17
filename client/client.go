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

func main() {
	url := "http://localhost:8080/rpc"

	// Prepare the request payload
	reqBody, err := json.Marshal(map[string]interface{}{
		"method": "add",
		"params": RPCArgs{"a": 1.22,
			"b": 2.33},
	})
	if err != nil {
		fmt.Println("Error encoding request:", err)
		return
	}

	// Send the request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var rpcResp RPCResponse
	err = json.Unmarshal(body, &rpcResp)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	if rpcResp.Error != nil {
		fmt.Printf("RPC Error: %s\n", rpcResp.Error.Message)
	} else {
		fmt.Printf("RPC Result: %+v\n", rpcResp.Result)
	}
}
