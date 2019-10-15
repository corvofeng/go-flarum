package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"goyoubbs/model"

	"goji.io/pat"
)

// CategoryDetailNew 新版的使用sql的页面
func (h *BaseHandler) CategoryDetailNew(w http.ResponseWriter, r *http.Request) {

	var start uint64
	var err error
	var cobj model.Category
	var page uint64
	var pageInfo model.ArticlePageInfo

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
	db := h.App.Db
	scf := h.App.Cf.Site
	sqlDB := h.App.MySQLdb
	logger := h.App.Logger

	if h.InAPI {
		cname := r.FormValue("tab")
		cobj, err = model.SQLCategoryGetByName(sqlDB, cname)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		p := r.FormValue("page")
		if len(p) > 0 {
			page, err = strconv.ParseUint(p, 10, 64)
			if err != nil {
				w.Write([]byte(`{"retcode":400,"retmsg":"page type err"}`))
				return
			}
		}
	} else {
		cid := pat.Param(r, "cid")
		_, err = strconv.Atoi(cid)
		if err != nil {
			w.Write([]byte(`{"retcode":400,"retmsg":"cid type err"}`))
			return
		}
		cobj, err = model.SQLCategoryGetByID(sqlDB, cid)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
	}

	tpl := h.CurrentTpl(r)

	type pageData struct {
		PageData
		Cobj     model.Category
		PageInfo model.ArticlePageInfo
	}

	evn := &pageData{}

	if h.InAPI {
		pageInfo = model.SQLCIDArticleListByPage(sqlDB, db, cobj.ID, page, uint64(scf.HomeShowNum), scf.TimeZone)
	} else {
		pageInfo = model.SQLCIDArticleList(sqlDB, db, cobj.ID, start, btn, uint64(scf.HomeShowNum), scf.TimeZone)
	}
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
					ID:    a.ID,
					UID:   a.UID,
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
