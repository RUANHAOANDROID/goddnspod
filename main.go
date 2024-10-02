package main

import (
	"dnspod_ddns_go/config"
	"dnspod_ddns_go/timer"
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"io"
	"log"
	"net/http"
)

func main() {
	fmt.Println("<(￣︶￣)↗[GO!]...")
	conf, err := config.Load("config.yml")
	if err != nil {
		config.CreateEmpty().Save()
		panic("请完善配置config.yml")
	}
	fmt.Println("----Config info")
	fmt.Printf("--UserAgent:%s\n", conf.UserAgent)
	fmt.Printf("--TokenId:%s\n", conf.TokenId)
	fmt.Printf("--LoginToken:%s\n", conf.LoginToken)
	fmt.Printf("--Domain:%s\n", conf.Domain)
	fmt.Printf("--SubDomain:%s\n", conf.SubDomain)
	fmt.Printf("--Timer:%s\n", conf.Timer)
	fmt.Printf("--Support:%s\n", conf.Support)
	fmt.Println("----Setup success")

	dnsTimer := timer.NewTimer(conf)
	defer dnsTimer.Stop()
	go dnsTimer.Start()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // 允许的源
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	// 创建您的http.Handler
	mux := http.NewServeMux()
	// 使用CORS中间件包装您的handler
	handler := c.Handler(mux)

	// 启动服务器
	mux.Handle("/", http.FileServer(http.Dir("./frontend/build/client")))
	mux.Handle("/config", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			http.StatusText(http.StatusOK)
			w.WriteHeader(http.StatusOK)
			resp, err := json.Marshal(conf)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(resp)
		case http.MethodPost:
			var newConfig *config.Config
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = json.Unmarshal(body, &newConfig)
			if err != nil {
				fmt.Println(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			fmt.Println("----New Config info")
			fmt.Printf("--UserAgent:%v\n", newConfig)
			newConfig.Save()
			http.StatusText(http.StatusOK)
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		}
	}))
	// 启动服务器
	fmt.Println("DDNSPod server on :6565")
	err = http.ListenAndServe(":6565", handler)
	if err != nil {
		log.Fatal(err)
	}
}
