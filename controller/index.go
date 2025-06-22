package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/corvofeng/go-flarum/model"
	"github.com/corvofeng/go-flarum/model/flarum"

	"gorm.io/gorm"

	"github.com/go-redis/redis/v7"
)

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
	FT         filterType
	CID        uint64
	UID        uint64
	Page       uint64
	PageOffset uint64
	pageLimit  uint64
}

func createFlarumPageAPIDoc(
	reqctx *ReqContext,
	redisDB *redis.Client, gormDB *gorm.DB,
	appConf model.AppConf,
	df dissFilter,
	tz int,
) (flarum.CoreData, error) {
	var err error
	var topics []model.Topic
	var hasNext bool = false
	coreData := flarum.NewCoreData()

	apiDoc := &coreData.APIDocument // 注意, 获取到的是指针

	inAPI := reqctx.inAPI
	siteInfo := model.GetSiteInfo(redisDB)
	currentUser := reqctx.currentUser
	logger := reqctx.GetLogger()
	allUsers := make(map[uint64]bool)
	logger.Debugf("query with %+v", df)

	// 添加当前用户的session信息
	if currentUser != nil {
		user := model.FlarumCreateCurrentUser(*currentUser)
		allUsers[user.GetID()] = true
		coreData.AddCurrentUser(user)
		if !inAPI { // 做API请求时, 不更新csrf信息
			coreData.AddSessionData(user, currentUser.RefreshCSRF(redisDB))
		}
	}

	if df.FT == eCategory {
		topics, err = model.SQLGetTopicByTag(gormDB, redisDB, df.CID, df.PageOffset, df.pageLimit+1)
	} else if df.FT == eUserPost {
		topics, err = model.SQLGetTopicByUser(gormDB, df.UID, df.PageOffset, df.pageLimit+1)
	}

	if err != nil {
		logger.Error("get topics with error:", err)
		return coreData, err
	}

	categories, err := model.SQLGetTags(gormDB)
	// logger.Debugf("Get topics %+v", topics)

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
	// 添加当前页面的的帖子与用户信息, 已经去重
	for idx, topic := range topics {
		logger.Debugf("Get topic %d with tags: %+v title: %s", topic.ID, func() []string {
			var urlNames []string
			for _, tag := range topic.Tags {
				urlNames = append(urlNames, tag.URLName)
			}
			return urlNames
		}(), topic.Title)
		if idx == int(df.pageLimit) {
			hasNext = true
			continue
		}

		diss := model.FlarumCreateDiscussion(topic)
		res = append(res, diss)
		coreData.AppendResources(diss)
		getUser := func(uid uint64) {
			// 用户不存在则添加, 已经存在的用户不会考虑
			// TODO: 多次执行SQL可能会有性能问题
			if _, ok := allUsers[uid]; !ok {
				var u model.User
				result := gormDB.First(&u, uid)
				allUsers[uid] = true
				if errors.Is(result.Error, gorm.ErrRecordNotFound) {
					logger.Warningf("Can't get user %d error: not fount", uid)
				}
				user := model.FlarumCreateUser(u)
				coreData.AppendResources(user)
			}
		}

		getUser(topic.UserID)
		if topic.BlogMetaData.ID != 0 {
			logger.Debugf("Create blog meta for article: %+v", topic.BlogMetaData)
			apiDoc.AppendResources(model.FlarumCreateBlogMeta(topic.BlogMetaData, currentUser))
		}

		if topic.LastPostUserID == 0 {
			// logger.Warning("Can't get last post user id for", topic.ID)
		} else {
			getUser(topic.LastPostUserID)
		}
	}
	apiDoc.SetData(res)

	scf := appConf.Site
	apiDoc.Links["first"] = scf.MainDomain + model.FlarumAPIPath + "/discussions?sort=&page%5Blimit%5D=" + fmt.Sprintf("%d", df.pageLimit)
	if df.PageOffset != 0 {
		apiDoc.Links["prev"] = scf.MainDomain + model.FlarumAPIPath + "/discussions?sort=&page%5Boffset%5D=" + fmt.Sprintf("%d", df.PageOffset-df.pageLimit)
	}
	if hasNext {
		apiDoc.Links["next"] = scf.MainDomain + model.FlarumAPIPath + "/discussions?sort=&page%5Boffset%5D=" + fmt.Sprintf("%d", df.PageOffset+df.pageLimit)
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

	redisDB := h.App.RedisDB
	gormDB := h.App.GormDB
	logger := ctx.GetLogger()
	page := uint64(0)

	tpl := h.CurrentTpl(r)

	_tag, _ := h.safeGetParm(r, "tag")
	var tag model.Tag
	if _tag != "" {
		tag, err = model.SQLGetTagByUrlName(gormDB, _tag)
		if err != nil {
			logger.Errorf("Can't get tag by url name: %s, %s", _tag, err)
		}
	}

	df := dissFilter{
		FT:        eCategory,
		CID:       tag.ID,
		Page:      page,
		pageLimit: uint64(h.App.Cf.Site.HomeShowNum),
	}
	coreData, err := createFlarumPageAPIDoc(ctx, redisDB, h.App.GormDB, *h.App.Cf, df, scf.TimeZone)
	if err != nil {
		h.flarumErrorMsg(w, "无法获取帖子信息")
		return
	}

	var pageInfo model.ArticlePageInfo
	for _, item := range coreData.APIDocument.Included {
		if item.Type == "discussions" {
			ab := model.Topic{
				ID:    item.GetID(),
				Title: item.Attributes.(*flarum.Discussion).Title,
			}
			pageInfo.Items = append(pageInfo.Items, model.ArticleListItem{Topic: ab})
		}
	}
	evn := InitPageData(r)
	evn.FlarumInfo = coreData
	evn.PageInfo = pageInfo
	evn.PluginHTML["analytics"] = `<script async src="https://www.googletagmanager.com/gtag/js"></script>`

	h.Render(w, tpl, evn, "layout.html", "index.html")
}

// FlarumAPIDiscussions flarum文章api
func FlarumAPIDiscussions(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h
	scf := h.App.Cf.Site
	gormDB := h.App.GormDB

	redisDB := h.App.RedisDB
	var err error

	logger := h.App.Logger
	coreData := flarum.NewCoreData()
	// const pageLimit = 20
	// apiDoc := &coreData.APIDocument // 注意, 获取到的是指针

	// 需要返回的relations TODO: use it
	_include := r.FormValue("include")
	strings.Split(_include, ",")

	// 当前的排序方式 TODO: use it
	// _sort := r.FormValue("sort")
	// strings.Split(_sort, ",")

	// 当前的过滤方式 filter[tag]:  tag:r_funny
	_tag_filter := r.FormValue("filter[tag]")
	_author_filter := r.FormValue("filter[author]")
	if _tag_filter == "" {
		_tag_filter, _ = h.safeGetParm(r, "tag")
	}

	// 当前的偏移数目, 可得到页码数目, 页码从1开始
	_offset := r.FormValue("page[offset]")
	if _offset == "" {
		_offset = "0"
	}
	pageOffset, err := strconv.ParseUint(_offset, 10, 64)
	if err != nil {
		logger.Error("Parse offset err:", err)
		h.flarumErrorJsonify(w, createSimpleFlarumError("Can't get offset"+err.Error()))
		return
	}
	logger.Debugf("Get _filter: `%s`, page: `%d`", _tag_filter, pageOffset)
	pageLimit := uint64(h.App.Cf.Site.HomeShowNum)

	if _tag_filter == "" {
		df := dissFilter{
			FT:         eCategory,
			PageOffset: pageOffset,
			pageLimit:  pageLimit,
			CID:        0,
		}
		coreData, err = createFlarumPageAPIDoc(ctx, redisDB, h.App.GormDB, *h.App.Cf, df, scf.TimeZone)
	} else {
		if _tag_filter != "" {
			cate, err := model.SQLGetTagByUrlName(gormDB, _tag_filter)
			if err != nil {
				h.flarumErrorJsonify(w, createSimpleFlarumError("Can't create category"+err.Error()))
				return
			}
			df := dissFilter{
				FT:         eCategory,
				PageOffset: pageOffset,
				CID:        cate.ID,
				pageLimit:  pageLimit,
			}
			coreData, err = createFlarumPageAPIDoc(ctx, redisDB, h.App.GormDB, *h.App.Cf, df, scf.TimeZone)
			if err != nil {
				h.flarumErrorJsonify(w, createSimpleFlarumError("Can't create category"+err.Error()))
				return
			}
		} else if _author_filter != "" {
			user, err := model.SQLUserGetByName(h.App.GormDB, _author_filter)
			if err != nil {
				h.flarumErrorJsonify(w, createSimpleFlarumError("Can't create user"+err.Error()))
				return
			}
			df := dissFilter{
				FT:         eUserPost,
				PageOffset: pageOffset,
				UID:        user.ID,
				pageLimit:  pageLimit,
			}
			coreData, err = createFlarumPageAPIDoc(ctx, redisDB, h.App.GormDB, *h.App.Cf, df, scf.TimeZone)
			if err != nil {
				h.flarumErrorJsonify(w, createSimpleFlarumError("Can't create flarum page"+err.Error()))
				return
			}
		} else {
			// logger.Warning("Can't use filter:", _filter)
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
