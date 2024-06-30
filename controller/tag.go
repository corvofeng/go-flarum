package controller

import (
	"net/http"

	"github.com/corvofeng/go-flarum/model"
	"github.com/corvofeng/go-flarum/model/flarum"

	"github.com/go-redis/redis/v7"
	"gorm.io/gorm"
)

func createFlarumTagAPIDoc(
	reqctx *ReqContext,
	gormDB *gorm.DB, redisDB *redis.Client,
	appConf model.AppConf,
	tz int,
) (flarum.CoreData, error) {
	var err error
	coreData := flarum.NewCoreData()
	currentUser := reqctx.currentUser
	inAPI := reqctx.inAPI
	siteInfo := model.GetSiteInfo(redisDB)
	apiDoc := &coreData.APIDocument // 注意, 获取到的是指针
	logger := reqctx.GetLogger()

	// 添加当前用户的session信息
	if currentUser != nil {
		user := model.FlarumCreateCurrentUser(*currentUser)
		coreData.AddCurrentUser(user)
		if !inAPI { // 做API请求时, 不更新csrf信息
			coreData.AddSessionData(user, currentUser.RefreshCSRF(redisDB))
		}
	}

	categories, err := model.SQLGetTags(gormDB)
	if err != nil {
		logger.Info("Can't get categories")
	}

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
	for _, category := range categories {
		tag := model.FlarumCreateTag(category)
		res = append(res, tag)
	}
	// article, err := model.SQLArticleGetByID(gormDB,    redisDB, 1)
	// if err != nil {
	// 	logger.Info("Can't get article", err.Error())
	// }
	// diss := model.FlarumCreateDiscussion(article)
	// apiDoc.AppendResources(diss)

	apiDoc.SetData(res)
	model.FlarumCreateLocale(&coreData, reqctx.locale)

	return coreData, err

}

// FlarumIndex flarum主页
func FlarumTagAll(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h
	inAPI := ctx.inAPI
	scf := h.App.Cf.Site

	redisDB := h.App.RedisDB
	logger := ctx.GetLogger()

	coreData, err := createFlarumTagAPIDoc(
		ctx, h.App.GormDB, redisDB, *h.App.Cf, scf.TimeZone)

	if err != nil {
		h.flarumErrorMsg(w, "查询标签信息错误:"+err.Error())
	}

	logger.Info(h.safeGetParm(r, "tag"))

	// 如果是API直接进行返回
	if inAPI {
		h.jsonify(w, coreData.APIDocument)
		return
	}

	tpl := h.CurrentTpl(r)
	evn := InitPageData(r)
	evn.FlarumInfo = coreData
	h.Render(w, tpl, evn, "layout.html", "article.html")
}
