package controller

import (
	"encoding/json"
	"goyoubbs/model"
	"net/http"
	"strconv"

	"github.com/rs/xid"
	"goji.io/pat"
)

func (h *BaseHandler) CommentEdit(w http.ResponseWriter, r *http.Request) {
	_aid, cid := pat.Param(r, "aid"), pat.Param(r, "cid")
	aid, err := strconv.ParseUint(_aid, 10, 64)
	if err != nil {
		w.Write([]byte(`{"retcode":400,"retmsg":"aid type err"}`))
		return
	}
	cidI, err := strconv.ParseUint(cid, 10, 64)
	if err != nil {
		w.Write([]byte(`{"retcode":400,"retmsg":"cid type err"}`))
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

	sqlDB := h.App.MySQLdb
	redisDB := h.App.RedisDB

	// aobj, _ := model.ArticleGetByID(db, aid)
	aobj, _ := model.SQLArticleGetByID(sqlDB, redisDB, aid)

	// comment
	cobj, err := model.SQLGetCommentByID(sqlDB, redisDB, cidI, h.App.Cf.Site.TimeZone)
	if err != nil {
		w.Write([]byte(`{"retcode":404,"retmsg":"` + err.Error() + `"}`))
		return
	}

	act := r.FormValue("act")

	if act == "del" {
		// remove
		// model.CommentDelByKey(db, aid, cidI)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	type pageData struct {
		PageData
		Aobj model.Article
		Cobj model.Comment
	}

	tpl := h.CurrentTpl(r)
	evn := &pageData{}
	evn.SiteCf = h.App.Cf.Site
	evn.Title = "修改评论"
	evn.IsMobile = tpl == "mobile"
	evn.CurrentUser = currentUser
	evn.ShowSideAd = true
	evn.PageName = "comment_edit"
	evn.SiteInfo = model.GetSiteInfo(redisDB)

	evn.Aobj = aobj
	evn.Cobj = cobj

	h.SetCookie(w, "token", xid.New().String(), 1)
	h.Render(w, tpl, evn, "layout.html", "admincommentedit.html")
}

func (h *BaseHandler) CommentEditPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	_aid, cid := pat.Param(r, "aid"), pat.Param(r, "cid")
	aid, err := strconv.ParseUint(_aid, 10, 64)
	if err != nil {
		w.Write([]byte(`{"retcode":400,"retmsg":"aid type err"}`))
		return
	}
	cidI, err := strconv.ParseUint(cid, 10, 64)
	if err != nil {
		w.Write([]byte(`{"retcode":400,"retmsg":"cid type err"}`))
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

	sqlDB := h.App.MySQLdb
	redisDB := h.App.RedisDB

	// comment
	cobj, err := model.SQLGetCommentByID(sqlDB, redisDB, cidI, h.App.Cf.Site.TimeZone)
	if err != nil {
		w.Write([]byte(`{"retcode":404,"retmsg":"` + err.Error() + `"}`))
		return
	}

	type recForm struct {
		Act     string `json:"act"`
		Content string `json:"content"`
	}

	decoder := json.NewDecoder(r.Body)
	var rec recForm
	err = decoder.Decode(&rec)
	if err != nil {
		w.Write([]byte(`{"retcode":400,"retmsg":"json Decode err:` + err.Error() + `"}`))
		return
	}
	defer r.Body.Close()

	if rec.Act == "preview" {
		tmp := struct {
			normalRsp
			Html string `json:"html"`
		}{
			normalRsp{200, ""},
			model.ContentFmt(rec.Content),
		}
		json.NewEncoder(w).Encode(tmp)
		return
	}

	oldContent := cobj.Content

	if oldContent == rec.Content {
		w.Write([]byte(`{"retcode":201,"retmsg":"nothing changed"}`))
		return
	}

	cobj.Content = rec.Content

	// model.CommentSetByKey(db, aid, cidI, cobj)

	h.DelCookie(w, "token")

	tmp := struct {
		normalRsp
		Aid uint64 `json:"aid"`
	}{
		normalRsp{200, "ok"},
		aid,
	}
	json.NewEncoder(w).Encode(tmp)
}
