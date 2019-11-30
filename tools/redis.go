package tools

import (
	"time"

	"github.com/go-redis/redis"
)

// RedisClientMap ...
var RedisClientMap = map[string]*redis.Client{}

// RedisOption ...
type RedisOption struct {
	Address  string
	Password string
	Db       int
}

// GetRedisInstance ...
func GetRedisInstance(redisOpt RedisOption) *redis.Client {
	address := redisOpt.Address
	db := redisOpt.Db
	password := redisOpt.Password
	if redisCli, ok := RedisClientMap[address]; ok {
		return redisCli
	}
	client := redis.NewClient(&redis.Options{
		Addr:       address,
		Password:   password,
		DB:         db,
		MaxConnAge: 20 * time.Second,
	})
	RedisClientMap[address] = client
	return RedisClientMap[address]
}
