package controller

import (
	"net/http"
	"strconv"
	"time"
	"zoe/lib/qqOAuth"
)

func (h *BaseHandler) QQOauthHandler(w http.ResponseWriter, r *http.Request) {
	scf := h.App.Cf.Site
	qq, err := qqOAuth.NewQQOAuth(strconv.Itoa(scf.QQClientID), scf.QQClientSecret, scf.MainDomain+"/oauth/qq/callback")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	// qqOAuth.Logging = true

	now := time.Now().UTC().Unix()
	qqURLState := strconv.FormatInt(now, 10)[6:]

	urlStr, err := qq.GetAuthorizationURL(qqURLState)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	h.SetCookie(w, "QQURLState", qqURLState, 1)
	http.Redirect(w, r, urlStr, http.StatusSeeOther)
}
