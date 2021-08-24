package gredis

import (
	"cloud-batch/configs"
	"fmt"
	"gopkg.in/redis.v5"
	"time"
)

var Rdb *redis.Client

//打开数据库
func init() {
	if Rdb == nil {
		Rdb = getRedisClient()
	}
}

func getRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     fmt.Sprintf("%s:%d", configs.Redis.Host, configs.Redis.Port),
		DB:       configs.Redis.DB,
		Password: configs.Redis.Pass,
	})
}

//func OpenQueue(name string) rmq.Queue {
//	conn := rmq.OpenConnectionWithRedisClient("gw", Rdb)
//	return conn.OpenQueue(name)
//}

func Exists(key string) *redis.BoolCmd {
	return Rdb.Exists(key)
}

func Get(key string) *redis.StringCmd {
	return Rdb.Get(key)
}

func Set(key string, value string, expiration time.Duration) error {
	var err error
	// set失败重试3次
	for i := 0; i <= 2; i++ {
		err = Rdb.Set(key, value, expiration).Err()
		if err == nil {
			break
		}
	}
	return err
}

func Del(key string) error {
	return Rdb.Del(key).Err()
}

func DelKeys(keys []string) *redis.IntCmd {
	args := make([]interface{}, 1+len(keys))
	args[0] = "DEL"
	for i, key := range keys {
		args[1+i] = key
	}
	cmd := redis.NewIntCmd(args...)
	Rdb.Process(cmd)
	return cmd
}

func TTL(key string) *redis.DurationCmd {
	return Rdb.TTL(key)
}

func Keys(pattern string) *redis.StringSliceCmd {
	return Rdb.Keys(pattern)
}

func Incr(key string) *redis.IntCmd {
	return Rdb.Incr(key)
}

func Expire(key string, expiration time.Duration) *redis.BoolCmd {
	ret := Rdb.Expire(key, expiration)
	// 失败重试3次
	for i := 0; i < 2; i++ {
		if ret.Err() == nil {
			break
		}
		ret = Rdb.Expire(key, expiration)
	}
	return ret
}

// 有ttl 的，删除ttl, 设置为永不过期
func Persist(key string) *redis.BoolCmd {
	return Rdb.Persist(key)
}

// 向集合添加一个成员
func SAdd(key string, value string) *redis.IntCmd {
	ret := Rdb.SAdd(key, value)
	// 失败重试3次
	for i := 0; i < 2; i++ {
		if ret.Err() == nil {
			break
		}
		ret = Rdb.SAdd(key, value)
	}
	return ret
}

// 获取集合的成员数
func SCard(key string) *redis.IntCmd {
	return Rdb.SCard(key)
}

// 返回集合中的所有成员
func SMembers(key string) *redis.StringSliceCmd {
	// set失败重试3次
	ret := Rdb.SMembers(key)
	i := 0
	for {
		if ret.Err() == nil || i == 2 {
			return ret
		}
		ret = Rdb.SMembers(key)
		i++
	}
}

// 删除集合中的一个成员
func SRem(key, value string) *redis.IntCmd {
	return Rdb.SRem(key, value)
}
