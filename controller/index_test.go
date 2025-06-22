package controller

import (
	"net/url"
	"testing"
)

func TestArticleGetter(t *testing.T) {
	t.Skip("Skipping testing in CI environment")
}

func TestFilePathJoin(t *testing.T) {
	f, err := url.JoinPath("https://baidu.com", "aaa", "bbb", "ccc")
	if err != nil || f != "https://baidu.com/aaa/bbb/ccc" {
		t.Errorf("File path join failed for %s", f)
	}

}
