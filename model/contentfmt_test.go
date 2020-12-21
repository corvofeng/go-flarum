package model

import (
	"strings"
	"testing"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	// "github.com/microcosm-cc/bluemonday"
	// "github.com/russross/blackfriday"
)

func TestContentFMT(t *testing.T) {
	// t.Error("Hello world")

	str := "```go" + `
func getTrue() bool {
    return true
}` + "```"

	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.SkipHTML
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	html := string(markdown.ToHTML([]byte(str), nil, renderer))
	html = strings.Replace(html, "\n", "<br>", -1)
	if false {
		t.Error(html)
	}

	// unsafe := blackfriday.Run([]byte(str))
	// html = string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
	// if false {
	// 	t.Error(html)
	// }
}

func TestNewLine(t *testing.T) {
	// t.Error("Hello world")

	str := `hello

world
`
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	html := string(markdown.ToHTML([]byte(str), nil, renderer))
	if false {
		t.Error(html)
	}
}

func TestSkipHTML(t *testing.T) {
	str := `<a href="/d/939/4" class="PostMention" data-id="141">@corvofeng</a>`
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	html := string(markdown.ToHTML([]byte(str), nil, renderer))
	t.Log(html)
}

func TestGistHTML(t *testing.T) {
	str := `Add gist:https://gist.github.com/corvofeng/f16a1fab54757f32879ae4d326c27518`
	html := ContentFmt(str)
	t.Log(html)
}

func TestBilibiliHTML(t *testing.T) {
	output := `<br><iframe style="min-width: 200px; width: 80%; height: 460px;" allowfullscreen="allowfullscreen" frameborder="0" src="//player.bilibili.com/player.html?aid=755694813&bvid=BV1R64y1f77F&cid=268573403&page=1" sandbox="allow-top-navigation allow-same-origin allow-forms allow-popups allow-scripts"></iframe><br>`
	str := `<iframe src="//player.bilibili.com/player.html?aid=755694813&bvid=BV1R64y1f77F&cid=268573403&page=1" scrolling="no" border="0" frameborder="no" framespacing="0" allowfullscreen="true"> </iframe>`
	result := ContentRich(str)

	if strings.Index(result, output) < 0 {
		t.Error(result)
	}
	// 两个iframe
	output = `<iframe style="min-width: 200px; width: 80%; height: 460px;" allowfullscreen="allowfullscreen" frameborder="0" src="//player.bilibili.com/player.html?aid=288149957&bvid=BV1Zf4y1Y7dE&cid=268125615&page=1" sandbox="allow-top-navigation allow-same-origin allow-forms allow-popups allow-scripts"></iframe><br>` + "\n" + `<br><iframe style="min-width: 200px; width: 80%; height: 460px;" allowfullscreen="allowfullscreen" frameborder="0" src="//player.bilibili.com/player.html?aid=372300761&bvid=BV1pZ4y1576P&cid=240628593&page=1" sandbox="allow-top-navigation allow-same-origin allow-forms allow-popups allow-scripts"></iframe>`
	str = `<iframe src="//player.bilibili.com/player.html?aid=288149957&bvid=BV1Zf4y1Y7dE&cid=268125615&page=1" scrolling="no" border="0" frameborder="no" framespacing="0" allowfullscreen="true"> </iframe><iframe src="//player.bilibili.com/player.html?aid=372300761&bvid=BV1pZ4y1576P&cid=240628593&page=1" scrolling="no" border="0" frameborder="no" framespacing="0" allowfullscreen="true"> </iframe>`
	result = ContentRich(str)
	if strings.Index(result, output) < 0 {
		t.Error(result)
	}
}

func TestYoutubeHTML(t *testing.T) {
	output := `<br><iframe style="min-width: 200px; width: 80%; height: 460px;" allowfullscreen="allowfullscreen" scrolling="no" frameborder="0" src="https://www.youtube.com/embed/NIfdoxxeJb4?start=15" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"></iframe><br>`
	str := `<iframe width="560" height="315" src="https://www.youtube.com/embed/NIfdoxxeJb4?start=15" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>`
	result := ContentRich(str)
	if strings.Index(result, output) < 0 {
		t.Error(result)
	}
}

func TestGistReplace(t *testing.T) {
	output := `<br><iframe style="min-width: 200px; width: 80%; height: 460px;" allowfullscreen="allowfullscreen" scrolling="no" frameborder="0" seamless="seamless" srcdoc='<html><body><style type="text/css">.gist .gist-data { height: 400px; }</style><script src="https://gist.github.com/corvofeng/f16a1fab54757f32879ae4d326c27518.js"></script></body></html>'></iframe><br>`
	str := `<script src="https://gist.github.com/corvofeng/f16a1fab54757f32879ae4d326c27518.js"></script>`
	result := ContentRich(str)
	if strings.Index(result, output) < 0 {
		t.Error(result)
	}

	str = `https://gist.github.com/corvofeng/f16a1fab54757f32879ae4d326c27518`

	result = ContentRich(str)
	if strings.Index(result, output) < 0 {
		t.Error(result)
	}
}

func TestPostMentions(t *testing.T) {
	output := `<p><a href="/d/939/4" class="PostMention" data-id="141">@corvofeng</a>这是对代码块评论的回复</p>`
	str := `<POSTMENTION discussionid="939" displayname="一枚小猿" id="141" number="4" username="corvofeng">@corvofeng</POSTMENTION>这是对代码块评论的回复`
	result := ContentRich(str)
	if strings.Index(result, output) < 0 {
		t.Error(result)
	}
}

func TestInject(t *testing.T) {

	output := `lt;script src=&ldquo;https&rdquo;&amp;gt;&amp;lt;/script&amp;gt`
	str := `<script src="https"></script>`
	result := ContentRich(str)

	if strings.Index(result, output) < 0 {
		t.Error(result)
	}
}

func TestMarkdownRender(t *testing.T) {
	// Not support yet...
	// https://github.com/gomarkdown/markdown/issues/149

	str := `> Blockquotes can also be nested...
>> ...by using additional greater-than signs right next to each other...
> > > ...or with spaces between arrows.
`
	result := ContentRich(str)
	if false {
		t.Log(result)
	}

	if false {
		t.Error(result)
	}
}
