package timer

import (
	"dnspod_ddns_go/config"
	"dnspod_ddns_go/dnspod"
	"time"
)

func Start(conf *config.Config) {
	dnspod.SetUp(conf)
	dnspod.RecordList()
	timer := time.NewTimer(conf.Timer)
	defer timer.Stop()
	for {
		timer.Reset(conf.Timer) // 这里复用了 timer
		select {
		case <-timer.C:
			dnspod.RecordList()
		}
	}
}
