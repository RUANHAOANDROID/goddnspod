package timer

import (
	"context"
	"dnspod_ddns_go/config"
	"dnspod_ddns_go/dnspod"
	"sync"
	"time"
)

type Timer struct {
	interval   time.Duration
	ctx        context.Context
	cancel     context.CancelFunc
	timerMutex sync.Mutex
	ticker     *time.Ticker
}

// NewTimer 创建一个新的 Timer 实例
func NewTimer() *Timer {
	ctx, cancel := context.WithCancel(context.Background())
	return &Timer{ctx: ctx, cancel: cancel}
}

// Start 启动定时器
func (t *Timer) Start() {
	t.timerMutex.Lock()
	defer t.timerMutex.Unlock()

	// 确保之前的定时器已经停止
	if t.ticker != nil {
		t.ticker.Stop()
	}

	conf, err := config.Load("config.yml")
	if err != nil {
		config.CreateEmpty().Save()
		panic("请完善配置config.yml")
	}
	interval, err := time.ParseDuration(conf.Timer)
	if err != nil {
		panic(err)
	}

	t.interval = interval
	dnspod.SetUp(conf)
	dnspod.RecordList()
	t.ticker = time.NewTicker(t.interval)

	go func() {
		for {
			select {
			case <-t.ticker.C:
				dnspod.RecordList()
			case <-t.ctx.Done():
				return
			}
		}
	}()
}

// Stop 停止定时器
func (t *Timer) Stop() {
	t.timerMutex.Lock()
	defer t.timerMutex.Unlock()
	if t.ticker != nil {
		t.ticker.Stop()
		t.ticker = nil
	}
	t.cancel()
}

// Restart 重新启动定时器
func (t *Timer) Restart() {
	t.Stop()  // 停止当前的定时器
	t.Start() // 启动新的定时器
}
