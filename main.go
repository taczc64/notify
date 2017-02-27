package main

import (
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	"notify/proxy"
	"notify/redis"
	"os"
	"runtime"
)

func loadConfig(config *proxy.Config) {
	filename := "config.json"
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
	config := proxy.Config{}
	loadConfig(&config)

	//use multi core if has, I set default 2, if have only 1 core, it will automatic change
	runtime.GOMAXPROCS(2)
	seelog.Info("running witch %d threads...", 2)

	redisClient := redis.NewRedisClient(&config.Redis)
	pong, err := redisClient.Ping().Result()
	if err != nil {
		seelog.Infof("Can't establish connection to redis: %v", err)
	} else {
		seelog.Infof("redis check reply: %v", pong)
	}

	//start
	notify := make(chan redis.BlockInfo, 16)

	for i := 0; i < len(config.Nodes); i++ {
		go proxy.MiningClient(config.Nodes[i], config.Timeout, &notify)
	}

	height := 0
	difficulty := ""

	for {
		select {
		case info := <-notify:
			if height == 0 {
				height = info.Height - 1
				difficulty = info.Difficulty
				continue
			}
			if info.Height > height {
				if difficulty != info.Difficulty && (info.Height%2016) != 0 {
					info.Difficulty = difficulty
				}
				//publish to redis
				redis.Publish(redisClient, info)
				height = info.Height
			}
		}
	}
}
