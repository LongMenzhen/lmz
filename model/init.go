package model

import (
	"log"
	"strconv"

	"github.com/cyrnicolase/lmz/config"
	redis "github.com/go-redis/redis"
)

var rds *redis.Client

func init() {
	config := config.Config.Redis
	rds = redis.NewClient(&redis.Options{
		Addr: config.Host + ":" + strconv.Itoa(config.Port),
	})

	if _, err := rds.Ping().Result(); nil != err {
		log.Println("redis连接失败" + err.Error())
	}
}

// Redis 返回redis客户端
func Redis() *redis.Client {
	return rds
}
