package proxy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	"net"
	"strconv"
)

type node struct {
	Name           string `json:"name"`
	Host           string `json:"host"`
	Port           int32  `json:"port"`
	WorkerName     string `json:"worker_name"`
	WorkerPassword string `json:"worker_password"`
	Enable         bool   `json:"enable"`
	Debug          bool   `json:"debug"`
}

type Pools struct {
	Timeout string `json:"timeout"`
	Nodes   []node `json:"nodes"`
}

type JSONRpcReq struct {
	Id     *json.RawMessage `json:"id"`
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
}

type JSONRpcResp struct {
	Id      *json.RawMessage `json:"id"`
	Version string           `json:"jsonrpc"`
	Result  interface{}      `json:"result"`
	Error   interface{}      `json:"error,omitempty"`
}

func MiningClient(conf node, timeout string, notify chan int) {
	if conf.Enable {
		seelog.Info("begin to dial node:", conf.Name)
		addr := conf.Host + ":" + strconv.Itoa(int(conf.Port))
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			seelog.Error("cannot connect to pool:", conf.Name)
			//TODO
			//reconnect
		}
		defer conn.Close()
		sendSubscribe(&conn)
		authorize(&conn, conf.WorkerName, conf.WorkerPassword)
	}
}

func sendSubscribe(conn *net.Conn) {
	enc := json.NewEncoder(*conn)

	id := []byte(string(strconv.Itoa(1)))
	// params := []byte(`{"Params":""}`)
	params := []byte(`{"method": "subscribe"}`)

	req := JSONRpcReq{Id: (*json.RawMessage)(&id), Method: "mining.subscribe", Params: (*json.RawMessage)(&params)}
	err := enc.Encode(&req)
	if err != nil {
		seelog.Error("send reqeust error:", err)
	}
	fmt.Println("send over")
	connbuf := bufio.NewReaderSize(*conn, 128)
	data, _, err := connbuf.ReadLine()
	if err != nil {
		seelog.Error("get response error:", err)
	}
	seelog.Info("data :", string(data), "data len:", len(data))
}

func authorize(conn *net.Conn, worker, pwd string) {
	enc := json.NewEncoder(*conn)

	id := []byte(string(strconv.Itoa(1)))
	// params := []byte(`{"Params":""}`)
	params := []byte(`{"method": "subscribe"}`)

	req := JSONRpcReq{Id: (*json.RawMessage)(&id), Method: "mining.authorize", Params: (*json.RawMessage)(&params)}
	err := enc.Encode(&req)
	if err != nil {
		seelog.Error("send reqeust error:", err)
	}
	fmt.Println("send over")
	connbuf := bufio.NewReaderSize(*conn, 128)
	data, _, err := connbuf.ReadLine()
	if err != nil {
		seelog.Error("get response error:", err)
	}
	seelog.Info("data :", string(data), "data len:", len(data))
}
