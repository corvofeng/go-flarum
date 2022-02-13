package controller

import (
	"database/sql"
	"net/http"
	"zoe/model"
	"zoe/model/flarum"

	"github.com/go-redis/redis/v7"
	"goji.io/pat"
	"gorm.io/gorm"
)

func (h *BaseHandler) TagDetail(w http.ResponseWriter, r *http.Request) {
	// btn, key := r.FormValue("btn"), r.FormValue("key")
	// if len(key) > 0 {
	// 	_, err := strconv.ParseUint(key, 10, 64)
	// 	if err != nil {
	// 		w.Write([]byte(`{"retcode":400,"retmsg":"key type err"}`))
	// 		return
	// 	}
	// }

	tag := pat.Param(r, "tag")
	// tagLow := strings.ToLower(tag)

	// cmd := "hrscan"
	// if btn == "prev" {
	// 	cmd = "hscan"
	// }

	// db := h.App.Db
	redisDB := h.App.RedisDB
	scf := h.App.Cf.Site
	// rs := db.Hscan("tag:"+tagLow, nil, 1)
	// if rs.State != "ok" {
	// 	w.Write([]byte(`{"retcode":404,"retmsg":"not found"}`))
	// 	return
	// }

	currentUser, _ := h.CurrentUser(w, r)

	// pageInfo := model.UserArticleList(db, cmd, "tag:"+tagLow, key, scf.PageShowNum, scf.TimeZone)
	pageInfo := model.ArticlePageInfo{}

	type tagDetail struct {
		Name   string
		Number uint64
	}
	// TODO: Delete this
	type pageData struct {
		BasePageData
		Tag      tagDetail
		PageInfo model.ArticlePageInfo
	}

	tpl := h.CurrentTpl(r)

	evn := &pageData{}
	evn.SiteCf = scf
	evn.Title = tag + " - " + scf.Name
	evn.Keywords = tag
	evn.Description = tag
	evn.IsMobile = tpl == "mobile"

	evn.CurrentUser = currentUser
	evn.ShowSideAd = true
	evn.PageName = "category_detail"
	// evn.HotNodes = model.CategoryHot(db, scf.CategoryShowNum)
	// evn.NewestNodes = model.CategoryNewest(db, scf.CategoryShowNum)

	evn.Tag = tagDetail{
		Name: tag,
	}
	evn.PageInfo = pageInfo
	evn.SiteInfo = model.GetSiteInfo(redisDB)

	h.Render(w, tpl, evn, "layout.html", "tag.html")
}

func createFlarumTagAPIDoc(
	reqctx *ReqContext,
	gormDB *gorm.DB, sqlDB *sql.DB, redisDB *redis.Client,
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
	// article, err := model.SQLArticleGetByID(gormDB, sqlDB, redisDB, 1)
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
	sqlDB := h.App.MySQLdb
	redisDB := h.App.RedisDB
	logger := ctx.GetLogger()

	coreData, err := createFlarumTagAPIDoc(
		ctx, h.App.GormDB, sqlDB, redisDB, *h.App.Cf, scf.TimeZone)

	if err != nil {
		h.flarumErrorMsg(w, "查询标签信息错误:"+err.Error())
	}

	// 如果是API直接进行返回
	if inAPI {
		h.jsonify(w, coreData.APIDocument)
		return
	}
	logger.Info(h.safeGetParm(r, "tag"))

	tpl := h.CurrentTpl(r)
	evn := InitPageData(r)
	evn.FlarumInfo = coreData
	h.Render(w, tpl, evn, "layout.html", "article.html")
}
