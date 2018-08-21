package main

import (
	"golang.org/x/net/html"
	"strings"
	"bytes"
	"io"
	"fmt"
	"golang.org/x/net/html/atom"
	"errors"
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

func NewBaseBuilder(t *Translate) *Translate {
	return &Translate{}
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
func (t *Translate) Parse() *Translate {
	doc, err := html.Parse(strings.NewReader(t.rawContent))
	if err != nil {
		errors.New(err.Error())
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

	return t
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

func main() {
	htmlString := `
<div>
	<p>This is p1</p>
	<p>This is p2</p>
	<a href="">Link1</a>
	<p>
		aaa<a href=""><em>Link2</em></a>
	</p>
	<p>       　　　</p>
    <ul>
		<li>T<p>p3</p></li>
		<li><p>p4</p></li>
		<li><p><img src="../img/one.png" alt="this is image" />IMAGE</p></li>
	</ul>
</div> 
`

	htmlString = `
<article itemprop="articleBody" data-uuid="64c76dce-d063-3f8a-9961-752c503879d8" data-type="story" data-reactid="3"><figure class="canvas-image Mx(a) canvas-atom Mt(0) Mt(20px)--sm Mb(24px) Mb(22px)--sm" style="max-width:100%;" data-type="image" itemscope="" itemprop="associatedMedia image" itemtype="http://schema.org/ImageObject" data-reactid="4"><meta itemprop="height" content="448" data-reactid="5"><meta itemprop="width" content="800" data-reactid="6"><div style="padding-bottom:56%;" class="Maw(100%) Pos(r) H(0)" data-reactid="7"><img alt="Trump Lashes Out At NYT&amp;apos;s Bombshell McGahn Report, Calls It &amp;apos;Fake News&amp;apos;" class=" StretchedBox W(100%) H(100%) ie-7_H(a)" itemprop="url" src="https://s.yimg.com/ny/api/res/1.2/qAfoUUPJKG6g0q2CJNaszw--~A/YXBwaWQ9aGlnaGxhbmRlcjtzbT0xO3c9ODAw/http://media.zenfs.com/en-US/homerun/the_huffington_post_584/9240b657bb573cd5b148e178332befc9" data-reactid="8"><noscript data-reactid="9"><img alt="Trump Lashes Out At NYT&amp;apos;s Bombshell McGahn Report, Calls It &amp;apos;Fake News&amp;apos;" class="StretchedBox W(100%) H(100%) ie-7_H(a)" src="https://s.yimg.com/ny/api/res/1.2/qAfoUUPJKG6g0q2CJNaszw--~A/YXBwaWQ9aGlnaGxhbmRlcjtzbT0xO3c9ODAw/http://media.zenfs.com/en-US/homerun/the_huffington_post_584/9240b657bb573cd5b148e178332befc9" itemprop="url"/></noscript></div><div itemprop="caption description" class="Ov(h) Pos(r) Ff(ss) Mah(80px)" data-reactid="10"><figcaption class="C($c-fuji-grey-h) Fz(13px) Py(5px) Lh(1.5)" title="President Donald Trump fired off several angry Tweets on Sunday aimed at" data-reactid="11"><div class="figure-caption" data-reactid="12">President Donald Trump fired off several angry Tweets on Sunday aimed at</div></figcaption><button class="C($c-fuji-blue-1-b) Cur(p) W(100%) T(63px) Bgc(#fff) Ta(start) Fz(13px) P(0) Bd(0) O(0) Lh(1.5) Pos(a)" data-reactid="13"><span data-reactid="14">More</span></button></div></figure><div class="canvas-body Wow(bw) Cl(start) Mb(20px) Lh(30px) Fz(18px) Ff(s) C(#000) D(i)" data-reactid="15"><p class="canvas-atom canvas-text Mb(1.0em) Mb(0)--sm Mt(0.8em)--sm" type="text" content="President <a href=&quot;https://www.huffingtonpost.com/topic/donald-trump&quot; rel=&quot;nofollow noopener&quot; target=&quot;_blank&quot;>Donald Trump </a>fired off several angry Tweets on Sunday aimed at discrediting&amp;nbsp;<a href=&quot;https://www.huffingtonpost.com/entry/don-mcgahn-white-house-counsel-mueller-russia_us_5b7872bfe4b05906b4144312&quot; rel=&quot;nofollow noopener&quot; target=&quot;_blank&quot;>a report that White House counsel</a> Donald McGahn is cooperating with the Russia investigation to protect himself, with Trump insisting that it is “just the opposite.”" data-reactid="16">President <a href="https://www.huffingtonpost.com/topic/donald-trump" rel="nofollow noopener" target="_blank">Donald Trump </a>fired off several angry Tweets on Sunday aimed at discrediting&nbsp;<a href="https://www.huffingtonpost.com/entry/don-mcgahn-white-house-counsel-mueller-russia_us_5b7872bfe4b05906b4144312" rel="nofollow noopener" target="_blank">a report that White House counsel</a> Donald McGahn is cooperating with the Russia investigation to protect himself, with Trump insisting that it is “just the opposite.”</p><p class="canvas-atom canvas-text Mb(1.0em) Mb(0)--sm Mt(0.8em)--sm" type="text" content="“I have demanded transparency so that this Rigged and Disgusting Witch Hunt can come to a close,”<a href=&quot;https://twitter.com/realDonaldTrump/status/1030940529037651968&quot; target=&quot;_blank&quot; rel=&quot;nofollow noopener&quot;> Trump tweeted</a>.&amp;nbsp;“I allowed [McGahn] and all others to testify ― I didn’t have to. I have nothing to hide.”" data-reactid="17">“I have demanded transparency so that this Rigged and Disgusting Witch Hunt can come to a close,”<a href="https://twitter.com/realDonaldTrump/status/1030940529037651968" target="_blank" rel="nofollow noopener"> Trump tweeted</a>.&nbsp;“I allowed [McGahn] and all others to testify ― I didn’t have to. I have nothing to hide.”</p><div class="canvas-tweet canvas-atom Ov(a) W(100%)" data-type="tweet" data-reactid="18"><twitterwidget class="twitter-tweet twitter-tweet-rendered" id="twitter-widget-0" style="position: static; visibility: visible; display: block; transform: rotate(0deg); max-width: 100%; width: 500px; min-width: 220px; margin-top: 10px; margin-bottom: 10px;" data-tweet-id="1031137499995930624"></twitterwidget></div><p class="canvas-atom canvas-text Mb(1.0em) Mb(0)--sm Mt(0.8em)--sm" type="text" content="Trump also said “some members of the media” were angry about the Times’ report and have called to “complain and apologize.” He did not identify who those members are.&amp;nbsp;" data-reactid="19">Trump also said “some members of the media” were angry about the Times’ report and have called to “complain and apologize.” He did not identify who those members are.&nbsp;</p><div class="canvas-tweet canvas-atom Ov(a) W(100%)" data-type="tweet" data-reactid="20"><twitterwidget class="twitter-tweet twitter-tweet-rendered" id="twitter-widget-1" style="position: static; visibility: visible; display: block; transform: rotate(0deg); max-width: 100%; width: 500px; min-width: 220px; margin-top: 10px; margin-bottom: 10px;" data-tweet-id="1031152483949838336"></twitterwidget></div><p class="canvas-atom canvas-text Mb(1.0em) Mb(0)--sm Mt(0.8em)--sm" type="text" content="&amp;nbsp;His outburst came the day after The <a href=&quot;https://www.nytimes.com/2018/08/18/us/politics/don-mcgahn-mueller-investigation.html&quot; target=&quot;_blank&quot; rel=&quot;nofollow noopener&quot;>New York Times reported</a> that McGahn is extensively cooperating with special counsel Robert Mueller’s investigation into possible links between Trump’s campaign and Russian interference in the 2016 presidential election and whether the president has obstructed justice. McGahn, according to the Times story, has voluntarily provided Mueller’s team a wealth of information about Trump’s behavior, beyond what might be expected from someone in his position." data-reactid="21">&nbsp;His outburst came the day after The <a href="https://www.nytimes.com/2018/08/18/us/politics/don-mcgahn-mueller-investigation.html" target="_blank" rel="nofollow noopener">New York Times reported</a> that McGahn is extensively cooperating with special counsel Robert Mueller’s investigation into possible links between Trump’s campaign and Russian interference in the 2016 presidential election and whether the president has obstructed justice. McGahn, according to the Times story, has voluntarily provided Mueller’s team a wealth of information about Trump’s behavior, beyond what might be expected from someone in his position.</p><p class="canvas-atom canvas-text Mb(1.0em) Mb(0)--sm Mt(0.8em)--sm" type="text" content="Any suggestion that McGahn has turned on him is false, <a href=&quot;https://twitter.com/realDonaldTrump/status/1031150465759633408&quot; target=&quot;_blank&quot; rel=&quot;nofollow noopener&quot;>Trump insisted </a>on Sunday." data-reactid="22">Any suggestion that McGahn has turned on him is false, <a href="https://twitter.com/realDonaldTrump/status/1031150465759633408" target="_blank" rel="nofollow noopener">Trump insisted </a>on Sunday.</p><p class="canvas-atom canvas-text Mb(1.0em) Mb(0)--sm Mt(0.8em)--sm" type="text" content="″In fact it is just the opposite,” <a href=&quot;https://twitter.com/realDonaldTrump/status/1031150465759633408&quot; target=&quot;_blank&quot; rel=&quot;nofollow noopener&quot;>he tweeted.</a> “This is why the Fake News Media has become the Enemy of the People. So bad for America!”" data-reactid="23">″In fact it is just the opposite,” <a href="https://twitter.com/realDonaldTrump/status/1031150465759633408" target="_blank" rel="nofollow noopener">he tweeted.</a> “This is why the Fake News Media has become the Enemy of the People. So bad for America!”</p><p class="canvas-atom canvas-text Mb(1.0em) Mb(0)--sm Mt(0.8em)--sm" type="text" content="The Times reported that McGahn’s cooperation in the investigation began in part because the president’s initial team of lawyers insisted that Trump had nothing to hide and they wanted the investigation to end quickly." data-reactid="24">The Times reported that McGahn’s cooperation in the investigation began in part because the president’s initial team of lawyers insisted that Trump had nothing to hide and they wanted the investigation to end quickly.</p><p class="canvas-atom canvas-text Mb(1.0em) Mb(0)--sm Mt(0.8em)--sm" type="text" content="But according to sources close to McGahn, cited in the Times’ report, both McGahn and his lawyer, William Burck, were confused by Trump’s willingness to allow him to speak so freely to the special counsel. Trump’s attitude reportedly led McGahn to suspect that the president was setting him up to take the blame for any possible illegal acts of obstruction." data-reactid="25">But according to sources close to McGahn, cited in the Times’ report, both McGahn and his lawyer, William Burck, were confused by Trump’s willingness to allow him to speak so freely to the special counsel. Trump’s attitude reportedly led McGahn to suspect that the president was setting him up to take the blame for any possible illegal acts of obstruction.</p><p class="canvas-atom canvas-text Mb(1.0em) Mb(0)--sm Mt(0.8em)--sm" type="text" content="Nixon administration attorney John Dean, who pleaded guilty to conspiracy to obstruct justice during the 1970s Watergate scandal,<a href=&quot;https://www.huffingtonpost.com/entry/nixon-attorney-john-dean-says-mcgahn-was-smart-to-cooperate-with-mueller-probe_us_5b78d9aae4b0a5b1febbf635&quot; rel=&quot;nofollow noopener&quot; target=&quot;_blank&quot;> told Slate on Saturday </a>that McGahn did the “right thing” to cooperate “extensively” with Mueller’s investigation." data-reactid="26">Nixon administration attorney John Dean, who pleaded guilty to conspiracy to obstruct justice during the 1970s Watergate scandal,<a href="https://www.huffingtonpost.com/entry/nixon-attorney-john-dean-says-mcgahn-was-smart-to-cooperate-with-mueller-probe_us_5b78d9aae4b0a5b1febbf635" rel="nofollow noopener" target="_blank"> told Slate on Saturday </a>that McGahn did the “right thing” to cooperate “extensively” with Mueller’s investigation.</p><p class="canvas-atom canvas-text Mb(1.0em) Mb(0)--sm Mt(0.8em)--sm" type="text" content="McGahn is “doing exactly the right thing, not merely to protect himself, but to protect his client. And his client is not Donald Trump; his client is the office of the president,” said Dean, who was White House counsel under Richard Nixon." data-reactid="27">McGahn is “doing exactly the right thing, not merely to protect himself, but to protect his client. And his client is not Donald Trump; his client is the office of the president,” said Dean, who was White House counsel under Richard Nixon.</p><p class="canvas-atom canvas-text Mb(1.0em) Mb(0)--sm Mt(0.8em)--sm" type="text" content="Trump, in his tweets on Sunday, referred to Dean as a “rat.”" data-reactid="28">Trump, in his tweets on Sunday, referred to Dean as a “rat.”</p><p class="canvas-atom canvas-text Mb(1.0em) Mb(0)--sm Mt(0.8em)--sm" type="text" content="Although Dean was part of the White House-led efforts under President Richard Nixon to cover up the crimes committed Watergate scandal, his decision to testify to Congress about it was crucial to exposing the affair and causing Nixon to resign from office 44 years ago." data-reactid="29">Although Dean was part of the White House-led efforts under President Richard Nixon to cover up the crimes committed Watergate scandal, his decision to testify to Congress about it was crucial to exposing the affair and causing Nixon to resign from office 44 years ago.</p><ul class="canvas-list List(d)" data-type="list" data-reactid="30"><li data-reactid="31">This article originally appeared on <a href="https://www.huffingtonpost.com/entry/trump-responds-to-mcgahn-report_us_5b7968e5e4b018b93e94c361">HuffPost</a>.</li></ul><div data-reactid="32"></div></div></article>`
	sohu := &SohuTranslate{}
	sohu.SetRawContent(htmlString).Parse()
	fmt.Println(sohu.Do().Render())

}
