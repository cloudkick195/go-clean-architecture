package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

var RedisStore *redis.Client

func InitRedis() {
	redisStore := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Env.REDIS_ADDRESS, Env.REDIS_PORT),
		Password: Env.REDIS_PASSWORD,
		DB:       0, // use default DB
	})

	_, err := redisStore.Ping().Result()
	if err != nil {
		panic(err)
	}
	RedisStore = redisStore
}

func RedisGet(key string) (string, error) {
	val, err := RedisStore.Get(key).Result()
	if err == redis.Nil {
		return "", errors.New("key not found")
	}
	return val, err
}

func RedisSet(key string, value interface{}, expiration time.Duration) error {
	err := RedisStore.Set(key, value, expiration).Err()
	return err
}

func RedisDel(key string) error {
	return RedisStore.Del(key).Err()
}
