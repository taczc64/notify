package proxy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	"net"
	"notify/redis"
	"strconv"
	"time"
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

type Config struct {
	Timeout string            `json:"timeout"`
	Redis   redis.RedisConfig `json:"redis"`
	Nodes   []node            `json:"nodes"`
}

type JSONRpcReq struct {
	Id     *json.RawMessage `json:"id"`
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
}

type Req struct {
	Id     int      `json:"id"`
	Method string   `json:"method"`
	Params []string `json:"params"`
}

type Resp struct {
	Id     int         `json:"id"`
	Result interface{} `json:"result"`
	Error  interface{} `json:"error,omitempty"`
}

type Resdata struct {
	Id     interface{}   `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

func MiningClient(conf node, timeout string, notify *chan redis.BlockInfo) {
	//recover
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println("routine critical error:", rec)
		}
	}()

	if conf.Enable {
		value, err := time.ParseDuration(timeout)
		if err != nil {
			panic(err)
		}
		for {
			seelog.Info("begin to dial node:", conf.Name)
			addr := conf.Host + ":" + strconv.Itoa(int(conf.Port))
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				seelog.Error("cannot connect to pool:", conf.Name)
				time.Sleep(value)
				continue
			}
			err = sendSubscribe(&conn)
			if err != nil {
				seelog.Errorf("error from node %s : %v", conf.Name, err)
			}
			sendAuthorize(&conn, conf.WorkerName, conf.WorkerPassword, notify, conf.Name)
			conn.Close()
			time.Sleep(value)
			seelog.Infof("redial %s will be happen in %s.", conf.Name, timeout)
		}
	}
}

func sendSubscribe(conn *net.Conn) error {
	enc := json.NewEncoder(*conn)

	id := []byte(string(strconv.Itoa(1)))
	params, _ := json.Marshal(map[string]interface{}{"params": ""})

	req := JSONRpcReq{Id: (*json.RawMessage)(&id), Method: "mining.subscribe", Params: (*json.RawMessage)(&params)}
	err := enc.Encode(&req)
	if err != nil {
		return err
	}
	connbuf := bufio.NewReaderSize(*conn, 128)
	data, _, err := connbuf.ReadLine()
	if err != nil {
		return err
	}
	resp := Resp{}
	json.Unmarshal(data, &resp)
	if resp.Error != nil {
		seelog.Error("send subscribe error, resp:", string(data))
		return err
	}
	return err
}

func sendAuthorize(conn *net.Conn, worker, pwd string, notify *chan redis.BlockInfo, nodename string) {
	enc := json.NewEncoder(*conn)

	// id, _ := json.Marshal(map[string]interface{}{"id": 2})

	req := Req{Id: 2, Method: "mining.authorize", Params: []string{worker, pwd}}
	err := enc.Encode(&req)
	if err != nil {
		seelog.Errorf("send reqeust to %s error: %v", nodename, err)
		return
	}
	connbuf := bufio.NewReaderSize(*conn, 2048)

	data, _, err := connbuf.ReadLine()
	if err != nil {
		seelog.Errorf("get response from %s error: %v", nodename, err)
	}

	resp := Resp{}
	json.Unmarshal(data, &resp)
	if resp.Error != nil {
		seelog.Infof("authorize from node %s error: %v", nodename, resp.Error)
	}

	lastntime := ""
	resdata := Resdata{}
	for {
		data, _, err = connbuf.ReadLine()
		if err != nil {
			seelog.Errorf("get response from %s error: %v", nodename, err)
			break
		}
		err = json.Unmarshal(data, &resdata)
		if err != nil {
			seelog.Error("Unmarshal authorize response error:", err)
			break
		}
		if resdata.Method == "mining.notify" {
			handleNotify(resdata.Params, &lastntime, notify)
		} else if resdata.Method == "mining.set_difficulty" {
			// TODO
		}
	}
}

func handleNotify(params []interface{}, lastntime *string, notify *chan redis.BlockInfo) {
	var difficulty, prevhash, ntime string
	var blockheight int64
	if clean, ok := params[8].(bool); ok && clean == true {
		if value, ok := params[2].(string); ok {
			// get block height
			hash := []byte(value)
			h := "0x" + string(hash[84:86]) //十六进制
			blockHeightWei, _ := strconv.ParseInt(h, 0, 4)
			height := hash[86:(86 + int(blockHeightWei)*2)]
			newslice := convert(height)
			newh := "0x" + string(newslice)
			blockheight, _ = strconv.ParseInt(newh, 0, 32)

		}
		//get difficulty , prev hash and Ntime
		difficulty = params[6].(string)
		prevhash = params[1].(string)
		ntime = params[7].(string)
		*notify <- redis.BlockInfo{Prevhash: prevhash, Height: int(blockheight), Difficulty: difficulty, Ntime: ntime, Mintime: *lastntime}
		//set last ntime
		*lastntime = ntime
	}
}

func convert(data []byte) []byte {
	num := len(data) / 2 //slice对数目
	newSlice := make([]byte, 0, len(data))
	for i := 1; i < num+1; i++ {
		newSlice = append(newSlice, data[len(data)-i*2:len(data)-i*2+2]...)
	}
	return newSlice
}
