package utils

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestRedisLock(t *testing.T) {
	err := InitClient()
	if err != nil {
		panic("connect error")
	}
	lock1 := NewRedisLock("test_key", Client)
	lock2 := NewRedisLock("test_key", Client)
	ctx := context.Background()
	err = lock1.Lock(ctx)
	if err != nil {
		fmt.Print("lock error")
	}
	err = lock2.Lock(ctx)
	if err != nil {
		fmt.Print("lock2 lock fail\n")
	}

	time.Sleep(time.Second * 3)
	err = lock1.Unlock()
	if err != nil {
		fmt.Print("unlock error")
	}
	fmt.Print("test success")
	return
}
