package dredis

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

// RedisClusterDriver is
type RedisClusterDriver struct {
	cluster bool
	client  *redis.ClusterClient
	timeout time.Duration
	ctx     context.Context
}

func NewClusterClient(opts *redis.ClusterOptions) (c *RedisClusterDriver, err error) {
	c = &RedisClusterDriver{
		cluster: true,
		client:  redis.NewClusterClient(opts),
		ctx:     context.TODO(),
	}
	err = c.Ping()
	return c, err
}

func (rd *RedisClusterDriver) Client() (bool, *redis.Client, *redis.ClusterClient) {
	return rd.cluster, nil, rd.client
}

func (rd *RedisClusterDriver) Ping() error {
	if err := rd.client.Set(rd.ctx, "ping", "pong", 0).Err(); err != nil {
		return err
	}
	return nil
}

// CacheGet 获取指定key的值,如果值不存在,就执行f方法将返回值存入redis
func (rd *RedisClusterDriver) CacheGet(key string, expiration time.Duration, f func() (res string, err error)) (string, error) {
	result, _ := rd.client.Get(rd.ctx, key).Result()
	if len(result) == 0 {
		res, err := f()
		if err == nil {
			rd.client.Set(rd.ctx, key, res, expiration)
		}
		return res, err
	}
	return result, nil
}

func (rd *RedisClusterDriver) Get(key string) (string, error) {
	value, err := rd.client.Get(rd.ctx, key).Result()
	if err != nil {
		return "", app.Err(app.Fail, "redis get key: %s err", key)
	}
	return value, nil
}

func (rd *RedisClusterDriver) Set(key string, value interface{}, expiration time.Duration) error {
	if err := rd.client.Set(rd.ctx, key, value, expiration).Err(); err != nil {
		return app.Err(app.Fail, "redis set key: %s err:%s", key, err.Error())
	}
	return nil
}

func (rd *RedisClusterDriver) Del(key string) (int64, error) {
	return rd.client.Del(rd.ctx, key).Result()
}

func (rd *RedisClusterDriver) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return rd.client.SetNX(rd.ctx, key, value, expiration).Result()
}

func (rd *RedisClusterDriver) LPush(key string, values ...interface{}) (int64, error) {
	return rd.client.LPush(rd.ctx, key, values).Result()
}

func (rd *RedisClusterDriver) LPushX(key string, values ...interface{}) (int64, error) {
	return rd.client.LPushX(rd.ctx, key, values).Result()
}

func (rd *RedisClusterDriver) Eval(script string, keys []string, args ...interface{}) (interface{}, error) {
	return rd.client.Eval(rd.ctx, script, keys, args).Result()
}

func (rd *RedisClusterDriver) Do(args ...interface{}) (interface{}, error) {
	return rd.client.Do(rd.ctx, args).Result()
}

func (rd *RedisClusterDriver) Expire(key string, expiration time.Duration) (bool, error) {
	return rd.client.Expire(rd.ctx, key, expiration).Result()
}

func (rd *RedisClusterDriver) ExpireAt(key string, ttl time.Time) (bool, error) {
	return rd.client.ExpireAt(rd.ctx, key, ttl).Result()
}

func (rd *RedisClusterDriver) TTL(key string) (time.Duration, error) {
	ttl, err := rd.client.TTL(rd.ctx, key).Result()
	if err != nil {
		return -1, app.Err(app.Fail, "redis get key: %s err", key)
	}
	return ttl, nil
}

func (rd *RedisClusterDriver) Exists(keys ...string) bool {
	if len(keys) == 0 {
		return true
	}
	value, _ := rd.client.Exists(rd.ctx, keys...).Result()
	return value > 0
}

func (rd *RedisClusterDriver) Incr(key string) int64 {
	value, _ := rd.client.Incr(rd.ctx, key).Result()
	return value
}

func (rd *RedisClusterDriver) Close() error {
	return rd.client.Close()
}

// ZAdd 将一个 member 元素及其 score 值加入到有序集 key 当中。
func (rd *RedisClusterDriver) ZAdd(key string, score float64, member string) (reply int64, err error) {
	return rd.client.ZAdd(rd.ctx, key, &redis.Z{
		Score:  score,
		Member: member,
	}).Result()
}

// ZIncrby 命令对有序集合中指定成员的分数加上增量 increment
func (rd *RedisClusterDriver) ZIncrby(key, member string, increment float64) (float64, error) {
	return rd.client.ZIncrBy(rd.ctx, key, increment, member).Result()
}

// ZRange 返回有序集中，指定区间内的成员。其中成员的位置按分数值递增(从小到大)来排序。具有相同分数值的成员按字典序(lexicographical order )来排列。
func (rd *RedisClusterDriver) ZRange(key string, start, end int64) ([]redis.Z, error) {
	return rd.client.ZRangeWithScores(rd.ctx, key, start, end).Result()
}

// ZRem 移除有序集 key 中的一个成员，不存在的成员将被忽略。
func (rd *RedisClusterDriver) ZRem(key string, member string) (reply int64, err error) {
	return rd.client.ZRem(rd.ctx, key, member).Result()
}

// ZCard 获取有序集合的成员数
func (rd *RedisClusterDriver) ZCard(key string) int64 {
	return rd.client.ZCard(rd.ctx, key).Val()
}

// ZScore 返回有序集 key 中，成员 member 的 score 值。 如果 member 元素不是有序集 key 的成员，或 key 不存在，返回 nil 。
func (rd *RedisClusterDriver) ZScore(key string, member string) (float64, error) {
	return rd.client.ZScore(rd.ctx, key, member).Result()
}

func (rd *RedisClusterDriver) BRPop(timeout time.Duration, keys ...string) ([]string, error) {
	return rd.client.BRPop(rd.ctx, timeout, keys...).Result()
}

// HGetAll 所有的域和值
func (rd *RedisClusterDriver) HGetAll(key string) (map[string]string, error) {
	return rd.client.HGetAll(rd.ctx, key).Result()
}

// HMGet 返回哈希表中
func (rd *RedisClusterDriver) HMGet(key string, fields ...string) ([]interface{}, error) {
	return rd.client.HMGet(rd.ctx, key, fields...).Result()
}

// HIncrBy 为哈希表中的字段值加上指定增量值
func (rd *RedisClusterDriver) HIncrBy(key, field string, incr int64) (int64, error) {
	return rd.client.HIncrBy(rd.ctx, key, field, incr).Result()
}

// HKeys 命令用于获取哈希表中的所有字段名
func (rd *RedisClusterDriver) HKeys(key string) ([]string, error) {
	return rd.client.HKeys(rd.ctx, key).Result()
}

// HDel 删除哈希表 key 中的一个或多个指定字段，不存在的字段将被忽略
func (rd *RedisClusterDriver) HDel(key string, fields ...string) (int64, error) {
	return rd.client.HDel(rd.ctx, key, fields...).Result()
}

// ZRangeByScore 返回有序集合中指定分数区间的成员列表
func (rd *RedisClusterDriver) ZRangeByScore(key string, opt *redis.ZRangeBy) ([]string, error) {
	return rd.client.ZRangeByScore(rd.ctx, key, opt).Result()
}

// ZRemRangeByScore 命令用于移除有序集中，指定分数（score）区间内的所有成员
func (rd *RedisClusterDriver) ZRemRangeByScore(key, min, max string) (int64, error) {
	return rd.client.ZRemRangeByScore(rd.ctx, key, min, max).Result()
}

// ZCount 返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。
func (rd *RedisClusterDriver) ZCount(key, min, max string) (int64, error) {
	return rd.client.ZCount(rd.ctx, key, min, max).Result()
}

// HSet 希表中的字段赋值
func (rd *RedisClusterDriver) HSet(key string, values ...interface{}) (int64, error) {
	if len(values) == 1 {
		return rd.client.HSet(rd.ctx, key, values[0]).Result()
	}
	return rd.client.HSet(rd.ctx, key, values).Result()
}

// Scan 搜索
func (rd *RedisClusterDriver) Scan(cursor uint64, match string, count int64) ([]string, uint64, error) {
	return rd.client.Scan(rd.ctx, cursor, match, count).Result()
}

// HScan 搜索
func (rd *RedisClusterDriver) HScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return rd.client.HScan(rd.ctx, key, cursor, match, count).Result()
}

// SAdd 向集合添加一个或多个成员
func (rd *RedisClusterDriver) SAdd(key string, args ...interface{}) (int64, error) {
	return rd.client.SAdd(rd.ctx, key, args).Result()
}

// SCard 获取集合的成员数
func (rd *RedisClusterDriver) SCard(key string) (int64, error) {
	return rd.client.SCard(rd.ctx, key).Result()
}

// SIsMember 判断 member 元素是否是集合 key 的成员
func (rd *RedisClusterDriver) SIsMember(key string, member interface{}) bool {
	return rd.client.SIsMember(rd.ctx, key, member).Val()
}

// SRandMember 返回集合中的一个随机元素。
func (rd *RedisClusterDriver) SRandMember(key string) string {
	return rd.client.SRandMember(rd.ctx, key).Val()
}

// SMembers 命令返回集合中的所有的成员。 不存在的集合 key 被视为空集合
func (rd *RedisClusterDriver) SMembers(key string) ([]string, error) {
	return rd.client.SMembers(rd.ctx, key).Result()
}

// SRem 用于移除集合中的一个或多个成员元素，不存在的成员元素会被忽略。
func (rd *RedisClusterDriver) SRem(key string, args ...interface{}) (int64, error) {
	return rd.client.SRem(rd.ctx, key, args).Result()
}

// Subscribe 订阅给定的一个或多个频道的信息。 每个模式以 * 作为匹配符
func (rd *RedisClusterDriver) Subscribe(channels ...string) *redis.PubSub {
	return rd.client.Subscribe(rd.ctx, channels...)
}

// PSubscribe 订阅一个或多个符合给定模式的频道。
func (rd *RedisClusterDriver) PSubscribe(channels ...string) *redis.PubSub {
	return rd.client.PSubscribe(rd.ctx, channels...)
}

// Publish 将信息 message 发送到指定的频道 channel
func (rd *RedisClusterDriver) Publish(channel string, message string) (int64, error) {
	return rd.client.Publish(rd.ctx, channel, message).Result()
}
