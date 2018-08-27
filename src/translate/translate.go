package translate

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"bytes"
	"io"
	"fmt"
	"log"
	"strings"
	"math/rand"
	"time"
	"github.com/spf13/viper"
	"strconv"
)

type IT interface {
	Parse() (string, error)
	Do()
	SetRawContent(s string)
	GetAccount()
}

type Translate struct {
	IT
	Config          *Config
	Viper           *viper.Viper // Config
	From            string       // 来源语言
	To              string       // 目标语言
	Doc             *html.Node
	Nodes           []*html.Node
	rawContent      string
	ignoreAtoms     []atom.Atom
	currentNode     *html.Node
	currentNodeText string
	Languages       map[string]string // 支持的语种
	Accounts        []Account         // 账号列表
}

// 设置要翻译的文本内容
func (t *Translate) SetRawContent(s string) *Translate {
	t.rawContent = s
	t.ignoreAtoms = []atom.Atom{atom.Ping, atom.Script, atom.Noscript, atom.Style}

	return t
}

func (t *Translate) GetRawContent() string {
	return t.rawContent
}

// 解析并翻译需要处理的内容
func (t *Translate) Parse() (*Translate, error) {
	doc, err := html.Parse(strings.NewReader(t.rawContent))
	if err != nil {
		return t, err
	}

	t.Doc = doc

	matcher := func(node *html.Node) (keep bool, exit bool) {
		if node.Type == html.TextNode && strings.TrimSpace(node.Data) != "" {
			keep = true
		}

		for _, v := range t.ignoreAtoms {
			if node.DataAtom == v {
				exit = true
			}
		}

		return
	}

	t.Nodes = t.traverseNode(doc, matcher)

	return t, nil
}

func (t *Translate) traverseNode(doc *html.Node, matcher func(node *html.Node) (bool, bool)) (nodes []*html.Node) {
	var keep, exit bool
	var f func(*html.Node)
	f = func(n *html.Node) {
		keep, exit = matcher(n)
		if keep {
			nodes = append(nodes, n)
		}
		if exit {
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return nodes
}

// 处理单个节点
func (t *Translate) Anatomy(n *html.Node) string {
	t.currentNode = n
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	t.currentNodeText = buf.String()

	return t.currentNodeText
}

// 输出处理后文本
func (t *Translate) Render() string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, t.Doc)
	s := buf.String()
	replacer := strings.NewReplacer("<html><head></head><body>", "", "</body></html>", "") // Remove <html><head></head><body></body></html>

	return replacer.Replace(s)
}

// 翻译处理
func (t *Translate) Do() *Translate {
	for i, node := range t.Nodes {
		s := t.Anatomy(node)
		if t.Config.Debug {
			log.Println(fmt.Sprintf("#%v: %#v", i+1, s))
		}
		t.currentNode.Data = t.currentNodeText
	}

	return t
}

// 更新账号状态
func (t *Translate) updateAccount(pid string, enable bool) (bool, error) {
	err := t.Viper.ReadInConfig()
	if err != nil {
		log.Panic("Config file not found(%v)", err)
		return false, err
	}

	key := "accounts"
	t.Viper.Get(key)
	d := time.Now()
	ym, _ := strconv.Atoi(fmt.Sprintf("%d%02d", d.Year(), int(d.Month())))
	for k, v := range t.Config.Accounts {
		if v.PID == pid {
			fmt.Println("Update " + pid)
			v.Enabled = enable
			v.YearMonth = ym
			t.Config.Accounts[k] = v
			continue
		}
		// 刷新其他账号
		if v.Enabled == false && v.YearMonth != ym {
			v.Enabled = true
			v.YearMonth = ym
			t.Config.Accounts[k] = v
		}
	}
	t.Viper.Set(key, t.Config.Accounts)
	err = t.Viper.WriteConfig()
	//err = t.Viper.WriteConfigAs("./src/config/config.bak.json")
	if err == nil {
		return true, nil
	} else {
		fmt.Println(err)
		return false, err
	}

}

// 获取一个随机有效账号
func (t *Translate) GetRandomAccount() Account {
	rawAccounts := t.Config.Accounts
	n := len(rawAccounts)
	if n == 0 {
		log.Panic("请设置翻译账号列表。")
	}

	accounts := make([]Account, 0)
	for _, v := range rawAccounts {
		if v.Enabled {
			accounts = append(accounts, v)
		}
	}

	n = len(accounts)
	if n == 0 {
		log.Panic("暂无有效的翻译账号。")
	}

	return accounts[rand.Intn(n)]
}
