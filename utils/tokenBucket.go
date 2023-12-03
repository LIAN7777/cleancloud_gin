package utils

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity     int        //桶容量
	fillRate     int        //令牌填充速度
	currentToken int        //当前令牌剩余数
	mu           sync.Mutex //锁，保证操作时只有一个goroutine能够进入bucket用于保护令牌数互斥访问
}

func NewTokenBucket(cap int, rate int, unit time.Duration) *TokenBucket {
	tb := &TokenBucket{
		capacity:     cap,
		fillRate:     rate,
		currentToken: cap,
		mu:           sync.Mutex{},
	}
	go tb.fill(unit)
	return tb
}

func (tb *TokenBucket) fill(unit time.Duration) {
	ticker := time.NewTicker(unit)
	for range ticker.C {
		tb.mu.Lock()
		tb.currentToken += tb.fillRate
		if tb.currentToken > tb.capacity {
			tb.currentToken = tb.capacity
		}
		tb.mu.Unlock()
	}
}

// TryAcquire 获取token失败直接返回失败
func (tb *TokenBucket) TryAcquire() bool {
	tb.mu.Lock()
	if tb.currentToken > 0 {
		tb.currentToken--
		//fmt.Print(tb.currentToken, "\n")
		tb.mu.Unlock()
		return true
	}
	tb.mu.Unlock()
	return false
}

// AcquireWithTimeOut 设置超时等待时间的获取token方法
func (tb *TokenBucket) AcquireWithTimeOut(timeOut time.Duration) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	if tb.currentToken > 0 {
		tb.currentToken--
		return true
	}
	timer := time.NewTimer(timeOut)
	for {
		tb.mu.Unlock()
		select {
		case <-timer.C: //超时
			return false
		default:
			tb.mu.Lock()
			if tb.currentToken > 0 { // 获取到令牌，返回true
				tb.currentToken--
				return true
			}
		}
	}
}

// Acquire 无限等待的获取token方法
func (tb *TokenBucket) Acquire() {
	tb.mu.Lock()
	for tb.currentToken <= 0 {
		tb.mu.Unlock()
		time.Sleep(time.Millisecond * 50)
		tb.mu.Lock()
	}
	tb.currentToken--
	tb.mu.Unlock()
}
