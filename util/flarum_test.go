package util

import (
	"testing"
)

func TestReadLocale(t *testing.T) {
	localeDataArr := FlarumReadLocale("../view/flarum", "../view/extensions", "../view/locale", "en")
	// fmt.Println(localeDataArr["core.forum.change_email.confirm_password_placeholder"])
	// fmt.Println(localeDataArr["core.ref.confirm_password"])
	// fmt.Println("page title", localeDataArr["core.lib.meta_titles.with_page_title"])

	if localeDataArr["core.forum.change_email.confirm_password_placeholder"] != "Confirm Password" {
		t.Errorf("Get wrong locale")
	}
	if localeDataArr["flarum-lock.forum.notifications.discussion_locked_text"] != "{username} locked" {
		t.Errorf("Get wrong locale")
	}
	if localeDataArr["flarum-tags.admin.edit_tag.submit_button"] != "Save Changes" {
		t.Errorf("Get wrong locale")
	}
	if localeDataArr["core.lib.meta_titles.with_page_title"] != "{pageNumber, plural, =1 {{pageTitle} - {forumName}} other {{pageTitle}: Page # - {forumName}}}" {
		t.Errorf("Get wrong locale")
	}
}

func TestReadLocaleZh(t *testing.T) {
	localeDataArr := FlarumReadLocale("../view/flarum", "../view/extensions", "../view/locale", "zh")

	if localeDataArr["core.forum.change_email.confirm_password_placeholder"] != "确认密码" {
		t.Errorf("Get wrong locale")
	}

	if (localeDataArr["fof-html-errors.admin.settings.error.403"]) != "403 拒绝访问" {
		t.Errorf("Get wrong locale")
	}

	if localeDataArr["core.lib.meta_titles.with_page_title"] != "{pageNumber, plural, =1 {{pageTitle} - {forumName}} other {{pageTitle}: Page # - {forumName}}}" {
		t.Errorf("Get wrong locale")
	}
}
