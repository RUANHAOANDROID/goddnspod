package main

import (
	"dnspod_ddns_go/api"
	"dnspod_ddns_go/config"
	"dnspod_ddns_go/dnspod"
	"dnspod_ddns_go/timer"
	"fmt"
)

func main() {
	fmt.Println("<(￣︶￣)↗[GO!]...")
	conf, err := config.Load("config.yml")
	if err != nil {
		config.CreateEmpty().Save()
		panic("请完善配置config.yml")
	}
	fmt.Printf("--config:%v\n", conf)
	dnspod.SetUp(conf)
	dnsTimer := timer.NewTimer()
	defer dnsTimer.Stop()
	dnsTimer.Start()
	api.Register(dnsTimer)
}
