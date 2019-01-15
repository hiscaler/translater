package translate

import (
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
	"errors"
	"runtime"
)

// 谷歌翻译
type GoogleTranslate struct {
	Translate
}

func (t *GoogleTranslate) Req(i int, s string, in chan<- string) (string, error) {
	if len(s) > 0 {
		req, err := http.NewRequest("GET", "https://translate.google.cn/translate_a/single?client=gtx&dt=t", nil)
		if err == nil {
			req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
			req.Header.Add("referer", "https://translate.google.com.ph/?hl=zh-CN")
			q := req.URL.Query()
			q.Add("sl", t.From)
			q.Add("tl", t.To)
			q.Add("q", s)
			req.URL.RawQuery = q.Encode()
			client := &http.Client{}
			resp, err := client.Do(req)
			defer resp.Body.Close()
			if err == nil {
				if resp.StatusCode == 200 {
					body, err := ioutil.ReadAll(resp.Body)
					if err == nil {
						translateAfterText := string(body)
						if index := strings.Index(string(translateAfterText), "]],"); index != -1 {
							translateAfterText = strings.Replace(translateAfterText[0:index+1], "[[[\"", "", -1)
							if index = strings.Index(translateAfterText, "\","); index != -1 {
								translateAfterText = strings.TrimSpace(translateAfterText[:index])
							}
						}

						in <- translateAfterText
						return translateAfterText, nil
					} else {
						t.Logger.ErrorLogger.Println(err)
					}
				} else {
					t.Logger.ErrorLogger.Println("HTTP CODE: " + string(resp.StatusCode))
				}
			} else {
				t.Logger.ErrorLogger.Println(err)
			}
		} else {
			t.Logger.ErrorLogger.Println(err)
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
