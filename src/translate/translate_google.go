package translate

import (
	"fmt"
	"strings"
	"net/http"
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

func (t *GoogleTranslate) Req(i int, s string, in chan<- string) (string, error) {
	if len(s) > 0 {
		rq := url.Values{}
		rq.Add("q", s)
		resp, err := http.Get("https://translate.google.cn/translate_a/single?client=gtx&sl=en&tl=zh&dt=t&" + rq.Encode())
		if err == nil && resp.StatusCode == 200 {
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
				log.Println(err)
			}
			resp.Body.Close()
		} else {
			log.Println(err)
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
