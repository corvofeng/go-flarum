package model

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

// 测试多个引用存在时的检查结果
func TestMentionGet(t *testing.T) {

	mentionPost := Comment{}
	mentionPost.ID = 86
	mentionPost.AID = 13
	mentionPost.UserName = "helloworld"
	mentionPost.Number = 3
	user := User{}
	user.ID = 123
	user.Name = "helloworld"
	user.Nickname = "nick"

	TESTCASES := []struct {
		UserData    string
		MentionPost Comment
		MentionUser User
		ExpectData  string
	}{
		{"@helloworld", mentionPost, user, `<USERMENTION displayname="nick" id="123" username="helloworld">@helloworld</USERMENTION>`},
		{"@helloworld#86", mentionPost, user, `<POSTMENTION discussionid="13" displayname="nick" id="86" number="3" username="helloworld">@helloworld</POSTMENTION>`},
		// {"@helloworld#54", mentionPost, user, `@helloworld#54`},
	}
	for _, data := range TESTCASES {
		for _, mentionStr := range mentionRegexp.FindAllStringSubmatch(data.UserData, -1) {
			renderData := makeMention(mentionStr, data.MentionPost, data.MentionUser)
			if data.ExpectData != renderData {
				t.Error("Can't process mention", mentionStr, renderData)
			}
		}
	}
}

func TestMention(t *testing.T) {

	mentionPost := Comment{}
	mentionPost.ID = 86
	mentionPost.AID = 13
	mentionPost.UserName = "helloworld"
	mentionPost.Number = 3
	user := User{}
	user.ID = 123
	user.Name = "helloworld"
	user.Nickname = "nick"

	TESTCASES := []struct {
		UserData    string
		MentionPost Comment
		MentionUser User
		ExpectData  string
	}{
		{
			"@helloworld",
			mentionPost,
			user,
			`<USERMENTION displayname="nick" id="123" username="helloworld">@helloworld</USERMENTION>`,
		},
		{
			"@helloworld#86",
			mentionPost,
			user,
			`<POSTMENTION discussionid="13" displayname="nick" id="86" number="3" username="helloworld">@helloworld</POSTMENTION>`,
		},
		{
			"@helloworld#86  测试输入 @helloworld#54",
			mentionPost,
			user,
			`<POSTMENTION discussionid="13" displayname="nick" id="86" number="3" username="helloworld">@helloworld</POSTMENTION> 测试输入 @helloworld#54`,
		},
		{
			`@"一枚小猿"#p86 测试引用`,
			mentionPost,
			user,
			`<POSTMENTION discussionid="13" displayname="nick" id="86" number="3" username="一枚小猿">一枚小猿</POSTMENTION> 测试输入 @helloworld#54`,
		},
	}

	for _, data := range TESTCASES {
		mentionDict := make(map[string]string)
		for _, mentionStr := range mentionRegexp.FindAllStringSubmatch(data.UserData, -1) {
			replData := makeMention(mentionStr, data.MentionPost, data.MentionUser)
			t.Error(mentionStr, replData)
			mentionDict[mentionStr[0]] = replData
		}

		newPost := replaceAllMentions(data.UserData, mentionDict)
		if newPost != data.ExpectData {
			t.Errorf("Can't process user data %s, %+v, %s", data.UserData, mentionDict, newPost)
		}
	}
}

func TestMentionRender(t *testing.T) {

	TESTCASES := []struct {
		UserData   string
		ExpectData string
	}{
		{
			`<POSTMENTION discussionid="13" displayname="nick" id="86" number="3" username="helloworld">@helloworld</POSTMENTION>`,
			`<a href="/d/13/3" class="PostMention" data-id="86">@helloworld</a>`,
		},
		{
			`<USERMENTION displayname="nick" id="123" username="helloworld">@helloworld</USERMENTION>`,
			`<a href="/u/helloworld" class="UserMention">@helloworld</a>`,
		},
	}

	for _, content := range TESTCASES {
		renderData := MentionToHTML(content.UserData)
		if renderData != content.ExpectData {
			t.Errorf(
				"Can't get render for %s: expected: %s, get: %s",
				content.UserData, content.ExpectData, renderData)
		}
	}
}

func TestFlarumMention(t *testing.T) {
	flarumMentionRegexp = regexp.MustCompile(`&lt;(USER|POST)MENTION(.+?)MENTION&gt;`)

	data := `&lt;USERMENTION displayname="corvofeng" id="87" username="corvofeng"&gt;@corvofeng&lt;/USERMENTION&gt;引用这条主题`
	replDict := make(map[string]string)
	for _, m := range flarumMentionRegexp.FindAllString(data, -1) {
		oldData := m
		m = strings.Replace(m, "&lt;", "<", -1)
		m = strings.Replace(m, "&gt;", ">", -1)
		replDict[oldData] = MentionToHTML(m)
	}
	for k, v := range replDict {
		data = strings.ReplaceAll(data, k, v)
	}

	fmt.Println(data)
}
