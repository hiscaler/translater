package translate

import "fmt"

// 搜狐翻译
type SohuTranslate struct {
	Translate
}

// 翻译处理
func (t *SohuTranslate) Do() *SohuTranslate {
	for i, node := range t.Nodes {
		fmt.Println(fmt.Sprintf("#%v: %#v", i+1, t.Anatomy(node)))
		t.currentNode.Data = t.currentNodeText + " [ Sohu 翻译后 ]"
	}
	return t
}
