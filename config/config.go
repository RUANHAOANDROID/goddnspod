package config

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"sync"
)

// Config 系统整体配置--
type Config struct {
	UserAgent  string `json:"userAgent" yaml:"userAgent"`
	TokenId    string `json:"tokenId" yaml:"tokenId"`
	LoginToken string `json:"loginToken" yaml:"loginToken"`
	Domain     string `json:"domain" yaml:"domain"`
	SubDomain  string `json:"subDomain" yaml:"subDomain"`
	Timer      string `json:"timer" yaml:"timer"`
	Support    string `json:"support" yaml:"Support"`
}

var path string
var configMux sync.RWMutex // 用于确保并发读写安全

// Load 加载配置
func Load(path string) (*Config, error) {
	configMux.RLock() // 读锁
	defer configMux.RUnlock()
	// 使用 viper 读取配置文件
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	c := &Config{}
	// 将配置绑定到结构体
	err = viper.Unmarshal(c)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return c, nil
}
func CreateEmpty() *Config {
	return &Config{
		UserAgent:  "Hao88 DDNS/1.0Alpha(52927295@qq.com) ",
		TokenId:    "",
		LoginToken: "",
		Domain:     "youdomain.com",
		SubDomain:  "sub.youdomain.com",
		Timer:      "5m30s",
	}
}
func (c *Config) Save() {
	configMux.Lock() // 写锁
	defer configMux.Unlock()
	// 将结构体转换为字节数组
	yamlBytes, err := yaml.Marshal(c)
	if err != nil {
		fmt.Println("无法将结构体转换为YAML格式：", err)
		return
	}

	// 将字节数组加载到Viper
	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(yamlBytes))
	if err != nil {
		fmt.Println("无法加载配置：", err)
		return
	}

	// 写入配置到文件
	filePath := "config.yml"
	err = viper.WriteConfigAs(filePath)
	if err != nil {
		fmt.Println("无法写入配置文件：", err)
		return
	}

	fmt.Println("配置已成功写入YAML文件：", filePath)
}
