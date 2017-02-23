package main

import (
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	"notify/proxy"
	"os"
	// "sync"
)

func loadNodes(config *proxy.Pools) {
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
	nodes := proxy.Pools{}
	loadNodes(&nodes)

	//start
	notify := make(chan proxy.BlockInfo, 32)
	for i := 0; i < len(nodes.Nodes); i++ {
		go proxy.MiningClient(nodes.Nodes[i], nodes.Timeout, &notify)
	}

	height := 0
	for {
		select {
		case info := <-notify:
			fmt.Println("main thread get block height from  node:", info.Height)
			if info.Height > height {
				height = info.Height
				seelog.Info("change block height")
				fmt.Println("change block!!!!!!")
			}
		}
	}
}
