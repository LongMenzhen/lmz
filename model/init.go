package model

import (
	"strconv"

	"github.com/cyrnicolase/lmz/config"
	rds "github.com/go-redis/redis"
)

var redis *rds.Client

func init() {
	config := config.Config.Redis
	redis = rds.NewClient(&rds.Options{
		Addr: config.Host + ":" + strconv.Itoa(config.Port),
	})
}

// Redis 返回redis客户端
func Redis() *rds.Client {
	return redis
}
