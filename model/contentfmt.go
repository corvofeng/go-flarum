package model

import (
	"fmt"
	"regexp"
	"strings"
	"zoe/util"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
)

var (
	// codeRegexp = regexp.MustCompile("(?s:```(.+?)```)")
	// imgRegexp     = regexp.MustCompile(`(https?://[\w./:]+/[\w./]+\.(jpg|jpe|jpeg|gif|png))`)
	// gistRegexp    = regexp.MustCompile(`(https?://gist\.github\.com/([a-zA-Z0-9-]+/)?[\d]+)`)
	gistRegexp    = regexp.MustCompile(`(https?://gist\.github\.com/([a-zA-Z0-9-_]+/)?[a-zA-Z\d]+)`)
	mentionRegexp = regexp.MustCompile(`\B@\"?([a-zA-Z0-9\p{Han}]{1,32})\"?#?p?([0-9]*)?`)
	// flarumMentionRegexp = regexp.MustCompile(`&lt;[USER|POST]MENTION(.+?)\/MENTION&gt;`)

	flarumMentionRegexp = regexp.MustCompile(`<(USER|POST)MENTION(.+?)MENTION>`)
	urlRegexp           = regexp.MustCompile(`([^;"='>])(https?://[^\s<]+[^\s<.)])`)
	nlineRegexp         = regexp.MustCompile(`\s{2,}`)
	youku1Regexp        = regexp.MustCompile(`https?://player\.youku\.com/player\.php/sid/([a-zA-Z0-9=]+)/v\.swf`)
	youku2Regexp        = regexp.MustCompile(`https?://v\.youku\.com/v_show/id_([a-zA-Z0-9=]+)(/|\.html?)?`)
	iframeStyle         = `style="min-width: 200px; width: 80%; height: 460px;" allowfullscreen="allowfullscreen" scrolling="no" frameborder="0"`
)

func htmlEscape(rs string) string {
	rs = strings.Replace(rs, "<", "&lt;", -1)
	rs = strings.Replace(rs, ">", "&gt;", -1)
	return rs
}

func genBilibili(url string) string {

	const style = `style="min-width: 200px; width: 80%; height: 460px;" allowfullscreen="allowfullscreen" frameborder="0"`
	return fmt.Sprintf(
		`<br><iframe %s src="%s" sandbox="allow-top-navigation allow-same-origin allow-forms allow-popups allow-scripts"></iframe><br>`,
		style,
		url,
	)
}

func genGist(url string) string {
	return fmt.Sprintf(
		`<br><iframe %s seamless="seamless" srcdoc='<html><body><style type="text/css">.gist .gist-data { height: 400px; }</style><script src="%s.js"></script></body></html>'></iframe><br>`,
		iframeStyle,
		url,
	)
}

func genYoutube(url string) string {
	return fmt.Sprintf(
		`<br><iframe %s src="%s" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"></iframe><br>`,
		iframeStyle,
		url,
	)
}

// ContentFmt 处理markdown样式
func ContentFmt(input string) string {
	return ContentRich(input)
}

type urlInfo struct {
	Href  string
	Click string
}

// ContentRich 用来转换文本, 转义以及允许用户添加一些富文本样式
// 该函数效率奇差, 但不会优化
func ContentRich(input string) string {
	input = strings.TrimSpace(input)
	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.Safelink
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	replaceDict := make(map[string]string) // 用来记录使用了uuid替换的数据的具体值

	// 处理mention信息
	// 首先获取应该被识别的mention信息
	// 参考: https://stackoverflow.com/a/39102969
	if strings.Contains(input, "USERMENTION") || strings.Contains(input, "POSTMENTION") { // flarum 的mention
		mentionDict := make(map[string]string)
		for _, m := range flarumMentionRegexp.FindAllString(input, -1) {
			uid := util.GetUUID()
			replaceDict[uid] = MentionToHTML(m)
			mentionDict[m] = uid
		}
		for k, v := range mentionDict {
			input = strings.ReplaceAll(input, k, v)
		}
	}
	if strings.Contains(input, "//player.bilibili.com") {
		bilibiliDict := make(map[string]string)
		bilibliRegexp := regexp.MustCompile(`<iframe src="(//player.bilibili.com[^"^\n]*)"[a-zA-Z0-9 ="]*>\s*</iframe>`)
		bilibiliURLRegexp := regexp.MustCompile(`(//player.bilibili.com[^"^\n]*)\n`)
		input = bilibliRegexp.ReplaceAllString(input, "$1\n") // 将原有的iframe包裹的块剥离开来
		for _, m := range bilibiliURLRegexp.FindAllString(input, -1) {
			m = strings.TrimSuffix(m, "\n")
			uid := util.GetUUID()
			replaceDict[uid] = genBilibili(m)
			bilibiliDict[m] = uid
		}

		for k, v := range bilibiliDict {
			input = strings.ReplaceAll(input, k, v)
		}
	}

	if strings.Contains(input, "://gist") {
		gistDict := make(map[string]string)
		embedRegexp := regexp.MustCompile(`<script src="(https://gist.github.com/([a-zA-Z0-9-_]+/)?[a-zA-Z\d]+).js"></script>`)
		gistURLRegexp := regexp.MustCompile(`(https?://gist\.github\.com/([a-zA-Z0-9-_]+/)?[a-zA-Z\d]+)`)
		input = embedRegexp.ReplaceAllString(input, "$1")
		for _, m := range gistURLRegexp.FindAllString(input, -1) {
			uid := util.GetUUID()
			replaceDict[uid] = genGist(m)
			gistDict[m] = uid
		}

		for k, v := range gistDict {
			input = strings.ReplaceAll(input, k, v)
		}
	}
	if strings.Contains(input, "://www.youtube.com") {
		youtubeDict := make(map[string]string)
		youtubeRegexp := regexp.MustCompile(`<iframe.*src="((https:)?//www.youtube.com[^"]*)".*>\s*</iframe>`)
		input = youtubeRegexp.ReplaceAllString(input, "$1\n") // 将原有的iframe包裹的块剥离开来
		youtubeURLRegexp := regexp.MustCompile(`((https:)?//www.youtube.com[^"^\n]*)\n`)

		for _, m := range youtubeURLRegexp.FindAllString(input, -1) {
			m = strings.TrimSuffix(m, "\n")
			uid := util.GetUUID()
			replaceDict[uid] = genYoutube(m)
			youtubeDict[m] = uid
		}

		for k, v := range youtubeDict {
			input = strings.ReplaceAll(input, k, v)
		}
	}

	// 将原有的字符串中的<>全部进行转义
	input = htmlEscape(input)

	// 对markdown文本进行解析
	input = string(markdown.ToHTML([]byte(input), nil, renderer))

	// 将原有被替换成uuid的内容进行恢复
	for k, v := range replaceDict {
		input = strings.ReplaceAll(input, k, v)
	}

	// input = strings.Replace(input, "\r\n", "\n", -1)
	// input = strings.Replace(input, "\r", "\n", -1)

	// input = nlineRegexp.ReplaceAllString(input, "</p><p>")
	// input = strings.Replace(input, "\n", "<br>", -1)

	input = "<p>" + input + "</p>"
	// input = strings.Replace(input, "<p></p>", "", -1)

	return input
}
