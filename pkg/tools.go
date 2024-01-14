package pkg

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
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
	ip, err := ExtractIP(string(body))
	if err != nil {
		fmt.Println("extract ip address error:", err)
		return "", err
	}
	return ip, nil
}
