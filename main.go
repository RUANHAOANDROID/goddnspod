package main

import (
	"dnspod_ddns_go/config"
	"dnspod_ddns_go/timer"
	"fmt"
)

func main() {
	conf, err := config.Load("config.yml")
	if err != nil {
		panic("Not find config.yml")
	}
	fmt.Println(conf)
	timer.Start(conf)
}
