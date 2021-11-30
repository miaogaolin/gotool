package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/miaogaolin/gotool/logx"

	"github.com/go-redis/redis"
)

var (
	prefix = "syt-crawler:"
	// Maximum number of socket connections.
	// Default is 10 connections per every CPU as reported by runtime.NumCPU.
	poolSize = 100
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	minIdleConns = 10
)

type Z struct {
	Score  float64
	Member interface{}
}

var DB *RedisDB

type RedisDB struct {
	client *redis.Client
	config *Config
	prefix string
}

func Connect(config *Config) (*RedisDB, error) {
	if config.Enabled == 0 {
		DB = &RedisDB{config: config}
		return DB, nil
	}
	if config.PoolSize > 0 {
		poolSize = config.PoolSize
	}

	if config.MinIdleConns > 0 {
		minIdleConns = config.MinIdleConns
	}

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.Db, // use default DB
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
	})
	_, err := client.Ping().Result()
	DB = &RedisDB{client: client, config: config, prefix: prefix}
	return DB, err
}

func (r *RedisDB) SetPrefix(prefix string) *RedisDB {
	r.prefix = prefix
	return r
}

// 获取
func (r *RedisDB) Get(key string) (string, error) {
	key = r.getKey(key)
	res, err := r.client.Get(key).Result()
	if err == redis.Nil {
		logx.Errorf("redis operation=%s key=%s not exist\n", "get", key)
		return "", nil
	}
	return res, err
}

// 包装Get方法
// 从json字符串解析为go语言变量类型
func (r *RedisDB) GetJson(key string, out interface{}) error {
	key = r.getKey(key)
	result, err := r.client.Get(key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(result), out)
	return err
}

// 判断key是否存在
func (r *RedisDB) Exist(key string) bool {
	if r.config.Enabled == 0 {
		return false
	}
	key = r.getKey(key)
	if count, _ := r.client.Exists(key).Result(); count > 0 {
		return true
	}
	return false
}

func (r *RedisDB) HGetJson(key string, field string, out interface{}) error {
	key = r.getKey(key)
	result, err := r.client.HGet(key, field).Result()
	if err == redis.Nil {
		logx.Errorf("redis operation=%s key=%s field=%s is not exist\n", "hget", key, field)
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(result), out)
}

func (r *RedisDB) HExist(key string, field string) bool {
	if r.config.Enabled == 0 {
		return false
	}
	key = r.getKey(key)
	if ok, _ := r.client.HExists(key, field).Result(); ok {
		return true
	}
	return false
}

func (r *RedisDB) HSetJson(key string, field string, value interface{}) *redis.BoolCmd {
	if r == nil || r.config.Enabled == 0 {
		return nil
	}

	key = r.getKey(key)
	res, _ := json.Marshal(value)
	return r.client.HSet(key, field, res)
}

func (r *RedisDB) Del(keys ...string) *redis.IntCmd {
	if r.config.Enabled == 0 {
		return nil
	}
	return r.client.Del(keys...)
}

func (r *RedisDB) HDel(key string, fields ...string) *redis.IntCmd {
	if r.config.Enabled == 0 {
		return nil
	}
	key = r.getKey(key)
	return r.client.HDel(key, fields...)
}

func (r *RedisDB) HGetAll(key string) (map[string]string, error) {
	if r.config.Enabled == 0 {
		return nil, nil
	}
	key = r.getKey(key)
	result, err := r.client.HGetAll(key).Result()
	if err == redis.Nil {
		logx.Errorf("redis operation=%s key=%s is not exist\n", "hgetall", key)
		return nil, nil
	}
	return result, nil
}

func (r *RedisDB) HMGet(key string, fields ...string) ([]interface{}, error) {
	if r.config.Enabled == 0 {
		return nil, nil
	}

	if len(fields) == 0 {
		return nil, nil
	}
	key = r.getKey(key)
	result, err := r.client.HMGet(key, fields...).Result()
	if err == redis.Nil {
		logx.Errorf("redis operation=%s key=%s fields=%v not exist\n", "hmget", key, fields)
		return nil, nil
	}
	return result, nil
}

// 是否启用缓存
func (r *RedisDB) IsEnabled() bool {
	return r.config.Enabled == 1
}

func (r *RedisDB) SMembers(key string) ([]string, error) {
	if r.config.Enabled == 0 {
		return nil, nil
	}

	key = r.getKey(key)
	result, err := r.client.SMembers(key).Result()
	if err == redis.Nil {
		logx.Errorf("redis operation=%s key=%s not exist\n", "smembers", key)
		return nil, nil
	}
	return result, err
}

func (r *RedisDB) ZRank(key string, member string) (int64, error) {
	if r.config.Enabled == 0 {
		return 0, nil
	}

	key = r.getKey(key)
	return r.client.ZRank(key, member).Result()
}

func (r *RedisDB) ZCard(key string) int64 {
	if r.config.Enabled == 0 {
		return 0
	}
	key = r.getKey(key)
	count, _ := r.client.ZCard(key).Result()
	return count
}

func (r *RedisDB) ZRange(key string, start int64, stop int64) ([]string, error) {
	if r.config.Enabled == 0 {
		return nil, nil
	}

	key = r.getKey(key)
	result, err := r.client.ZRange(key, start, stop).Result()
	if err == redis.Nil {
		logx.Errorf("redis operation=%s key=%s not exist\n", "zrange", key)
		return nil, nil
	}
	return result, err
}

func (r *RedisDB) SAdd(key string, members ...interface{}) (int64, error) {
	if r.config.Enabled == 0 {
		return 0, nil
	}
	key = r.getKey(key)
	result, err := r.client.SAdd(key, members).Result()
	if err == redis.Nil {
		logx.Errorf("redis operation=%s key=%s not exist\n", "sadd", key)
		return 0, nil
	}
	return result, err
}

func (r *RedisDB) ZAddJson(key string, members ...Z) (int64, error) {
	if r.config.Enabled == 0 {
		return 0, nil
	}

	key = r.getKey(key)
	var z []redis.Z
	for i := range members {
		b, _ := json.Marshal(members[i].Member)
		z = append(z, redis.Z{
			Score:  members[i].Score,
			Member: b,
		})
	}

	result, err := r.client.ZAdd(key, z...).Result()
	if err == redis.Nil {
		logx.Errorf("redis operation=%s key=%s not exist\n", "sadd", key)
		return 0, nil
	}
	return result, err
}

func (r *RedisDB) Expire(key string, expiration time.Duration) (bool, error) {
	if r.config.Enabled == 0 {
		return false, nil
	}
	key = r.getKey(key)
	return r.client.Expire(key, expiration).Result()
}

func (r *RedisDB) TTL(key string) (time.Duration, error) {
	if r.config.Enabled == 0 {
		return 0, nil
	}
	return r.client.TTL(key).Result()
}

func (r *RedisDB) Set(key string, value interface{}, expiration time.Duration) (string, error) {
	if r.config.Enabled == 0 {
		return "", nil
	}
	key = r.getKey(key)
	res, err := r.client.Set(key, value, expiration).Result()
	if err == redis.Nil {
		logx.Errorf("redis operation=%s key=%s not exist\n", "set", key)
		return "", nil
	}
	return res, err
}

func (r *RedisDB) Eval(script string, keys []string, args ...interface{}) (interface{}, error) {
	if r.config.Enabled == 0 {
		return nil, nil
	}
	for i := range keys {
		keys[i] = r.getKey(keys[i])
	}
	return r.client.Eval(script, keys, args...).Result()
}

func (r *RedisDB) getKey(key string) string {
	return prefix + key
}
