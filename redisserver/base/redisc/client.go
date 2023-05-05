// Package redisc .

package redisc

import (
	"strconv"
	"sync"
	"time"

	"redisserver/base/conf"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

// 全局变量
var (
	once           sync.Once
	redisSingleton *Client
)

// Client 客户端
type Client struct {
	client *redis.Client
}

// DC 默认连接
func DC() *Client {
	once.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:        conf.RedisDBUrl,
			PoolTimeout: 5 * time.Second,
		})

		pong, err := client.Ping().Result()
		if err != nil {
			logrus.Errorf("redis host %v  ping %v error %v", conf.RedisDBUrl, pong, err)
		}

		redisSingleton = &Client{client: client}
	})
	return redisSingleton
}

// NewClient 创建一个连接
func NewClient(addr, password string, db int) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		logrus.Errorf("redis host %v  ping %v error %v", addr, pong, err)
	}

	return &Client{client: client}, err
}

// Close ..
func (r *Client) Close() {
	r.client.Close()
}

// PoolStats 连接池状态
func (r *Client) PoolStats() *redis.PoolStats {
	return r.client.PoolStats()
}

// Pipeline Pipeline
func (r *Client) Pipeline() *Pipeline {
	rp := &Pipeline{Pipeliner: r.client.Pipeline()}
	return rp
}

// TxPipeline TxPipeline
func (r *Client) TxPipeline() *Pipeline {
	rp := &Pipeline{Pipeliner: r.client.TxPipeline()}
	return rp
}

// FlushDB 清除DB
func (r *Client) FlushDB() {
	_, _ = r.client.FlushDB().Result()
}

// SelectDB SelectDB
func (r *Client) SelectDB(index int32) error {
	cmd := redis.NewStatusCmd("select", index)
	if err := r.client.Process(cmd); err != nil {
		return err
	}
	_, err := cmd.Result()
	return err
}

// BgSave BgSave
func (r *Client) BgSave() (string, error) {
	return r.client.BgSave().Result()
}

// Exists Exists
func (r *Client) Exists(key string) (int64, error) {
	return r.client.Exists(key).Result()
}

// Get Get
func (r *Client) Get(key string) (string, error) {
	return r.client.Get(key).Result()
}

// Set Set
func (r *Client) Set(key string, value interface{}, expiration time.Duration) (string, error) {
	return r.client.Set(key, value, expiration).Result()
}

// SetNX SetNX
func (r *Client) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.client.SetNX(key, value, expiration).Result()
}

// Del Del
func (r *Client) Del(key string) (int64, error) {
	return r.client.Del(key).Result()
}

// Incr Incr
func (r *Client) Incr(key string) (int64, error) {
	return r.client.Incr(key).Result()
}

// LRem LRem
func (r *Client) LRem(key string, count int64, value interface{}) (int64, error) {
	return r.client.LRem(key, count, value).Result()
}

// LPush LPush
func (r *Client) LPush(key string, values ...interface{}) (int64, error) {
	return r.client.LPush(key, values...).Result()
}

// LPop LPop
func (r *Client) LPop(key string, values ...interface{}) (string, error) {
	return r.client.LPop(key).Result()
}

// RPush RPush
func (r *Client) RPush(key string, values ...interface{}) (int64, error) {
	return r.client.RPush(key, values...).Result()
}

// RPop RPop
func (r *Client) RPop(key string) (string, error) {
	return r.client.RPop(key).Result()
}

// LLen LLen
func (r *Client) LLen(key string) (int64, error) {
	return r.client.LLen(key).Result()
}

// LRange LRange
func (r *Client) LRange(key string, start, stop int64) ([]string, error) {
	return r.client.LRange(key, start, stop).Result()
}

// LTrim LTrim
func (r *Client) LTrim(key string, start, stop int64) (string, error) {
	return r.client.LTrim(key, start, stop).Result()
}

// HMSet HMSet
func (r *Client) HMSet(key string, fields map[string]interface{}) bool {
	_, err := r.client.HMSet(key, fields).Result()
	return err == nil
}

// HMSetObject  HMSetObject
func (r *Client) HMSetObject(key string, val interface{}, expire int) error {
	args := Flat(val, RedisCmdHMSET, key)
	cmd := redis.NewStatusCmd(args...)
	if err := r.client.Process(cmd); err != nil {
		return err
	}
	_, err := cmd.Result()
	return err
}

// HMGet HMGet
func (r *Client) HMGet(key string, fields ...string) ([]interface{}, error) {
	return r.client.HMGet(key, fields...).Result()
}

// HMGetObject HMGetObject
func (r *Client) HMGetObject(key string, val interface{}) (interface{}, error) {
	res, err := r.client.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}
	return ConvertTo(res, val)
}

// HSet HSet
func (r *Client) HSet(key, field string, value interface{}) bool {
	_, err := r.client.HSet(key, field, value).Result()
	if err != nil {
		logrus.Errorf("HExists  failed key %s, field %s, err: %s", key, field, err.Error())
		return false
	}

	return true
}

// HGet HGet
func (r *Client) HGet(key, member string) (string, error) {
	return r.client.HGet(key, member).Result()
}

// HDel HDel
func (r *Client) HDel(key, member string) (reply interface{}, err error) {
	return r.client.HDel(key, member).Result()
}

// HExists HExists
func (r *Client) HExists(key, field string) bool {
	ok, err := r.client.HExists(key, field).Result()
	if err != nil {
		logrus.Errorf("HExists  failed key %s, field %s, err: %s", key, field, err.Error())
		return false
	}

	return ok
}

// HSetNX HSetNX
func (r *Client) HSetNX(key, field, value string) bool {
	ok, err := r.client.HSetNX(key, field, value).Result()
	if err != nil {
		logrus.Errorf("HSetnx write failed key %s, field %s, value %s, err: %s", key, field, value, err.Error())
		return false
	}

	return ok
}

// HGetAll HGetAll
func (r *Client) HGetAll(key string) (map[string]string, error) {
	return r.client.HGetAll(key).Result()
}

// HIncrBy HIncrBy
func (r *Client) HIncrBy(key, field string, incr int64) (int64, error) {
	return r.client.HIncrBy(key, field, incr).Result()
}

// SAdd 增加
func (r *Client) SAdd(key string, members ...interface{}) (int64, error) {
	return r.client.SAdd(key, members...).Result()
}

// SRem 删除
func (r *Client) SRem(key string, members ...interface{}) (int64, error) {
	return r.client.SRem(key, members...).Result()
}

// SIsMember 成员
func (r *Client) SIsMember(key string, member interface{}) (bool, error) {
	return r.client.SIsMember(key, member).Result()
}

// ZAdd 将一个成员元素及其分数值加入到有序集当中
func (r *Client) ZAdd(key string, members ...redis.Z) (reply interface{}, err error) {
	return r.client.ZAdd(key, members...).Result()
}

// ZRem 移除有序集中的一个，不存在的成员将被忽略。
func (r *Client) ZRem(key string, member interface{}) (reply interface{}, err error) {
	return r.client.ZRem(key, member).Result()
}

// ZRemRangeByScore 根据Score移除有序集中的元素。
func (r *Client) ZRemRangeByScore(key, min, max string) (reply interface{}, err error) {
	return r.client.ZRemRangeByScore(key, min, max).Result()
}

// ZRemRangeByRank 根据Rank移除有序集中的元素。
func (r *Client) ZRemRangeByRank(key string, min, max int64) (reply interface{}, err error) {
	return r.client.ZRemRangeByRank(key, min, max).Result()
}

// ZExist 判断元素是否在有序集合中
func (r *Client) ZExist(key, member string) (bool, error) {
	_, err := r.client.ZRank(key, member).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ZScore 返回有序集中，成员的分数值。 如果成员元素不是有序集 key 的成员，或 key 不存在，返回 nil
func (r *Client) ZScore(key, member string) int64 {
	val, _ := r.client.ZScore(key, member).Result()
	return int64(val)
}

// ZIncrBy 添加单位成员的有序集合存储增量键比分
func (r *Client) ZIncrBy(key string, score int64, member string) (int64, error) {
	val, err := r.client.ZIncrBy(key, float64(score), member).Result()
	return int64(val), err
}

// ZRank 返回有序集中指定成员的排名。其中有序集成员按分数值递增(从小到大)顺序排列。score 值最小的成员排名为 0
func (r *Client) ZRank(key, member string) (int64, error) {
	return r.client.ZRank(key, member).Result()
}

// ZCount 返回集合数量
func (r *Client) ZCount(key, min, max string) (int64, error) {
	return r.client.ZCount(key, min, max).Result()
}

// ZRevRank 返回有序集中成员的排名。其中有序集成员按分数值递减(从大到小)排序。分数值最大的成员排名为 0 。
func (r *Client) ZRevRank(key, member string) (int64, error) {
	return r.client.ZRevRank(key, member).Result()
}

// ZRange 返回有序集中，指定区间内的成员。其中成员的位置按分数值递增(从小到大)来排序。具有相同分数值的成员按字典序(lexicographical order )来排列。
// 以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。或 以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
func (r *Client) ZRange(key string, from, to int64) ([]string, error) {
	values, err := r.client.ZRange(key, from, to).Result()
	return values, err
}

// ZRangeByScore 返回有序集合中指定分数区间的成员列表。有序集成员按分数值递增(从小到大)次序排列。
// 具有相同分数值的成员按字典序来排列
func (r *Client) ZRangeByScore(key string, max, min, offset, count int64) ([]string, error) {
	opt := redis.ZRangeBy{
		Min:    strconv.FormatInt(min, 10),
		Max:    strconv.FormatInt(max, 10),
		Offset: offset,
		Count:  count,
	}
	values, err := r.client.ZRangeByScore(key, opt).Result()
	return values, err
}

// ZRangeByScoreWithScores ....
func (r *Client) ZRangeByScoreWithScores(key string, max, min, offset, count int64) ([]redis.Z, error) {
	opt := redis.ZRangeBy{
		Min:    strconv.FormatInt(min, 10),
		Max:    strconv.FormatInt(max, 10),
		Offset: offset,
		Count:  count,
	}
	values, err := r.client.ZRangeByScoreWithScores(key, opt).Result()
	return values, err
}

// ZRevRange 返回有序集中，指定区间内的成员。其中成员的位置按分数值递减(从大到小)来排列。具有相同分数值的成员按字典序(lexicographical order )来排列。
// 以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。或 以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
func (r *Client) ZRevRange(key string, from, to int64) ([]string, error) {
	values, err := r.client.ZRevRange(key, from, to).Result()
	return values, err
}

// ZRevRangeByScore 返回有序集中指定分数区间内的所有的成员。有序集成员按分数值递减(从大到小)的次序排列。
// 具有相同分数值的成员按字典序来排列
func (r *Client) ZRevRangeByScore(key string, max, min, offset, count int64) ([]string, error) {
	opt := redis.ZRangeBy{
		Min:    strconv.FormatInt(min, 10),
		Max:    strconv.FormatInt(max, 10),
		Offset: offset,
		Count:  count,
	}
	values, err := r.client.ZRevRangeByScore(key, opt).Result()
	return values, err
}

// ZRevRangeByScoreWithScores ..
func (r *Client) ZRevRangeByScoreWithScores(key string, max, min, offset, count int64) ([]redis.Z, error) {
	opt := redis.ZRangeBy{
		Min:    strconv.FormatInt(min, 10),
		Max:    strconv.FormatInt(max, 10),
		Offset: offset,
		Count:  count,
	}
	values, err := r.client.ZRevRangeByScoreWithScores(key, opt).Result()
	return values, err
}

// ZRevRangeMinScore [MinScore, +inf)
func (r *Client) ZRevRangeMinScore(key string, min, offset, count int64) ([]string, error) {
	opt := redis.ZRangeBy{
		Min:    strconv.FormatInt(min, 10),
		Max:    "+inf",
		Offset: offset,
		Count:  count,
	}
	values, err := r.client.ZRevRangeByScore(key, opt).Result()
	return values, err
}

// ZRevRangeMinScoreWithScores [MinScore, +inf)
func (r *Client) ZRevRangeMinScoreWithScores(key string, min, offset, count int64) ([]redis.Z, error) {
	opt := redis.ZRangeBy{
		Min:    strconv.FormatInt(min, 10),
		Max:    "+inf",
		Offset: offset,
		Count:  count,
	}
	values, err := r.client.ZRevRangeByScoreWithScores(key, opt).Result()
	return values, err
}

// ZRevRangeMaxScore (-inf, max]
func (r *Client) ZRevRangeMaxScore(key string, max, offset, count int64) ([]string, error) {
	opt := redis.ZRangeBy{
		Min:    "-inf",
		Max:    strconv.FormatInt(max, 10),
		Offset: offset,
		Count:  count,
	}
	values, err := r.client.ZRevRangeByScore(key, opt).Result()
	return values, err
}

// ZRevRangeMaxScoreWithScores (-inf, max]
func (r *Client) ZRevRangeMaxScoreWithScores(key string, max, offset, count int64) ([]redis.Z, error) {
	opt := redis.ZRangeBy{
		Min:    "-inf",
		Max:    strconv.FormatInt(max, 10),
		Offset: offset,
		Count:  count,
	}
	values, err := r.client.ZRevRangeByScoreWithScores(key, opt).Result()
	return values, err
}

// ZCard 获取set中元素数量
func (r *Client) ZCard(key string) (int64, error) {
	values, err := r.client.ZCard(key).Result()
	return values, err
}

// Publish 将信息发送到指定的频道，返回接收到信息的订阅者数量
func (r *Client) Publish(channel, message string) (int64, error) {
	return r.client.Publish(channel, message).Result()
}

// Subscribe 订阅
func (r *Client) Subscribe(channels ...string) *redis.PubSub {
	return r.client.Subscribe(channels...)
}
