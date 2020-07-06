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
	"github.com/op/go-logging"
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

// 记录当前的过滤器内容
type filterType string

const (
	eUserPost filterType = "userpost"
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
	logger *logging.Logger,
	sqlDB *sql.DB, redisDB *redis.Client,
	appConf model.AppConf, siteInfo model.SiteInfo,
	currentUser *model.User,
	inAPI bool,
	// page uint64,
	// cid uint64,
	df dissFilter,
	tz int,
) (flarum.CoreData, error) {
	var err error
	var pageInfo model.ArticlePageInfo
	coreData := flarum.NewCoreData()
	apiDoc := &coreData.APIDocument // 注意, 获取到的是指针
	page := df.Page

	if df.FT == eCategory {
		pageInfo = model.SQLArticleGetByCID(sqlDB, redisDB, df.CID, page, 20, tz)
	} else if df.FT == eUserPost {
		pageInfo = model.SQLArticleGetByUID(sqlDB, redisDB, df.UID, page, 20, tz)
	}

	categories, err := model.SQLGetNotEmptyCategory(sqlDB, redisDB)

	// 添加所有分类的信息
	var flarumTags []flarum.Resource
	for _, category := range categories {
		tag := model.FlarumCreateTag(category)
		coreData.AppendResourcs(tag)
		flarumTags = append(flarumTags, tag)
	}

	// 添加主站点信息
	coreData.AppendResourcs(model.FlarumCreateForumInfo(appConf, siteInfo, flarumTags))

	var res []flarum.Resource
	allUsers := make(map[uint64]bool)

	// 添加当前页面的的帖子与用户信息, 已经去重
	for _, article := range pageInfo.Items {
		lastComent, err := model.SQLGetCommentByID(sqlDB, redisDB, article.LastPostID, tz)
		if err != nil {
			logger.Warningf("Can't get article comment(%d, %d) err: %s", article.ID, article.LastPostID, err.Error())
		}
		diss := model.FlarumCreateDiscussion(article, lastComent)
		res = append(res, diss)
		coreData.AppendResourcs(diss)

		if currentUser != nil && article.UID == currentUser.ID { // 当前用户单独进行添加
			continue
		}

		// 用户不存在则添加, 已经存在的用户不会考虑
		if _, ok := allUsers[article.UID]; !ok {
			user := model.FlarumCreateUser(article)
			// apiDoc.AppendResourcs(user)
			allUsers[article.UID] = true
			coreData.AppendResourcs(user)
		}
	}
	apiDoc.SetData(res)

	// 添加当前用户的session信息
	if currentUser != nil {
		user := model.FlarumCreateCurrentUser(*currentUser)
		coreData.AddCurrentUser(user)
		if !inAPI { // 做API请求时, 不更新csrf信息
			coreData.AddSessionData(user, currentUser.RefreshCSRF(redisDB))
		}
	}

	scf := appConf.Site
	apiDoc.Links["first"] = scf.MainDomain + model.FlarumAPIPath + "/discussions?sort=&page%5Blimit%5D=20"
	if page != 1 {
		apiDoc.Links["prev"] = scf.MainDomain + model.FlarumAPIPath + "/discussions?sort=&page%5Boffset%5D=" + fmt.Sprintf("%d", page*20)
	}

	if pageInfo.HasNext {
		apiDoc.Links["next"] = scf.MainDomain + model.FlarumAPIPath + "/discussions?sort=&page%5Boffset%5D=" + fmt.Sprintf("%d", (page+1)*20)
	}

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
	logger := ctx.GetLogger()
	page := uint64(1)

	tpl := h.CurrentTpl(r)
	evn := &pageData{}
	evn.SiteCf = scf
	si := model.GetSiteInfo(redisDB)

	df := dissFilter{
		FT:   eCategory,
		CID:  0,
		Page: page,
	}
	coreData, err := createFlarumPageAPIDoc(logger, sqlDB, redisDB, *h.App.Cf, si, ctx.currentUser, ctx.inAPI, df, scf.TimeZone)
	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("无法获取帖子信息"))
		return
	}

	// 设置语言信息
	coreData.Locales = make(map[string]string)
	coreData.Locales["en"] = "English"
	coreData.Locales["zh"] = "中文"
	coreData.Locale = "en"

	evn.FlarumInfo = coreData

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
	apiDoc := &coreData.APIDocument // 注意, 获取到的是指针

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
			h.jsonify(w, apiDoc)
			return
		}
		page = data / 20
	}
	page = page + 1
	si := model.GetSiteInfo(redisDB)

	if _filter == "" {
		df := dissFilter{
			FT:   eCategory,
			Page: page,
			CID:  0,
		}
		coreData, err = createFlarumPageAPIDoc(logger, sqlDB, redisDB, *h.App.Cf, si, ctx.currentUser, ctx.inAPI, df, scf.TimeZone)
	} else {
		data := strings.Trim(_filter, " ")
		if strings.HasPrefix(data, "tag:") {
			cate, err := model.SQLCategoryGetByURLName(sqlDB, data[4:])
			if err != nil {
				h.flarumErrorJsonify(w, createSimpleFlarumError("Can't create category"+err.Error()))
				return
			}
			df := dissFilter{
				FT:   eCategory,
				Page: page,
				CID:  cate.ID,
			}
			coreData, err = createFlarumPageAPIDoc(logger, sqlDB, redisDB, *h.App.Cf, si, ctx.currentUser, ctx.inAPI, df, scf.TimeZone)
		} else if strings.HasPrefix(data, "author:") {
			user, err := model.SQLUserGetByName(sqlDB, data[7:])
			if err != nil {
				h.flarumErrorJsonify(w, createSimpleFlarumError("Can't create user"+err.Error()))
				return
			}
			df := dissFilter{
				FT:   eUserPost,
				Page: page,
				UID:  user.ID,
			}
			coreData, err = createFlarumPageAPIDoc(logger, sqlDB, redisDB, *h.App.Cf, si, ctx.currentUser, ctx.inAPI, df, scf.TimeZone)
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
