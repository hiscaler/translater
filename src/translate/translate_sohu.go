package translate

import (
	"fmt"
	"log"
	"strings"
	"net/http"
	"math/rand"
	"crypto/md5"
	"io/ioutil"
	"encoding/hex"
	"encoding/json"
	"errors"
)

// 搜狐翻译
type SohuTranslate struct {
	Translate
}

type SohuResponse struct {
	Zly         string
	Query       string
	Translation string
	ErrorCode   string
}

var errorCodes = map[string]string{
	"1001":  "不支持的语言类型",
	"1002":  "文本过长",
	"1003":  "无效PID",
	"1004":  "试用Pid限额已满",
	"1005":  "Pid请求流量过高",
	"1006":  "余额不足",
	"1007":  "随机数不存在",
	"1008":  "签名不存在",
	"1009":  "签名不正确",
	"10010": "文本不存在",
	"1050":  "内部服务错误",
}

// 翻译处理
func (t *SohuTranslate) Do() (*SohuTranslate, error) {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	randSeq := func(n int) string {
		b := make([]rune, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return string(b)
	}
	client := &http.Client{}
	for i, node := range t.Nodes {
		t.Anatomy(node)
		s := strings.Trim(t.currentNodeText, "")
		if t.Debug {
			log.Println(fmt.Sprintf("#%v: %#v", i+1, s))
		}
		if len(s) > 0 {
			salt := randSeq(12)
			account := t.GetAccount()
			mdx := md5.New()
			mdx.Write([]byte(account.PID + s + salt + account.SecretKey))
			sign := hex.EncodeToString(mdx.Sum(nil))
			fields := map[string]string{
				"q":    s,
				"from": t.From,
				"to":   t.To,
				"pid":  account.PID,
				"salt": salt,
				"sign": sign,
			}
			payload := make([]string, len(fields))
			for k, v := range fields {
				payload = append(payload, fmt.Sprintf("%s=%s", k, v))
			}

			req, err := http.NewRequest("POST", "http://fanyi.sogou.com/reventondc/api/sogouTranslate", strings.NewReader(strings.Join(payload, "&")))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Accept", "application/json")
			if err == nil {
				resp, err := client.Do(req)
				if err == nil {
					body, err := ioutil.ReadAll(resp.Body)
					if err == nil {
						sohuResponse := &SohuResponse{}
						err = json.Unmarshal([]byte(body), &sohuResponse)
						if err == nil {
							if sohuResponse.ErrorCode == "0" {
								s = sohuResponse.Translation
							} else {
								msg, exists := errorCodes[sohuResponse.ErrorCode]
								if !exists {
									msg = sohuResponse.ErrorCode
								}
								return t, errors.New(msg)
							}
						}
					}
					resp.Body.Close()
				}
				t.currentNode.Data = s
			} else {
				return t, err
			}
		}
	}

	return t, nil
}
