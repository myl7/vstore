package dao

import "github.com/go-redis/redis/v8"

var kv *redis.Client

func GetKV() *redis.Client {
	if kv != nil {
		return kv
	}
	kv = redis.NewClient(&redis.Options{})
	return kv
}
