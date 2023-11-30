package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type RedisLock struct {
	key    string
	token  string
	client *redis.Client
}

// GetCurrentProcessID 获取当前进程id
func GetCurrentProcessID() string {
	return strconv.Itoa(os.Getpid())
}

// GetCurrentGoroutineID 获取当前的协程ID
func GetCurrentGoroutineID() string {
	buf := make([]byte, 128)
	buf = buf[:runtime.Stack(buf, false)]
	stackInfo := string(buf)
	return strings.TrimSpace(strings.Split(strings.Split(stackInfo, "[running]")[0], "goroutine")[1])
}

func GetLockID() string {
	return GetCurrentProcessID() + GetCurrentGoroutineID()
}

func NewRedisLock(key string, client *redis.Client) *RedisLock {
	return &RedisLock{
		key:    key,
		token:  GetLockID(),
		client: client,
	}
}

func (r *RedisLock) Lock(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			return
		}
	}()

	err = r.tryLock()
	if err == nil {
		return nil
	}

	err = r.blockingLock(ctx)
	return
}

func (r *RedisLock) tryLock() error {
	reply, err := r.client.SetNX("RedisLock:"+r.key, r.token, time.Second*10).Result()
	if err != nil {
		return err
	}
	if !reply {
		return fmt.Errorf("lock fail")
	}
	return nil
}

func (r *RedisLock) blockingLock(ctx context.Context) error {
	// 阻塞模式等锁时间上限-2s
	timeoutCh := time.After(time.Second * 2)
	// 轮询 ticker，每隔 50 ms 尝试取锁一次
	ticker := time.NewTicker(time.Duration(50) * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		select {
		// ctx 终止了
		case <-ctx.Done():
			return fmt.Errorf("lock failed, ctx timeout, err: %w", ctx.Err())
			// 阻塞等锁达到上限时间
		case <-timeoutCh:
			return fmt.Errorf("block waiting time out")
		// 放行
		default:
		}

		// 尝试取锁
		err := r.tryLock()
		if err == nil {
			// 加锁成功，返回结果
			return nil
		}
	}

	// 不可达
	return nil
}

const LuaCheckAndDeleteDistributionLock = `
  local lockerKey = KEYS[1]
  local targetToken = ARGV[1]
  local getToken = redis.call('get',lockerKey)
  if (not getToken or getToken ~= targetToken) then
    return 0
  else
    return redis.call('del',lockerKey)
  end
`

func (r *RedisLock) Unlock() (err error) {
	defer func() {
		if err != nil {
			return
		}
	}()
	keys := []string{"RedisLock:" + r.key}
	reply, _err := r.client.Eval(LuaCheckAndDeleteDistributionLock, keys, r.token).Result()
	if _err != nil {
		err = _err
		return
	}
	if ret, _ := reply.(int64); ret != 1 {
		err = errors.New("can not unlock without ownership of lock")
	}

	return nil
}
