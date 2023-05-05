// Package redisc .
 
package redisc

import (
	"strconv"

	"github.com/go-redis/redis"
)

// Pipeline pipeline
type Pipeline struct {
	redis.Pipeliner
}

// HMSetObject 操作struct
func (rp *Pipeline) HMSetObject(key string, val interface{}, expire int) *redis.StatusCmd {
	args := Flat(val, RedisCmdHMSET, key)
	cmd := redis.NewStatusCmd(args...)
	_ = rp.Process(cmd)
	return cmd
}

// HMGetObject 获取struct
func (rp *Pipeline) HMGetObject(key string, val interface{}) (interface{}, error) {
	res, err := rp.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}
	return ConvertTo(res, val)
}

// ZRevRangeMinScore ZRevRangeMinScore
func (rp *Pipeline) ZRevRangeMinScore(key string, min, offset, count int64) *redis.StringSliceCmd {
	opt := redis.ZRangeBy{
		Min:    strconv.FormatInt(min, 10),
		Max:    "+inf",
		Offset: offset,
		Count:  count,
	}
	return rp.ZRevRangeByScore(key, opt)
}

// ZRevRangeMaxScoreWithScores (-inf, max]
func (rp *Pipeline) ZRevRangeMaxScoreWithScores(key string, max, offset, count int64) *redis.ZSliceCmd {
	opt := redis.ZRangeBy{
		Min:    "-inf",
		Max:    strconv.FormatInt(max, 10),
		Offset: offset,
		Count:  count,
	}
	return rp.ZRevRangeByScoreWithScores(key, opt)
}
