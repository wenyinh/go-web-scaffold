package verify

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go-web-scaffold/dao/redis"
	"go.uber.org/zap"
	"time"
)

var ctx = context.Background()

func TestRedis() (err error) {
	uid := uuid.New().String()
	key := fmt.Sprint(uid[:8])
	val := fmt.Sprintf("%s %s", uid, "PONG")
	exp := 5 * time.Second
	if err = redis.RDB.Set(key, val, exp).Err(); err != nil {
		zap.L().Error("Redis.Set failed", zap.String("key", key), zap.Error(err))
		return
	}
	fmt.Println("Redis Set Success")
	result, err := redis.RDB.Get(key).Result()
	if err != nil {
		zap.L().Error("Redis.Get failed", zap.String("key", key), zap.Error(err))
		return
	}
	if result != val {
		return errors.New("Redis value mismatch")
	}
	fmt.Println("Redis Get Success, Pass Redis Test")
	return
}
