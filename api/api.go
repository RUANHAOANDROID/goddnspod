package api

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

func Register(dnsTimer *timer.Timer) {
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
	conf, err := config.Load("config.yml")
	if err != nil {
		config.CreateEmpty().Save()
		panic("请完善配置config.yml")
	}
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
			fmt.Println(newConfig)
			newConfig.Save()
			dnsTimer.Restart()
			http.StatusText(http.StatusOK)
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		}
	}))
	// 启动服务器
	port := ":5173"
	fmt.Println("DDNSPod server on " + port)
	err = http.ListenAndServe(port, handler)
	if err != nil {
		log.Fatal(err)
	}
}
