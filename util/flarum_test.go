package util

import (
	"testing"
)

func TestReadLocale(t *testing.T) {
	localeDataArr := FlarumReadLocale("../view/locale", "en")
	// fmt.Println(localeDataArr["core.forum.change_email.confirm_password_placeholder"])
	// fmt.Println(localeDataArr["core.ref.confirm_password"])

	if localeDataArr["core.forum.change_email.confirm_password_placeholder"] != "Confirm Password" {
		t.Errorf("Get wrong locale")
	}
	if localeDataArr["flarum-lock.forum.notifications.discussion_locked_text"] != "{username} locked" {
		t.Errorf("Get wrong locale")
	}
}

func TestMention(t *testing.T) {
	input := "@corvofeng#70 0715第一条评论"
	output := `<a href="/d/927" class="UserMention" data-id="70">corvofeng</a>0715第一条评论`
	input = mentionRegexp.ReplaceAllString(input, mentionReplaceStr)
	if input != output {
		t.Errorf("Get wrong mention render data:\n%s\n%s", input, output)
	}
}
