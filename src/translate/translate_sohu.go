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
)

// 搜狐翻译
type SohuTranslate struct {
	Translate
	PID       string
	SecretKey string
}

// 翻译处理
func (t *SohuTranslate) Do() (*SohuTranslate, error) {
	for i, node := range t.Nodes {
		s := strings.TrimSpace(t.Anatomy(node))
		if t.Debug {
			log.Println(fmt.Sprintf("#%v: %#v", i+1, s))
		}

		s = t.currentNodeText
		if len(strings.TrimSpace(s)) > 0 {
			pid := t.PID
			letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
			randSeq := func(n int) string {
				b := make([]rune, n)
				for i := range b {
					b[i] = letters[rand.Intn(len(letters))]
				}
				return string(b)
			}
			salt := randSeq(12)
			key := t.SecretKey
			mdx := md5.New()
			mdx.Write([]byte(pid + s + salt + key))
			sign := hex.EncodeToString(mdx.Sum(nil))
			fields := map[string]string{
				"q":    s,
				"from": "zh-CHS",
				"to":   "en",
				"pid":  pid,
				"salt": salt,
				"sign": string(sign),
			}
			payload := make([]string, len(fields))
			for k, v := range fields {
				payload = append(payload, fmt.Sprintf("%s=%s", k, v))
			}
			client := &http.Client{}
			req, err := http.NewRequest("POST", "http://fanyi.sogou.com/reventondc/api/sogouTranslate", strings.NewReader(strings.Join(payload, "&")))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Accept", "application/json")
			if err == nil {
				resp, err := client.Do(req)
				if err == nil {
					body, err := ioutil.ReadAll(resp.Body)
					if err == nil {
						s = string(body)
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
