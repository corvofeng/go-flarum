package controller

import (
	"fmt"
	"goyoubbs/model"
	"goyoubbs/model/flarum"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// ArticleHomeList 文章主页
func (h *BaseHandler) ArticleHomeList(w http.ResponseWriter, r *http.Request) {
	btn, key, score := r.FormValue("btn"), r.FormValue("key"), r.FormValue("score")
	var start uint64
	var err error
	var count uint64

	rsp := response{}
	if len(key) > 0 {
		start, err = strconv.ParseUint(key, 10, 64)
		if err != nil {
			rsp = response{400, "key type err"}
			h.Jsonify(w, rsp)
			return
		}
	}
	if len(score) > 0 {
		_, err = strconv.ParseUint(score, 10, 64)
		if err != nil {
			rsp = response{400, "scope type err"}
			h.Jsonify(w, rsp)
			return
		}
	}

	scf := h.App.Cf.Site

	type pageData struct {
		PageData
		SiteInfo   model.SiteInfo
		PageInfo   model.ArticlePageInfo
		Links      []model.Link
		FlarumInfo interface{}
	}

	sqlDB := h.App.MySQLdb
	redisDB := h.App.RedisDB
	// 获取全部的帖子数目
	err = sqlDB.QueryRow("SELECT COUNT(*) FROM topic").Scan(&count)
	if err != nil {
		log.Printf("Error %s", err)
		rsp = response{400, "Failed to get the count"}
		h.Jsonify(w, rsp)
		return
	}

	// 获取贴子列表
	pageInfo := model.SQLArticleList(sqlDB, redisDB, start, btn, uint64(scf.HomeShowNum), scf.TimeZone)
	categories, err := model.SQLGetAllCategory(sqlDB)

	tpl := h.CurrentTpl(r)
	evn := &pageData{}
	evn.SiteCf = scf
	evn.Title = scf.Name
	evn.Keywords = evn.Title
	evn.Description = scf.Desc
	evn.IsMobile = tpl == "mobile"
	currentUser, _ := h.CurrentUser(w, r)
	evn.CurrentUser = currentUser
	evn.ShowSideAd = false
	evn.PageName = "home"
	evn.NewestNodes = categories
	// evn.HotNodes = model.CategoryHot(db, scf.CategoryShowNum)
	// evn.NewestNodes = model.CategoryNewest(db, scf.CategoryShowNum)
	evn.SiteInfo = model.GetSiteInfo(redisDB)
	evn.PageInfo = pageInfo

	// 右侧的链接
	evn.Links = model.RedisLinkList(redisDB, false)

	h.Render(w, tpl, evn, "layout.html", "index.html")
}

// FlarumIndex flarum主页
func (h *BaseHandler) FlarumIndex(w http.ResponseWriter, r *http.Request) {
	btn, key, score := r.FormValue("btn"), r.FormValue("key"), r.FormValue("score")
	var start uint64
	var err error
	var count uint64

	rsp := response{}
	if len(key) > 0 {
		start, err = strconv.ParseUint(key, 10, 64)
		if err != nil {
			rsp = response{400, "key type err"}
			h.Jsonify(w, rsp)
			return
		}
	}
	if len(score) > 0 {
		_, err = strconv.ParseUint(score, 10, 64)
		if err != nil {
			rsp = response{400, "scope type err"}
			h.Jsonify(w, rsp)
			return
		}
	}

	scf := h.App.Cf.Site

	type pageData struct {
		PageData
		SiteInfo   model.SiteInfo
		PageInfo   model.ArticlePageInfo
		Links      []model.Link
		FlarumInfo interface{}
	}

	sqlDB := h.App.MySQLdb
	redisDB := h.App.RedisDB
	// 获取全部的帖子数目
	err = sqlDB.QueryRow("SELECT COUNT(*) FROM topic").Scan(&count)
	if err != nil {
		log.Printf("Error %s", err)
		rsp = response{400, "Failed to get the count"}
		h.Jsonify(w, rsp)
		return
	}

	// 获取贴子列表
	pageInfo := model.SQLArticleList(sqlDB, redisDB, start, btn, uint64(scf.HomeShowNum), scf.TimeZone)
	categories, err := model.SQLGetAllCategory(sqlDB)

	tpl := h.CurrentTpl(r)
	evn := &pageData{}
	evn.SiteCf = scf
	evn.Title = scf.Name
	evn.Keywords = evn.Title
	// evn.Description = scf.Desc
	// evn.IsMobile = tpl == "mobile"
	// evn.CurrentUser = currentUser
	// evn.ShowSideAd = false
	// evn.PageName = "home"
	// evn.NewestNodes = categories
	// evn.HotNodes = model.CategoryHot(db, scf.CategoryShowNum)
	// evn.NewestNodes = model.CategoryNewest(db, scf.CategoryShowNum)
	coreData := flarum.CoreData{}

	// 添加主站点信息
	coreData.AppendResourcs(model.FlarumCreateForumInfo(*h.App.Cf, evn.SiteInfo))

	// 添加所有分类的信息
	for _, category := range categories {
		coreData.AppendResourcs(
			model.FlarumCreateTag(category))
	}

	// 添加当前页面的帖子信息
	var res []flarum.Resource
	for _, article := range pageInfo.Items {
		diss := model.FlarumCreateDiscussion(article)
		coreData.AppendResourcs(diss)
		res = append(res, diss)
	}
	coreData.APIDocument.SetData(res)

	// 添加当前页面的帖子的用户信息, TODO: 用户有可能重复, 这里理论是需要优化的
	for _, article := range pageInfo.Items {
		user := model.FlarumCreateUser(article)
		coreData.AppendResourcs(user)
	}

	coreData.APIDocument.Links = make(map[string]string)
	coreData.APIDocument.Links["first"] = "https://flarum.yjzq.fun/api/v1/flarum/discussions?sort=&page%5Blimit%5D=20"

	// 设置语言信息
	coreData.Locales = make(map[string]string)
	coreData.Locales["en"] = "English"
	coreData.Locales["zh"] = "中文"
	coreData.Locale = "en"

	// 添加当前用户的session信息
	currentUser, err := h.CurrentUser(w, r)
	if err == nil {
		user := model.FlarumCreateCurrentUser(currentUser)
		coreData.AddSessionData(user, currentUser.RefreshCSRF(redisDB))
	}

	evn.FlarumInfo = coreData
	evn.SiteInfo = model.GetSiteInfo(redisDB)
	evn.PageInfo = pageInfo

	// 右侧的链接
	evn.Links = model.RedisLinkList(redisDB, false)

	h.Render(w, tpl, evn, "layout.html", "index.html")
}

// FlarumAPIDiscussions flarum文章api
func (h *BaseHandler) FlarumAPIDiscussions(w http.ResponseWriter, r *http.Request) {
	apiDoc := flarum.NewAPIDoc()
	var err error
	var page uint64

	// 需要返回的relations
	_include := r.FormValue("include")
	strings.Split(_include, ",")
	// fmt.Printf(include)
	page = 1

	scf := h.App.Cf.Site
	sqlDB := h.App.MySQLdb
	redisDB := h.App.RedisDB

	pageInfo := model.SQLCIDArticleListByPage(sqlDB, redisDB, 0, page, uint64(scf.HomeShowNum), scf.TimeZone)
	// categories, err := model.SQLGetAllCategory(sqlDB)
	if err != nil {
		fmt.Println(err)
	}

	// for _, category := range categories {
	// 	coreData.Resources = append(coreData.Resources,
	// 		model.FlarumCreateTag(category))
	// }

	var dissArr []flarum.Resource
	for _, article := range pageInfo.Items {
		diss := model.FlarumCreateDiscussion(article)
		dissArr = append(dissArr, diss)
	}
	apiDoc.SetData(dissArr)

	for _, article := range pageInfo.Items {
		user := model.FlarumCreateUser(article)
		apiDoc.Included = append(
			apiDoc.Included,
			user,
		)
	}

	apiDoc.Links["first"] = "https://flarum.yjzq.fun/api/v1/flarum/discussions?sort=&page%5Blimit%5D=20"
	apiDoc.Links["next"] = "https://flarum.yjzq.fun/api/v1/flarum/discussions?sort=&page%5Blimit%5D=20"

	h.Jsonify(w, apiDoc)
}
