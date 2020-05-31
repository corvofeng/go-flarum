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
	if h.App.IsFlarum() {
		coreData := flarum.CoreData{}

		coreData.Resources = append(coreData.Resources,
			model.FlarumCreateForumInfo(*h.App.Cf.Site, evn.SiteInfo))

		for _, category := range categories {
			coreData.Resources = append(coreData.Resources,
				model.FlarumCreateTag(category))
		}
		for _, article := range pageInfo.Items {
			diss := model.FlarumCreateDiscussion(article)
			coreData.Resources = append(
				coreData.Resources,
				diss)
			coreData.APIDocument.Data = append(
				coreData.APIDocument.Data,
				diss)
		}

		for _, article := range pageInfo.Items {
			user := model.FlarumCreateUser(article)
			coreData.Resources = append(coreData.Resources, user)
			coreData.APIDocument.Included = append(
				coreData.APIDocument.Included,
				user,
			)
		}
		coreData.APIDocument.Links = make(map[string]string)
		coreData.APIDocument.Links["first"] = "https://flarum.yjzq.fun/api/v1/flarum/discussions?sort=&page%5Blimit%5D=20"

		coreData.Locales = make(map[string]string)
		coreData.Locales["en"] = "English"
		coreData.Locale = "en"
		coreData.Sessions = flarum.Session{
			UserID:    1,
			CsrfToken: "hello world",
		}

		evn.FlarumInfo = coreData
	}

	evn.SiteInfo = model.GetSiteInfo(redisDB)
	evn.PageInfo = pageInfo

	// 右侧的链接
	evn.Links = model.RedisLinkList(redisDB, false)

	h.Render(w, tpl, evn, "layout.html", "index.html")
}

// FlarumIndex flarum主页
func (h *BaseHandler) FlarumIndex(w http.ResponseWriter, r *http.Request) {

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
	for _, article := range pageInfo.Items {
		diss := model.FlarumCreateDiscussion(article)
		apiDoc.Data = append(
			apiDoc.Data,
			diss)
	}

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
