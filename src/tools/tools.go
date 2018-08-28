package main

import (
	"io/ioutil"
	"fmt"
	"strings"
	"encoding/json"
)

// 帐号 json 格式文件生成
func main() {
	b, err := ioutil.ReadFile("./src/config/accounts.txt")
	if err != nil {
		fmt.Println("Read file error")
	}

	type Account struct {
		PID       string
		SecretKey string
		Enabled   bool
		YearMonth int
	}
	accounts := make([]Account, 0)
	s := string(b)
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		t := strings.Split(line, "\t") // Format: sn, PID, SecretKey, ...
		if len(t) >= 3 {
			account := &Account{
				PID:       t[1],
				SecretKey: t[2],
				Enabled:   true,
				YearMonth: 0,
			}
			accounts = append(accounts, *account)
		}
	}
	if len(accounts) != 0 {
		jsonString, err := json.MarshalIndent(accounts, "", "    ")
		if err == nil {
			ioutil.WriteFile("./src/config/accounts.js", jsonString, 0644)
			fmt.Println("帐号文件写入成功。")
		} else {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("没有可写入的帐号。")
	}
}
