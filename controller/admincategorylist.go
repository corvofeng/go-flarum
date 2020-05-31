package controller

import (
	"encoding/json"
	"goyoubbs/model"
	"net/http"
	"strconv"

	"github.com/rs/xid"
)

// AdminCategoryList 获取category列表
/*
 * w (http.ResponseWriter): TODO
 * r (*http.Request): TODO
 */
func (h *BaseHandler) AdminCategoryList(w http.ResponseWriter, r *http.Request) {
	cid, _, key := r.FormValue("cid"), r.FormValue("btn"), r.FormValue("key")
	if len(key) > 0 {
		_, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			w.Write([]byte(`{"retcode":400,"retmsg":"key type err"}`))
			return
		}
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

	var err error
	var cobj model.Category
	if len(cid) > 0 {
		cobj, err = model.SQLCategoryGetByID(sqlDB, cid)
		if err != nil {
			cobj = model.Category{}
		}
	}

	pageInfo := model.SQLCategoryList(sqlDB)

	type pageData struct {
		PageData
		PageInfo model.CategoryPageInfo
		Cobj     model.Category
	}

	tpl := h.CurrentTpl(r)
	evn := &pageData{}
	evn.SiteCf = h.App.Cf.Site
	evn.Title = "分类列表"
	evn.IsMobile = tpl == "mobile"
	evn.CurrentUser = currentUser
	evn.ShowSideAd = true
	evn.PageName = "category_list"

	evn.PageInfo = pageInfo
	evn.Cobj = cobj
	evn.SiteInfo = model.GetSiteInfo(redisDB)

	token := h.GetCookie(r, "token")
	if len(token) == 0 {
		token := xid.New().String()
		h.SetCookie(w, "token", token, 1)
	}

	h.Render(w, tpl, evn, "layout.html", "admincategorylist.html")
}

// AdminCategoryListPost 添加标签
/*
 * w (http.ResponseWriter): TODO
 * r (*http.Request): TODO
 */
func (h *BaseHandler) AdminCategoryListPost(w http.ResponseWriter, r *http.Request) {
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

	type recForm struct {
		CID    uint64 `json:"cid"`
		Name   string `json:"name"`
		About  string `json:"about"`
		Hidden string `json:"hidden"`
	}

	type response struct {
		normalRsp
	}

	decoder := json.NewDecoder(r.Body)
	var rec recForm
	err := decoder.Decode(&rec)
	if err != nil {
		w.Write([]byte(`{"retcode":400,"retmsg":"json Decode err:` + err.Error() + `"}`))
		return
	}
	defer r.Body.Close()

	if len(rec.Name) == 0 {
		w.Write([]byte(`{"retcode":400,"retmsg":"name is empty"}`))
		return
	}

	var hidden bool
	if rec.Hidden == "1" {
		hidden = true
	}

	var cobj model.Category
	// if rec.CID > 0 {
	// 	// edit
	// 	cobj, err = model.CategoryGetByID(db, strconv.FormatUint(rec.CID, 10))
	// 	if err != nil {
	// 		w.Write([]byte(`{"retcode":404,"retmsg":"cid not found"}`))
	// 		return
	// 	}
	// } else {
	// 	// add
	// 	newCID, _ := db.HnextSequence("category")
	// 	cobj.ID = newCID
	// }

	cobj.Name = rec.Name
	cobj.About = rec.About
	cobj.Hidden = hidden

	// jb, _ := json.Marshal(cobj)
	// db.Hset("category", youdb.I2b(cobj.ID), jb)

	rsp := response{}
	rsp.Retcode = 200
	json.NewEncoder(w).Encode(rsp)
}
