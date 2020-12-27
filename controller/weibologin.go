package controller

import (
	"net/http"
	"strconv"
	"time"
	"zoe/lib/weiboOAuth"
)

func (h *BaseHandler) WeiboOauthHandler(w http.ResponseWriter, r *http.Request) {
	scf := h.App.Cf.Site
	weibo, err := weiboOAuth.NewWeiboOAuth(strconv.Itoa(scf.WeiboClientID), scf.WeiboClientSecret, scf.MainDomain+"/oauth/wb/callback")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	// weiboOAuth.Logging = true

	now := time.Now().UTC().Unix()
	WeiboURLState := strconv.FormatInt(now, 10)[6:]

	urlStr, err := weibo.GetAuthorizationURL(WeiboURLState)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	h.SetCookie(w, "WeiboURLState", WeiboURLState, 1)
	http.Redirect(w, r, urlStr, http.StatusSeeOther)
}
