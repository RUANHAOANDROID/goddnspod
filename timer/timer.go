package timer

import (
	"context"
	"dnspod_ddns_go/config"
	"dnspod_ddns_go/dnspod"
	"time"
)

type Timer struct {
	conf   *config.Config
	ctx    context.Context
	cancel context.CancelFunc
}

func NewTimer(conf *config.Config) *Timer {
	ctx, cancel := context.WithCancel(context.Background())
	return &Timer{conf: conf, ctx: ctx, cancel: cancel}
}

func (t *Timer) Start() {
	dnspod.SetUp(t.conf)
	dnspod.RecordList()
	timer, err := time.ParseDuration(t.conf.Timer)
	if err != nil {
		panic(err)
	}
	ticker := time.NewTicker(timer)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			dnspod.RecordList()
		case <-t.ctx.Done():
			return
		}
	}
}
func (t *Timer) Stop() {
	t.cancel()
}
