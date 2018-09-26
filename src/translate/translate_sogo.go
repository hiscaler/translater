package translate

import (
	"fmt"
	"strings"
	"net/http"
	"math/rand"
	"crypto/md5"
	"io/ioutil"
	"encoding/json"
	"errors"
	"encoding/hex"
	"net/url"
	"runtime"
)

// 搜狗翻译
type SogoTranslate struct {
	Translate
}

type SogoResponse struct {
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

func (t *SogoTranslate) Req(i int, s string, in chan<- string) (string, error) {
	if len(s) > 0 {
		letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		randSeq := func(n int) string {
			b := make([]rune, n)
			for i := range b {
				b[i] = letters[rand.Intn(len(letters))]
			}
			return string(b)
		}
		client := &http.Client{}
		salt := randSeq(12)
		account := t.GetRandomAccount()
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
		body := url.Values{}
		for k, v := range fields {
			body.Set(k, v)
		}
		req, err := http.NewRequest("POST", "http://fanyi.sogou.com/reventondc/api/sogouTranslate", strings.NewReader(body.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept", "application/json")
		if err == nil {
			resp, err := client.Do(req)
			if err == nil && resp.StatusCode == 200 {
				body, err := ioutil.ReadAll(resp.Body)
				if err == nil {
					sogoResponse := &SogoResponse{}
					err = json.Unmarshal([]byte(body), &sogoResponse)
					if err == nil {
						if sogoResponse.ErrorCode == "0" {
							if t.Config.Debug {
								t.Logger.InfoLogger.Println(fmt.Sprintf("#%v After: %#v", i+1, sogoResponse.Translation))
							}
							in <- sogoResponse.Translation
							return sogoResponse.Translation, nil
						} else {
							if sogoResponse.ErrorCode == "1003" || sogoResponse.ErrorCode == "1004" || sogoResponse.ErrorCode == "1005" {
								_, err = t.updateAccount(account.PID, false)
								if err != nil {
									t.Logger.ErrorLogger.Println(err)
								}
							}
							msg, exists := errorCodes[sogoResponse.ErrorCode]
							if !exists {
								msg = sogoResponse.ErrorCode
							}
							if t.Config.Debug {
								msg = fmt.Sprintf("%v (%v)", msg, s)
							}

							err = errors.New(msg)
							if t.Config.Debug {
								t.Logger.ErrorLogger.Println(fmt.Sprintf("#%v After: %#v", i+1, err.Error()))
							}
						}
					} else {
						if t.Config.Debug {
							t.Logger.ErrorLogger.Println(fmt.Sprintf("#%v After: %#v", i+1, string(body)))
						}
					}
				}
				resp.Body.Close()
			}
		}
		in <- s
		return s, err
	}

	return s, errors.New("Text is empty.")

}

// 翻译处理
func (t *SogoTranslate) Do() (*SogoTranslate, error) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	out := make(chan string, len(t.Nodes))
	for i, node := range t.Nodes {
		t.Anatomy(node)
		s := strings.Trim(t.currentNodeText, " \r\n\t")
		if t.Config.Debug {
			t.Logger.InfoLogger.Println(fmt.Sprintf("#%v Before: %#v", i+1, s))
		}

		go t.Req(i, s, out)
		t.currentNode.Data = <-out
	}
	close(out)

	return t, nil
}

// 翻译处理
// Deprecated: 弃用，改为使用 goroutine 方式
func (t *SogoTranslate) _Do() (*SogoTranslate, error) {
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
		s := strings.Trim(t.currentNodeText, " \r\n\t")
		if t.Config.Debug {
			t.Logger.InfoLogger.Println(fmt.Sprintf("#%v Before: %#v", i+1, s))
		}
		if len(s) > 0 {
			salt := randSeq(12)
			account := t.GetRandomAccount()
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
			body := url.Values{}
			for k, v := range fields {
				body.Set(k, v)
			}
			req, err := http.NewRequest("POST", "http://fanyi.sogou.com/reventondc/api/sogouTranslate", strings.NewReader(body.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Accept", "application/json")
			if err == nil {
				resp, err := client.Do(req)
				if err == nil && resp.StatusCode == 200 {
					body, err := ioutil.ReadAll(resp.Body)
					if err == nil {
						sogoResponse := &SogoResponse{}
						err = json.Unmarshal([]byte(body), &sogoResponse)
						if err == nil {
							if sogoResponse.ErrorCode == "0" {
								t.currentNode.Data = sogoResponse.Translation
								if t.Config.Debug {
									t.Logger.InfoLogger.Println(fmt.Sprintf("#%v After: %#v", i+1, sogoResponse.Translation))
								}
							} else {
								if sogoResponse.ErrorCode == "1003" || sogoResponse.ErrorCode == "1004" || sogoResponse.ErrorCode == "1005" {
									_, err = t.updateAccount(account.PID, false)
									if err != nil {
										t.Logger.ErrorLogger.Println(err)
									}
								}
								msg, exists := errorCodes[sogoResponse.ErrorCode]
								if !exists {
									msg = sogoResponse.ErrorCode
								}
								if t.Config.Debug {
									msg = fmt.Sprintf("%v (%v)", msg, s)
								}

								err = errors.New(msg)
								if t.Config.Debug {
									t.Logger.ErrorLogger.Println(fmt.Sprintf("#%v After: %#v", i+1, err.Error()))
								}

								return t, err
							}
						} else {
							if t.Config.Debug {
								t.Logger.ErrorLogger.Println(fmt.Sprintf("#%v After: %#v", i+1, string(body)))
							}
						}
					}
					resp.Body.Close()
				}
			} else {
				return t, err
			}
		}
	}

	return t, nil
}
