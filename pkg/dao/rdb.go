package dao

import "github.com/go-redis/redis/v8"

func GetKV() *redis.Client {
	kv := redis.NewClient(&redis.Options{})
	return kv
}
