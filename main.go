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
	fmt.Println("-Config INFO")
	fmt.Printf("--UserAgent:%s\n", conf.UserAgent)
	fmt.Printf("--TokenId:%s\n", conf.TokenId)
	fmt.Printf("--LoginToken:%s\n", conf.LoginToken)
	fmt.Printf("--Domain:%s\n", conf.Domain)
	fmt.Printf("--SubDomain:%s\n", conf.SubDomain)
	fmt.Printf("--Timer:%s\n", conf.Timer)
	fmt.Printf("--Support:%s\n", conf.Support)
	fmt.Println("-Start Success")
	timer.Start(conf)
}
