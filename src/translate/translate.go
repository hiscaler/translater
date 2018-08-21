package translate

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"strings"
	"bytes"
	"io"
	"fmt"
)

type IT interface {
	Parse() (string, error)
	Do()
	SetRawContent(s string)
}

type Translate struct {
	IT
	Doc             *html.Node
	Nodes           []*html.Node
	rawContent      string
	ignoreAtoms     []atom.Atom
	currentNode     *html.Node
	currentNodeText string
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

func (t *Translate) Render() string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, t.Doc)

	return buf.String()
}

// 翻译处理
func (t *Translate) Do() *Translate {
	for i, node := range t.Nodes {
		fmt.Println(fmt.Sprintf("#%v: %#v", i+1, t.Anatomy(node)))
		t.currentNode.Data = t.currentNodeText
	}

	return t
}
