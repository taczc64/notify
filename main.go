package main

import (
	"fmt"
	"github.com/cihub/seelog"
)

type node struct {
    name:             string  `json:"name"`
    host:             string  `json:"host"`
    port:             int32   `json:"port"`
    worker_name:      string  `json:"worker_name"`
    worker_password:  string  `json:"worker_password"`
    enable:           bool    `json:"enable"`
    debug:            bool    `json:"debug"`
}

type nodes struct {
	Nodes []node
}

func loadNodes() {

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
	loadNodes()

	//start
}
