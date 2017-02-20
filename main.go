package main

import (
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	"os"
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

type pools struct {
	Nodes []node `json:"nodes"`
}

func loadNodes(config *pools) {
	filename := "nodes.json"
	seelog.Info("loading nodes file:", filename)
	file, err := os.Open(filename)
	if err != nil {
		seelog.Error("cannot find nodes file:", err.Error())
		panic(err.Error())
	}
	defer file.Close()

	jsonParse := json.NewDecoder(file)
	if err := jsonParse.Decode(&config); err != nil {
		seelog.Critical("Config error: ", err.Error())
		panic(err.Error())
	}
}

func listenNode(conf node, notify chan int) {
	notify <- 1
}

func main() {
	//for recover
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("critical error, recover:", r)
		}
	}()

	//load log config
	logger, err := seelog.LoggerFromConfigAsFile("./logconfig.xml")
	if err != nil {
		panic(err)
	}
	seelog.ReplaceLogger(logger)
	defer seelog.Flush()

	//load all node message
	nodes := pools{}
	loadNodes(&nodes)

	//start
	notify := make(chan int, 1)
	for i := 0; i < len(nodes.Nodes); i++ {
		go listenNode(nodes.Nodes[i], notify)
	}

	for {
		select {
		case data := <-notify:
			fmt.Println("data:", data)
		}
	}
}
