package database

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"
)

type RedisDB struct {
	client  *redis.Client
	context context.Context
}

func NewRedisDB(context context.Context, client *redis.Client) *RedisDB {
	rdb := new(RedisDB)
	rdb.client = client
	rdb.context = context

	return rdb
}

func (r *RedisDB) Set(key string, value interface{}) (string, error) {
	_, err := r.client.Set(r.context, key, value, 0).Result()
	return key, err
}

func (r *RedisDB) Get(key string, resultType string) (interface{}, error) {
	result, err := r.client.Get(r.context, key).Result()
	if err != nil {
		return result, err
	}

	switch resultType {
	case "uint":
		temp, err := strconv.ParseUint(result, 10, 32)
		if err != nil {
			glog.Fatal(err)
			return temp, err
		}
		return uint(temp), nil
	default:
		return result, err
	}
}

func (r *RedisDB) Delete(key string) (string, error) {
	_, err := r.client.Del(r.context, key).Result()
	return key, err
}

func (r *RedisDB) Incr(key string) (int64, error) {
	return r.client.Incr(r.context, key).Result()
}

func (r *RedisDB) SAdd(key string, member string) (int64, error) {
	return r.client.SAdd(r.context, key, member).Result()
}

func (r *RedisDB) SMembers(key string) ([]string, error) {
	return r.client.SMembers(r.context, key).Result()
}

func (r *RedisDB) LPush(key string, value string) (int64, error) {
	return r.client.LPush(r.context, key, value).Result()
}

func (r *RedisDB) LPop(key string) (string, error) {
	return r.client.LPop(r.context, key).Result()
}
