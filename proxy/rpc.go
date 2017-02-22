package proxy

import (
	// "bytes"
	// "encoding/json"
	// "errors"
	// "etc-pool-admin/utils/common"
	"github.com/cihub/seelog"
	// "math/big"
	"net/http"
	// "strconv"
	"time"
)

type Upstream struct {
	Url     string `json:"url"`
	Timeout string `json:"timeout"`
}

type RPCClient struct {
	Url    string
	client *http.Client
}

//NewRPCClient make a new rpc lient for get data from node
func NewRPCClient(url string, timeout string) *RPCClient {
	rpcClient := &RPCClient{Url: url}
	timeoutIntv, err := time.ParseDuration(timeout)
	if err != nil {
		seelog.Info("parse timeout when new rpc clien error:", err)
		panic(err)
	}
	rpcClient.client = &http.Client{
		Timeout: timeoutIntv,
	}
	return rpcClient
}

// func (r *RPCClient) Subscribe() (ex1, ex2 string) {
// rpcResp, err := r.doPost(r.Url, "mining.subscribe", nil)
// if err != nil {
// 	seelog.Info("subscribe error:", err)
// 	return "", ""
// }
// var reply string
// err = json.Unmarshal(*rpcResp.Result, &reply)
// if err != nil {
// 	seelog.Info("Unmarshal error:", err)
// 	return "", ""
// }
// seelog.Info("result is:", reply)
// return "", ""
// }

func (r *RPCClient) Authorize(user, pwd string) {

}

// //GetAccountBalance get given address balance by rpc
// func (r *RPCClient) GetAccountBalance(account string) (int64, error) {
// 	rpcResp, err := r.doPost(r.Url, "eth_getBalance", []string{account, "latest"})
// 	if err != nil {
// 		return 0, err
// 	}
// 	var reply string
// 	err = json.Unmarshal(*rpcResp.Result, &reply)
// 	if err != nil {
// 		return 0, err
// 	}
//
// 	balance := new(big.Rat).SetInt(common.String2Big(reply))
// 	return weiToShannonInt64(balance), err
// }

// func (r *RPCClient) doPost(url string, method string, params interface{}) (*JSONRpcResp, error) {
// 	jsonReq := map[string]interface{}{"jsonrpc": "2.0", "method": method, "params": params, "id": 0}
// 	data, _ := json.Marshal(jsonReq)
//
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
// 	req.Header.Set("Content-Length", (string)(len(data)))
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Accept", "application/json")
//
// 	resp, err := r.client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
//
// 	var rpcResp *JSONRpcResp
// 	err = json.NewDecoder(resp.Body).Decode(&rpcResp)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if rpcResp.Error != nil {
// 		return nil, errors.New(rpcResp.Error["message"].(string))
// 	}
// 	return rpcResp, err
// }
//
// func weiToShannonInt64(wei *big.Rat) int64 {
// 	shannon := new(big.Rat).SetInt(common.Shannon)
// 	inShannon := new(big.Rat).Quo(wei, shannon)
// 	value, _ := strconv.ParseInt(inShannon.FloatString(0), 10, 64)
// 	return value
// }
