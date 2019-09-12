package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"goyoubbs/model"
	"github.com/ego008/youdb"
	"goji.io/pat"
)

// CategoryDetailNew 新版的使用sql的页面
func (h *BaseHandler) CategoryDetailNew(w http.ResponseWriter, r *http.Request) {

	var start uint64
	var err error

	btn, key, score := r.FormValue("btn"), r.FormValue("key"), r.FormValue("score")
	if len(key) > 0 {
		start, err = strconv.ParseUint(key, 10, 64)
		if err != nil {
			w.Write([]byte(`{"retcode":400,"retmsg":"key type err"}`))
			return
		}
	}
	if len(score) > 0 {
		_, err = strconv.ParseUint(score, 10, 64)
		if err != nil {
			w.Write([]byte(`{"retcode":400,"retmsg":"score type err"}`))
			return
		}
	}
	cid := pat.Param(r, "cid")
	_, err = strconv.Atoi(cid)
	if err != nil {
		w.Write([]byte(`{"retcode":400,"retmsg":"cid type err"}`))
		return
	}

	tpl := h.CurrentTpl(r)

	type pageData struct {
		PageData
		Cobj     model.Category
		PageInfo model.ArticlePageInfo
	}

	evn := &pageData{}
	db := h.App.Db
	scf := h.App.Cf.Site
	sqlDB := h.App.MySQLdb
	logger := h.App.Logger

	cobj, err := model.SQLCategoryGetById(sqlDB, cid)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	pageInfo := model.SQLCidArticleList(sqlDB, db, cobj.Id, start, btn, scf.HomeShowNum, scf.TimeZone)
	categories, err := model.SQLGetAllCategory(sqlDB)
	evn.Cobj = cobj
	evn.PageInfo = pageInfo

	evn.SiteCf = scf
	evn.Title = cobj.Name + " - " + scf.Name
	evn.Keywords = cobj.Name
	evn.Description = cobj.About
	evn.IsMobile = tpl == "mobile"

	currentUser, _ := h.CurrentUser(w, r)

	evn.CurrentUser = currentUser
	evn.ShowSideAd = true
	evn.PageName = "category_detail"
	evn.HotNodes = model.CategoryHot(db, scf.CategoryShowNum)
	evn.NewestNodes = categories
	if h.InAPI {
		type NodeData struct {
			model.RestfulAPIBase
			Data []model.RestfulTopic `json:"data"`
		}
		nodeData := NodeData{
			RestfulAPIBase: model.RestfulAPIBase{
				State: true,
			},
			Data: []model.RestfulTopic{},
		}
		// pageInfo
		for _, a := range pageInfo.Items {
			nodeData.Data = append(nodeData.Data,
				model.RestfulTopic{
					ID:    a.Id,
					UID:   a.Uid,
					Title: a.Title,
					Author: model.RestfulUser{
						Name:   a.Name,
						Avatar: a.Avatar,
					},
					CreateAt:   a.EditTimeFmt,
					VisitCount: a.ClickCnt,
				},
			)
		}
		logger.Debug("This is in api version")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(nodeData)
	} else {
		h.Render(w, tpl, evn, "layout.html", "category.html")
	}
}

func (h *BaseHandler) CategoryDetail(w http.ResponseWriter, r *http.Request) {
	btn, key, score := r.FormValue("btn"), r.FormValue("key"), r.FormValue("score")
	if len(key) > 0 {
		_, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			w.Write([]byte(`{"retcode":400,"retmsg":"key type err"}`))
			return
		}
	}
	if len(score) > 0 {
		_, err := strconv.ParseUint(score, 10, 64)
		if err != nil {
			w.Write([]byte(`{"retcode":400,"retmsg":"score type err"}`))
			return
		}
	}

	cid := pat.Param(r, "cid")
	_, err := strconv.Atoi(cid)
	if err != nil {
		w.Write([]byte(`{"retcode":400,"retmsg":"cid type err"}`))
		return
	}

	cmd := "zrscan"
	if btn == "prev" {
		cmd = "zscan"
	}

	db := h.App.Db
	scf := h.App.Cf.Site
	cobj, err := model.CategoryGetById(db, cid)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	currentUser, _ := h.CurrentUser(w, r)

	if cobj.Hidden && currentUser.Flag < 99 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"retcode":404,"retmsg":"not found"}`))
		return
	}
	cobj.Articles = db.Zget("category_article_num", youdb.I2b(cobj.Id)).Uint64()
	pageInfo := model.ArticleList(db, cmd, "category_article_timeline:"+cid, key, score, scf.HomeShowNum, scf.TimeZone)

	type pageData struct {
		PageData
		Cobj     model.Category
		PageInfo model.ArticlePageInfo
	}

	tpl := h.CurrentTpl(r)

	evn := &pageData{}
	evn.SiteCf = scf
	evn.Title = cobj.Name + " - " + scf.Name
	evn.Keywords = cobj.Name
	evn.Description = cobj.About
	evn.IsMobile = tpl == "mobile"

	evn.CurrentUser = currentUser
	evn.ShowSideAd = true
	evn.PageName = "category_detail"
	evn.HotNodes = model.CategoryHot(db, scf.CategoryShowNum)
	evn.NewestNodes = model.CategoryNewest(db, scf.CategoryShowNum)

	evn.Cobj = cobj
	evn.PageInfo = pageInfo

	// pageInfo = model.SqlArticleList(sqlDB, db, start, scf.HomeShowNum, scf.TimeZone)
	// End mysql

	h.Render(w, tpl, evn, "layout.html", "category.html")
}
