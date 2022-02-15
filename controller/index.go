package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"zoe/model"
	"zoe/model/flarum"

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
	FT    filterType
	CID   uint64
	UID   uint64
	Page  uint64
	Limit uint64
}

func createFlarumPageAPIDoc(
	reqctx *ReqContext,
	sqlDB *sql.DB, redisDB *redis.Client, gormDB *gorm.DB,
	appConf model.AppConf,
	df dissFilter,
	tz int,
) (flarum.CoreData, error) {
	var err error
	var articlePageInfo model.ArticlePageInfo
	var topics []model.Topic
	coreData := flarum.NewCoreData()

	apiDoc := &coreData.APIDocument // 注意, 获取到的是指针
	page := df.Page

	inAPI := reqctx.inAPI
	siteInfo := model.GetSiteInfo(redisDB)
	currentUser := reqctx.currentUser
	logger := reqctx.GetLogger()
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

	if df.FT == eCategory {
		topics = model.SQLTopicGetByTag(gormDB, sqlDB, redisDB, df.CID, page, df.Limit, tz)
	} else if df.FT == eUserPost {
		articlePageInfo = model.SQLTopicGetByUID(gormDB, sqlDB, redisDB, df.UID, page, df.Limit, tz)
		// topics = articlePageInfo.Items
	}

	categories, err := model.SQLGetTags(gormDB)

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
	for _, topic := range topics {
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

		if topic.LastPostUserID == 0 {
			logger.Warning("Can't get last post uer id for", topic.ID)
		} else {
			getUser(topic.LastPostUserID)
		}
	}
	apiDoc.SetData(res)

	scf := appConf.Site
	apiDoc.Links["first"] = scf.MainDomain + model.FlarumAPIPath + "/discussions?sort=&page%5Blimit%5D=" + fmt.Sprintf("%d", df.Limit)
	if page != 1 {
		apiDoc.Links["prev"] = scf.MainDomain + model.FlarumAPIPath + "/discussions?sort=&page%5Boffset%5D=" + fmt.Sprintf("%d", page*df.Limit)
	}

	if articlePageInfo.HasNext {
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
	gormDB := h.App.GormDB
	logger := ctx.GetLogger()
	page := uint64(1)

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
		FT:    eCategory,
		CID:   tag.ID,
		Page:  page,
		Limit: 10,
	}
	coreData, err := createFlarumPageAPIDoc(ctx, sqlDB, redisDB, h.App.GormDB, *h.App.Cf, df, scf.TimeZone)
	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("无法获取帖子信息"))
		return
	}
	// 设置语言信息

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
	evn.PluginHTML["analytics"] = `<script async src="https://www.googletagmanager.com/gtag/js?id=%%TRACKING_CODE%%"></script>`

	h.Render(w, tpl, evn, "layout.html", "index.html")
}

// FlarumAPIDiscussions flarum文章api
func FlarumAPIDiscussions(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h
	scf := h.App.Cf.Site
	gormDB := h.App.GormDB
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

	// 当前的过滤方式 filter[tag]:  tag:r_funny
	_tag_filter := r.FormValue("filter[tag]")
	_author_filter := r.FormValue("filter[author]")

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
	logger.Debugf("Get _filter: `%s`, page: `%d`", _tag_filter, page)

	if _tag_filter == "" {
		df := dissFilter{
			FT:    eCategory,
			Page:  page,
			Limit: pageLimit,
			CID:   0,
		}
		coreData, err = createFlarumPageAPIDoc(ctx, sqlDB, redisDB, h.App.GormDB, *h.App.Cf, df, scf.TimeZone)
	} else {
		// data := strings.Trim(_filter, " ")
		if _tag_filter != "" {
			cate, err := model.SQLGetTagByUrlName(gormDB, _tag_filter)
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
			coreData, err = createFlarumPageAPIDoc(ctx, sqlDB, redisDB, h.App.GormDB, *h.App.Cf, df, scf.TimeZone)
			// } else if strings.HasPrefix(data, "author:") {
			// 	user, err := model.SQLUserGetByName(h.App.GormDB, data[7:])
			// 	if err != nil {
			// 		h.flarumErrorJsonify(w, createSimpleFlarumError("Can't create user"+err.Error()))
			// 		return
			// 	}
			// 	df := dissFilter{
			// 		FT:    eUserPost,
			// 		Page:  page,
			// 		UID:   user.ID,
			// 		Limit: pageLimit,
			// 	}
			// 	coreData, err = createFlarumPageAPIDoc(ctx, sqlDB, redisDB, h.App.GormDB, *h.App.Cf, df, scf.TimeZone)

		} else if _author_filter != "" {
			user, err := model.SQLUserGetByName(h.App.GormDB, _author_filter)
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
			coreData, err = createFlarumPageAPIDoc(ctx, sqlDB, redisDB, h.App.GormDB, *h.App.Cf, df, scf.TimeZone)
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
