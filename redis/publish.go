package redis

import (
	"encoding/json"
	"github.com/cihub/seelog"
	"gopkg.in/redis.v3"
)

type RedisConfig struct {
	Endpoint string `json:"endpoint"`
	Password string `json:"password"`
	Database int64  `json:"database"`
	PoolSize int    `json:"poolSize"`
}

type BlockInfo struct {
	Prevhash   string `json:"prevhash"`
	Height     int    `json:"height"`
	Difficulty string `json:"nbits"`
	Ntime      string `json:"ntime"`
	Mintime    string `json:"mintime"`
}

func NewRedisClient(cfg *RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Endpoint,
		Password: cfg.Password,
		DB:       cfg.Database,
		PoolSize: cfg.PoolSize,
	})
	return client
}

func Publish(client *redis.Client, info BlockInfo) {

	data, err := json.Marshal(info)
	if err != nil {
		seelog.Error("json Marshal error:", err)
	}
	client.Publish("pool.blocknotify1st", string(data))
}
