package main

import (
	"dnspod_ddns_go/config"
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
	fmt.Println("Timer interval", conf.Timer)
	timer.Start(conf)
}
