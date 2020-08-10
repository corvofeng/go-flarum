package controller

import (
	"encoding/json"
	"goyoubbs/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/rs/xid"
)

func (h *BaseHandler) AdminLinkList(w http.ResponseWriter, r *http.Request) {
	lid := r.FormValue("lid")

	db := h.App.Db
	redisDB := h.App.RedisDB

	var lobj model.Link
	if len(lid) > 0 {
		_, err := strconv.ParseUint(lid, 10, 64)
		if err != nil {
			w.Write([]byte(`{"retcode":400,"retmsg":"key type err"}`))
			return
		}

		lobj = model.LinkGetByID(db, lid)
		if lobj.ID == 0 {
			w.Write([]byte(`{"retcode":404,"retmsg":"id not found"}`))
			return
		}
	} else {
		lobj.Score = 10
	}

	currentUser, _ := h.CurrentUser(w, r)
	if currentUser.ID == 0 {
		w.Write([]byte(`{"retcode":401,"retmsg":"authored err"}`))
		return
	}
	if !currentUser.IsAdmin() {
		w.Write([]byte(`{"retcode":403,"retmsg":"flag forbidden"}`))
		return
	}

	type pageData struct {
		BasePageData
		Items []model.Link
		Lobj  model.Link
	}

	tpl := h.CurrentTpl(r)
	evn := &pageData{}
	evn.SiteCf = h.App.Cf.Site
	evn.Title = "链接列表"
	evn.IsMobile = tpl == "mobile"
	evn.CurrentUser = currentUser
	evn.ShowSideAd = true
	evn.PageName = "user_list"

	evn.SiteInfo = model.GetSiteInfo(redisDB)

	evn.Items = model.LinkList(db, true)
	evn.Lobj = lobj

	token := h.GetCookie(r, "token")
	if len(token) == 0 {
		token := xid.New().String()
		h.SetCookie(w, "token", token, 1)
	}

	h.Render(w, tpl, evn, "layout.html", "adminlinklist.html")
}

func (h *BaseHandler) AdminLinkListPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	token := h.GetCookie(r, "token")
	if len(token) == 0 {
		w.Write([]byte(`{"retcode":400,"retmsg":"token cookie missed"}`))
		return
	}

	currentUser, _ := h.CurrentUser(w, r)
	if currentUser.ID == 0 {
		w.Write([]byte(`{"retcode":401,"retmsg":"authored err"}`))
		return
	}
	if !currentUser.IsAdmin() {
		w.Write([]byte(`{"retcode":403,"retmsg":"flag forbidden}`))
		return
	}

	type response struct {
		normalRsp
	}

	decoder := json.NewDecoder(r.Body)
	var rec model.Link
	err := decoder.Decode(&rec)
	if err != nil {
		w.Write([]byte(`{"retcode":400,"retmsg":"json Decode err:` + err.Error() + `"}`))
		return
	}
	defer r.Body.Close()

	rec.Name = strings.TrimSpace(rec.Name)
	rec.URL = strings.TrimSpace(rec.URL)

	if len(rec.Name) == 0 || len(rec.URL) == 0 {
		w.Write([]byte(`{"retcode":400,"retmsg":"missed arg"}`))
		return
	}

	model.LinkSet(h.App.Db, rec)

	rsp := response{}
	rsp.Retcode = 200
	json.NewEncoder(w).Encode(rsp)
}
