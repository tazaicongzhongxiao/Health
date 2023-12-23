// Package dredis
// @Description: REDIS 分布式并发锁
package dredis

import (
	"encoding/base64"
	"errors"
	"math/rand"
	"sync"
	"time"
)

const (
	DefaultRetry        = 5                      //默认重试
	DefaultInterval     = 100 * time.Millisecond //默认间隔
	DefaultExpire       = 5 * time.Second        //默认过期
	DefaultExpireUnsame = 600 * time.Second      //默认过期时间
)

var ErrFailed = errors.New("failed to acquire lock")

type RedisLock struct {
	driver        Driver
	Expiry        time.Duration
	Retry         int
	RetryInterval time.Duration
	lockKey       string
	lockValue     string
	sameLock      bool
	mutex         sync.Mutex
}

func NewRedisLock(lk string, sameLock bool) *RedisLock {
	if sameLock {
		return &RedisLock{
			lockKey:       lk,
			Expiry:        DefaultExpire,
			Retry:         DefaultRetry,
			RetryInterval: DefaultInterval,
			driver:        Get(),
		}
	} else {
		return &RedisLock{
			lockKey:       lk,
			Expiry:        DefaultExpireUnsame,
			Retry:         1,
			RetryInterval: 0,
			driver:        Get(),
		}
	}
}

func NewRedisLockWithParam(lk string, expire int, retry int, interval int) *RedisLock {
	return &RedisLock{
		lockKey:       lk,
		Expiry:        time.Duration(expire) * time.Second,
		Retry:         retry,
		RetryInterval: time.Duration(interval) * time.Millisecond,
		driver:        Get(),
	}
}

func (rl *RedisLock) Lock() error {
	if rl.sameLock {
		rl.mutex.Lock()
		err := rl.lock()
		if err != nil {
			rl.mutex.Unlock()
		}
		return err
	} else {
		err := rl.lock()
		return err
	}
}

func (rl *RedisLock) lock() error {
	if rl.sameLock {
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			return err
		}
		rl.lockValue = base64.StdEncoding.EncodeToString(b)
	} else {
		rl.lockValue = ""
	}
	for i := 0; i < rl.Retry; i++ {
		ok, err := rl.driver.SetNX(rl.lockKey, rl.lockValue, rl.Expiry)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
		if i < rl.Retry-1 {
			time.Sleep(rl.RetryInterval)
		}
	}
	return ErrFailed
}

const delScript = `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end`

func (rl *RedisLock) UnLock() error {
	if rl.sameLock {
		_, err := rl.driver.Eval(delScript, []string{rl.lockKey}, rl.lockValue)
		rl.mutex.Unlock()
		return err
	} else {
		_, err := rl.driver.Del(rl.lockKey)
		return err
	}
}

// RedLock 简单锁 name锁名称 expire 有效期 retry 重试次数
func RedLock(name string, expire int, retry int) (lock *RedisLock, res bool) {
	lock = NewRedisLockWithParam("redLock:"+name, expire, retry, 300)
	if err := lock.Lock(); err != nil {
		return nil, false
	}
	return lock, true
}
