package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var RDB *redis.Client

func Init() (err error) {
	RDB = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			viper.GetString("redis.host"),
			viper.GetInt("redis.port"),
		),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
		PoolSize: viper.GetInt("redis.pool_size"),
	})
	pong, err := RDB.Ping().Result()
	if err != nil {
		zap.L().Error("ping redis Error:", zap.Error(err))
		return
	}
	fmt.Printf("%s! %s\n", pong, "Redis Connect Success!")
	return
}

func CloseRedis() {
	if RDB != nil {
		RDB.Close()
	}
}
