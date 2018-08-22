package translate

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"bytes"
	"io"
	"fmt"
	"log"
	"strings"
)

type IT interface {
	Parse() (string, error)
	Do()
	SetRawContent(s string)
}

type Translate struct {
	IT
	Debug           bool
	From            string // 来源语言
	To              string // 目标语言
	Doc             *html.Node
	Nodes           []*html.Node
	rawContent      string
	ignoreAtoms     []atom.Atom
	currentNode     *html.Node
	currentNodeText string
	Languages       map[string]string // 支持的语种
}

func NewTranslate(debug bool) *Translate {
	return &Translate{
		Debug: debug,
	}
}

// 设置要翻译的文本内容
func (t *Translate) SetRawContent(s string) *Translate {
	t.rawContent = s
	t.ignoreAtoms = []atom.Atom{atom.Ping}

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
	html.Render(w, t.Doc.FirstChild.LastChild.FirstChild) // Remove <html><head></head><body></body></html>

	return buf.String()
}

// 翻译处理
func (t *Translate) Do() *Translate {
	for i, node := range t.Nodes {
		s := t.Anatomy(node)
		if t.Debug {
			log.Println(fmt.Sprintf("#%v: %#v", i+1, s))
		}
		t.currentNode.Data = t.currentNodeText
	}

	return t
}
