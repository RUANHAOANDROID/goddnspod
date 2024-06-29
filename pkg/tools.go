package pkg

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
)

func ExtractIP(input string) (string, error) {
	ipRegex := regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)
	match := ipRegex.FindString(input)
	if match == "" {
		return "", fmt.Errorf("not find ip address")
	}
	return match, nil
}

func PublicIP() (string, error) {
	url := "https://myip.ipip.net/"
	resp, err := http.Get(url)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("反差IP发生错误:", err)
		return "", err
	}
	ip, err := ExtractIP(string(body))
	if err != nil {
		fmt.Println("extract ip address error:", err)
		return "", err
	}
	return ip, nil
}
func PublicIPV6() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	var ipv6 string
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() == nil {
			v6 := ipNet.IP
			if v6.IsLoopback() || v6.IsLinkLocalUnicast() || v6.IsLinkLocalMulticast() || v6.IsMulticast() || v6.IsUnspecified() {
				continue
			}
			// 检查是否是唯一本地地址 (Unique local addresses)
			if strings.HasPrefix(v6.String(), "fc") || strings.HasPrefix(v6.String(), "fd") {
				continue
			}
			if strings.HasPrefix(v6.String(), "2001:") {
				continue
			}
			if strings.HasPrefix(v6.String(), "2409:") {
				ipv6 = v6.String()
				break
			}
		}
	}

	return ipv6, nil
}

// isPublicIPv6 判断一个IPv6地址是否是公网地址
func isPublicIPv6(ipv6 string) bool {
	if strings.HasPrefix(ipv6, "f") {
		return false
	}
	return true
}
