package dnspod

import (
	"bytes"
	"dnspod_ddns_go/config"
	"dnspod_ddns_go/pkg"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const dnsPodUrl = "https://dnsapi.cn"
const recordList = "/Record.List"
const recordModify = "/Record.Modify"
const infoVersion = "/Info.Version"

var conf *config.Config

func SetUp(config *config.Config) {
	conf = config
}
func crtBaseParams() url.Values {
	params := url.Values{}
	params.Set("login_token", conf.TokenId+","+conf.LoginToken) //必须
	params.Set("format", conf.Format)                           //json
	return params
}
func addHeaders(header *http.Header) {
	header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	header.Set("User-Agent", conf.UserAgent)
}
func InfoVersion() {
	params := crtBaseParams().Encode()
	req, err := http.NewRequest("POST", dnsPodUrl+infoVersion, bytes.NewBuffer([]byte(params)))
	if err != nil {
		fmt.Println(err)
	}
	addHeaders(&req.Header)
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}
func RecordList() {
	if conf.LoginToken == "" || conf.TokenId == "" {
		panic("tokenId或loginToken未填写!")
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("发生错误", r)
		}
	}()
	params := crtBaseParams()
	params.Set("domain", conf.Domain)
	req, err := creatRequest(recordList, params)
	if err != nil {
		fmt.Println("Error:", err)
	}
	addHeaders(&req.Header)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	jsonData := string(body)
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		fmt.Println(jsonData)
		fmt.Println("JSON 解析失败:", err)
		return
	}
	// 获取 records 数组
	records, ok := data["records"].([]interface{})
	if !ok {
		fmt.Println("找不到 records 数组")
		return
	}
	// 遍历 records 数组，获取每个记录的 name 字段的值
	for _, record := range records {
		recordMap, ok := record.(map[string]interface{})
		if !ok {
			fmt.Println("记录格式错误")
			continue
		}

		name, ok := recordMap["name"].(string)
		if !ok {
			fmt.Println("找不到 name 字段")
			continue
		}
		if name == conf.SubDomain {
			value, ok := recordMap["value"].(string)
			if !ok {
				fmt.Println("not find value")
				break
			}
			currPIP, err := pkg.PublicIP()
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("Remote IP:", value, "Current IP:", currPIP, "time:", time.Now().Local().Format("2006-01-02 15:04:05"))
			if value != currPIP {
				recordId := recordMap["id"].(string)
				recordType := recordMap["type"].(string)
				recordMx := recordMap["mx"].(string)
				RecordModify(currPIP, recordId, recordType, recordMx)
			}
		}
	}
}

func RecordModify(publicIp string, recordId string, recordType string, mx string) {
	params := crtBaseParams()
	params.Set("record_id", recordId) //记录ID。
	params.Set("record_line", "默认")
	params.Set("domain", conf.Domain)        //域名
	params.Set("sub_domain", conf.SubDomain) //子域名
	params.Set("value", publicIp)            //ipv4 or v6
	params.Set("record_type", recordType)    //A
	params.Set("mx", mx)                     //0
	req, err := creatRequest(recordModify, params)
	if err != nil {
		fmt.Println("Create request error:", err)
		return
	}
	addHeaders(&req.Header)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("HTTP POST request error:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read response error:", err)
		return
	}

	// 打印响应内容
	fmt.Println("response body:", string(body))
}

func creatRequest(path string, params url.Values) (*http.Request, error) {
	return http.NewRequest("POST", dnsPodUrl+path, bytes.NewBuffer([]byte(params.Encode())))
}
