package controller

import (
	"database/sql"
	"fmt"
	"goyoubbs/model"
	"goyoubbs/model/flarum"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v7"
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

	// type pageData struct {
	// 	BasePageData
	// 	SiteInfo   model.SiteInfo
	// 	PageInfo   model.ArticlePageInfo
	// 	Links      []model.Link
	// 	FlarumInfo interface{}
	// }

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
	categories, err := model.SQLGetAllCategory(sqlDB, redisDB)

	tpl := h.CurrentTpl(r)
	evn := InitPageData(r)
	evn.Keywords = evn.Title
	evn.IsMobile = tpl == "mobile"
	evn.ShowSideAd = false
	evn.PageName = "home"
	evn.NewestNodes = categories
	// evn.HotNodes = model.CategoryHot(db, scf.CategoryShowNum)
	// evn.NewestNodes = model.CategoryNewest(db, scf.CategoryShowNum)
	evn.PageInfo = pageInfo

	h.Render(w, tpl, evn, "layout.html", "index.html")
}

// 记录当前的过滤器内容
type filterType string

const (
	eUserPost filterType = "userpost"
	ePost     filterType = "post"
	ePosts    filterType = "posts"
	eCategory filterType = "category"
	eArticle  filterType = "article"
	eReply    filterType = "reply"
)

type dissFilter struct {
	FT    filterType
	CID   uint64
	UID   uint64
	Page  uint64
	Limit uint64
}

func createFlarumPageAPIDoc(
	reqctx *ReqContext,
	sqlDB *sql.DB, redisDB *redis.Client,
	appConf model.AppConf, siteInfo model.SiteInfo,
	df dissFilter,
	tz int,
) (flarum.CoreData, error) {
	var err error
	var pageInfo model.ArticlePageInfo
	coreData := flarum.NewCoreData()
	apiDoc := &coreData.APIDocument // 注意, 获取到的是指针
	page := df.Page

	inAPI := reqctx.inAPI
	currentUser := reqctx.currentUser
	logger := reqctx.GetLogger()

	if df.FT == eCategory {
		pageInfo = model.SQLArticleGetByCID(sqlDB, redisDB, df.CID, page, df.Limit, tz)
	} else if df.FT == eUserPost {
		pageInfo = model.SQLArticleGetByUID(sqlDB, redisDB, df.UID, page, df.Limit, tz)
	}

	categories, err := model.SQLGetNotEmptyCategory(sqlDB, redisDB)

	// 添加所有分类的信息
	var flarumTags []flarum.Resource
	for _, category := range categories {
		tag := model.FlarumCreateTag(category)
		coreData.AppendResources(tag)
		flarumTags = append(flarumTags, tag)
	}

	// 添加主站点信息
	coreData.AppendResources(model.FlarumCreateForumInfo(
		currentUser,
		appConf, siteInfo, flarumTags,
	))

	var res []flarum.Resource
	allUsers := make(map[uint64]bool)
	// 添加当前用户的session信息
	if currentUser != nil {
		user := model.FlarumCreateCurrentUser(*currentUser)
		allUsers[user.GetID()] = true
		coreData.AddCurrentUser(user)
		if !inAPI { // 做API请求时, 不更新csrf信息
			coreData.AddSessionData(user, currentUser.RefreshCSRF(redisDB))
		}
	}

	// 添加当前页面的的帖子与用户信息, 已经去重
	for _, article := range pageInfo.Items {
		diss := model.FlarumCreateDiscussion(article)
		res = append(res, diss)
		coreData.AppendResources(diss)
		getUser := func(uid uint64) {
			// 用户不存在则添加, 已经存在的用户不会考虑
			// TODO: 多次执行SQL可能会有性能问题
			if _, ok := allUsers[uid]; !ok {
				u, err := model.SQLUserGetByID(sqlDB, uid)
				if err != nil {
					logger.Warningf("Get user %d error: %s", uid, err)
				} else {
					user := model.FlarumCreateUser(u)
					allUsers[uid] = true
					coreData.AppendResources(user)
				}
			}
		}
		getUser(article.UID)
		d := diss.Attributes.(*flarum.Discussion)
		if d.LastPostID != 0 {
			getUser(d.LastUserID)
		}
	}
	apiDoc.SetData(res)

	scf := appConf.Site
	apiDoc.Links["first"] = scf.MainDomain + model.FlarumAPIPath + "/discussions?sort=&page%5Blimit%5D=" + fmt.Sprintf("%d", df.Limit)
	if page != 1 {
		apiDoc.Links["prev"] = scf.MainDomain + model.FlarumAPIPath + "/discussions?sort=&page%5Boffset%5D=" + fmt.Sprintf("%d", page*df.Limit)
	}

	if pageInfo.HasNext {
		apiDoc.Links["next"] = scf.MainDomain + model.FlarumAPIPath + "/discussions?sort=&page%5Boffset%5D=" + fmt.Sprintf("%d", (page+1)*20)
	}
	model.FlarumCreateLocale(&coreData, reqctx.locale)

	return coreData, err
}

// FlarumIndex flarum主页
func FlarumIndex(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := GetRetContext(r)
	h := ctx.h
	scf := h.App.Cf.Site
	sqlDB := h.App.MySQLdb
	redisDB := h.App.RedisDB
	// logger := ctx.GetLogger()
	page := uint64(1)

	tpl := h.CurrentTpl(r)
	si := model.GetSiteInfo(redisDB)

	df := dissFilter{
		FT:    eCategory,
		CID:   0,
		Page:  page,
		Limit: 10,
	}
	coreData, err := createFlarumPageAPIDoc(ctx, sqlDB, redisDB, *h.App.Cf, si, df, scf.TimeZone)
	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("无法获取帖子信息"))
		return
	}
	// 设置语言信息

	var pageInfo model.ArticlePageInfo
	for _, item := range coreData.APIDocument.Included {
		if item.Type == "discussions" {
			ab := model.ArticleBase{
				ID:    item.GetID(),
				Title: item.Attributes.(*flarum.Discussion).Title,
			}
			pageInfo.Items = append(pageInfo.Items, model.ArticleListItem{ArticleBase: ab})
		}
	}
	evn := InitPageData(r)
	evn.FlarumInfo = coreData
	evn.PageInfo = pageInfo
	evn.PluginHTML["analytics"] = `<script async src="https://www.googletagmanager.com/gtag/js?id=%%TRACKING_CODE%%"></script>`

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
	var err error

	logger := h.App.Logger
	coreData := flarum.NewCoreData()
	const pageLimit = 10
	// apiDoc := &coreData.APIDocument // 注意, 获取到的是指针

	// 需要返回的relations TODO: use it
	_include := r.FormValue("include")
	strings.Split(_include, ",")

	// 当前的排序方式 TODO: use it
	// _sort := r.FormValue("sort")
	// strings.Split(_sort, ",")

	// 当前的过滤方式 filter[q]:  tag:r_funny
	_filter := r.FormValue("filter[q]")

	// 当前的偏移数目, 可得到页码数目, 页码从1开始
	_offset := r.FormValue("page[offset]")
	if _offset != "" {
		data, err := strconv.ParseUint(_offset, 10, 64)
		if err != nil {
			logger.Error("Parse offset err:", err)
			h.flarumErrorJsonify(w, createSimpleFlarumError("Can't get offset"+err.Error()))
			return
		}
		page = data / pageLimit
	}
	page = page + 1
	si := model.GetSiteInfo(redisDB)
	logger.Debugf("Get _filter: `%s`, page: `%d`", _filter, page)

	if _filter == "" {
		df := dissFilter{
			FT:    eCategory,
			Page:  page,
			Limit: pageLimit,
			CID:   0,
		}
		coreData, err = createFlarumPageAPIDoc(ctx, sqlDB, redisDB, *h.App.Cf, si, df, scf.TimeZone)
	} else {
		data := strings.Trim(_filter, " ")
		if strings.HasPrefix(data, "tag:") {
			cate, err := model.SQLCategoryGetByURLName(sqlDB, data[4:])
			if err != nil {
				h.flarumErrorJsonify(w, createSimpleFlarumError("Can't create category"+err.Error()))
				return
			}
			df := dissFilter{
				FT:    eCategory,
				Page:  page,
				CID:   cate.ID,
				Limit: pageLimit,
			}
			coreData, err = createFlarumPageAPIDoc(ctx, sqlDB, redisDB, *h.App.Cf, si, df, scf.TimeZone)
		} else if strings.HasPrefix(data, "author:") {
			user, err := model.SQLUserGetByName(sqlDB, data[7:])
			if err != nil {
				h.flarumErrorJsonify(w, createSimpleFlarumError("Can't create user"+err.Error()))
				return
			}
			df := dissFilter{
				FT:    eUserPost,
				Page:  page,
				UID:   user.ID,
				Limit: pageLimit,
			}
			coreData, err = createFlarumPageAPIDoc(ctx, sqlDB, redisDB, *h.App.Cf, si, df, scf.TimeZone)
		} else {
			logger.Warning("Can't use filter:", _filter)
			h.flarumErrorJsonify(w, createSimpleFlarumError("过滤器未实现"))
			return
		}
	}

	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("无法获取帖子信息"))
		return
	}

	h.jsonify(w, coreData.APIDocument)
}
