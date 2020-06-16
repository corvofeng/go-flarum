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
			h.jsonify(w, rsp)
			return
		}
	}
	if len(score) > 0 {
		_, err = strconv.ParseUint(score, 10, 64)
		if err != nil {
			rsp = response{400, "scope type err"}
			h.jsonify(w, rsp)
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
		h.jsonify(w, rsp)
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
func FlarumIndex(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := GetRetContext(r)
	h := ctx.h
	scf := h.App.Cf.Site
	sqlDB := h.App.MySQLdb
	redisDB := h.App.RedisDB
	page := uint64(1)

	// 获取贴子列表
	pageInfo := model.SQLCIDArticleListByPage(sqlDB, redisDB, 0, page, uint64(scf.HomeShowNum), scf.TimeZone)
	// categories, err := model.SQLGetAllCategory(sqlDB)

	tpl := h.CurrentTpl(r)
	evn := &pageData{}
	evn.SiteCf = scf
	coreData := flarum.CoreData{}

	categories, err := model.SQLGetNotEmptyCategory(sqlDB, redisDB)
	// 添加所有分类的信息
	var flarumTags []flarum.Resource
	for _, category := range categories {
		tag := model.FlarumCreateTag(category)
		coreData.AppendResourcs(tag)
		flarumTags = append(flarumTags, tag)
	}

	// 添加主站点信息
	coreData.AppendResourcs(model.FlarumCreateForumInfo(*h.App.Cf, evn.SiteInfo, flarumTags))

	// 添加当前页面的帖子信息
	var res []flarum.Resource
	for _, article := range pageInfo.Items {
		lastComent := model.SQLGetCommentByID(sqlDB, redisDB, article.LastPostID, scf.TimeZone)
		diss := model.FlarumCreateDiscussion(article, lastComent)
		coreData.AppendResourcs(diss)
		res = append(res, diss)
	}
	coreData.APIDocument.SetData(res)

	// 添加当前页面的帖子的用户信息, TODO: 用户有可能重复, 这里理论是需要优化的
	for _, article := range pageInfo.Items {
		user := model.FlarumCreateUser(article)
		coreData.AppendResourcs(user)
	}

	// 设置语言信息
	coreData.Locales = make(map[string]string)

	coreData.Locales["en"] = "English"
	coreData.Locales["zh"] = "中文"
	coreData.Locale = "en"

	coreData.APIDocument.Links = make(map[string]string)
	coreData.APIDocument.Links["first"] = scf.MainDomain + model.FlarumAPIPath + "/discussions?sort=&page%5Blimit%5D=20"
	coreData.APIDocument.Links["next"] = scf.MainDomain + model.FlarumAPIPath + "/discussions?sort=&page%5Boffset%5D20"

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
func FlarumAPIDiscussions(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h
	scf := h.App.Cf.Site
	sqlDB := h.App.MySQLdb
	redisDB := h.App.RedisDB
	var page uint64
	var pageInfo model.ArticlePageInfo

	logger := h.App.Logger

	apiDoc := flarum.NewAPIDoc()

	// 需要返回的relations TODO: use it
	_include := r.FormValue("include")
	strings.Split(_include, ",")
	// fmt.Println(_include)

	// 当前的排序方式 TODO: use it
	// _sort := r.FormValue("sort")
	// strings.Split(_sort, ",")
	// fmt.Println(_sort)

	// 当前的过滤方式 filter[q]:  tag:r_funny
	_filter := r.FormValue("filter[q]")
	// fmt.Println(_filter)

	// 当前的偏移数目, 可得到页码数目, 页码从1开始
	_offset := r.FormValue("page[offset]")
	if _offset != "" {
		data, err := strconv.ParseUint(_offset, 10, 64)
		if err != nil {
			logger.Error("Parse offset err:", err)
			h.jsonify(w, apiDoc)
			return
		}
		page = data / 20
	}
	page = page + 1

	if _filter == "" {
		pageInfo = model.SQLCIDArticleListByPage(sqlDB, redisDB, 0, page, uint64(scf.HomeShowNum), scf.TimeZone)
	} else {
		data := strings.Trim(_filter, " ")
		if strings.HasPrefix(data, "tag:") {
			cate, err := model.SQLCategoryGetByURLName(sqlDB, data[4:])
			if err != nil {
				logger.Error("Can't get category", err)
				h.jsonify(w, apiDoc)
				return
			}
			pageInfo = model.SQLCIDArticleListByPage(sqlDB, redisDB, cate.ID, page, uint64(scf.HomeShowNum), scf.TimeZone)
		}
	}

	var dissArr []flarum.Resource
	for _, article := range pageInfo.Items {
		var lastComent model.Comment
		if article.LastPostID != 0 {
			lastComent = model.SQLGetCommentByID(sqlDB, redisDB, article.LastPostID, scf.TimeZone)
		} else {
			lastComent = model.Comment{}
		}
		diss := model.FlarumCreateDiscussion(article, lastComent)
		dissArr = append(dissArr, diss)
	}
	apiDoc.SetData(dissArr)

	for _, article := range pageInfo.Items {
		user := model.FlarumCreateUser(article)
		apiDoc.AppendResourcs(user)
	}
	// categories, err := model.SQLGetAllCategory(sqlDB)
	categories, err := model.SQLGetNotEmptyCategory(sqlDB, redisDB)

	// 添加所有分类的信息
	for _, category := range categories {
		apiDoc.AppendResourcs(model.FlarumCreateTag(category))
	}

	// 添加当前用户的session信息
	currentUser, err := h.CurrentUser(w, r)
	if err == nil {
		user := model.FlarumCreateCurrentUser(currentUser)
		apiDoc.AppendResourcs(user)
	}

	apiDoc.Links["first"] = scf.MainDomain + model.FlarumAPIPath + "/discussions?sort=&page%5Blimit%5D=20"
	if page != 1 {
		apiDoc.Links["prev"] = scf.MainDomain + model.FlarumAPIPath + "/discussions?sort=&page%5Boffset%5D=" + fmt.Sprintf("%d", page*20)
	}

	if pageInfo.HasNext {
		apiDoc.Links["next"] = scf.MainDomain + model.FlarumAPIPath + "/discussions?sort=&page%5Boffset%5D=" + fmt.Sprintf("%d", (page+1)*20)
	}

	h.jsonify(w, apiDoc)
}
