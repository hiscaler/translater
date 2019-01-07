package translate

import (
	"fmt"
	"strings"
	"net/http"
	"math/rand"
	"crypto/md5"
	"io/ioutil"
	"errors"
	"net/url"
	"runtime"
	"log"
)

// 谷歌翻译
type GoogleTranslate struct {
	Translate
}

type GoogleResponse struct {
	Zly         string
	Query       string
	Translation string
	ErrorCode   string
}

var googleErrorCodes = map[string]string{
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

func (t *GoogleTranslate) Req(i int, s string, in chan<- string) (string, error) {
	if len(s) > 0 {
		letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		randSeq := func(n int) string {
			b := make([]rune, n)
			for i := range b {
				b[i] = letters[rand.Intn(len(letters))]
			}
			return string(b)
		}
		salt := randSeq(12)
		account := t.GetRandomAccount()
		mdx := md5.New()
		mdx.Write([]byte(account.PID + s + salt + account.SecretKey))

		var rq = url.Values{}
		rq.Add("q", s)
		resp, err := http.Get("https://translate.google.cn/translate_a/single?client=gtx&sl=en&tl=zh&dt=t&" + rq.Encode())
		fmt.Println("Code:", resp.StatusCode)
		if err == nil && resp.StatusCode == 200 {
			body, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				googleResponse := &GoogleResponse{}
				afterString := string(body)
				if index := strings.Index(string(afterString), "]],"); index != -1 {
					afterString = strings.Replace(afterString[0:index+1], "[[[\"", "", -1)
					if index = strings.Index(afterString, "\","); index != -1 {
						afterString = strings.TrimSpace(afterString[:index])
						googleResponse.Translation = afterString
					}
				}

				in <- googleResponse.Translation
				return googleResponse.Translation, nil
			} else {
				log.Println("Google error, ", err)
			}
			resp.Body.Close()
		} else {
			log.Println("Google error, ", err)
		}
		in <- s
		return s, err
	}

	return s, errors.New("Text is empty.")

}

// 翻译处理
func (t *GoogleTranslate) Do() (*GoogleTranslate, error) {
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
