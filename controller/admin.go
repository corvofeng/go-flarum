package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/corvofeng/go-flarum/model"
	"github.com/corvofeng/go-flarum/model/flarum"

	"github.com/go-redis/redis/v7"
	"gorm.io/gorm"
)

func createFlarumAdminAPIDoc(
	reqctx *ReqContext,
	gormDB *gorm.DB,
	redisDB *redis.Client,
	appConf model.AppConf,
	siteInfo model.SiteInfo,
	tz int,
) (flarum.AdminCoreData, error) {
	var err error
	coreData := flarum.NewAdminCoreData()
	inAPI := reqctx.inAPI
	currentUser := reqctx.currentUser
	logger := reqctx.GetLogger()

	// 所有分类的信息, 用于整个站点的信息
	var flarumTags []flarum.Resource

	if currentUser != nil {
		user := model.FlarumCreateCurrentUser(*currentUser)
		coreData.AddCurrentUser(user)
		if !inAPI { // 做API请求时, 不更新csrf信息
			coreData.AddSessionData(user, currentUser.RefreshCSRF(redisDB))
		}
	}
	// 添加当前站点信息
	categories, err := model.SQLGetTags(gormDB)
	if err != nil {
		logger.Error("Get all categories error", err)
	}
	for _, category := range categories {
		flarumTags = append(flarumTags, model.FlarumCreateTag(category))
	}
	coreData.AppendResources(model.FlarumCreateForumInfo(
		currentUser,
		appConf, siteInfo, flarumTags,
	))
	model.FlarumCreateLocale(&coreData.CoreData, reqctx.locale)
	coreData.Extensions, err = flarum.ReadExtensionMetadata(appConf.Main.ExtensionsDir)
	if err != nil {
		return coreData, err
	}
	// 需要给出
	coreData.Settings.DefaultRoute = "/all"
	coreData.Settings.FlarumMarkdownMdarea = "1"
	coreData.Settings.FlarumMentionsAllowUsernameFormat = "1"
	coreData.Settings.FlarumTagsMaxPrimaryTags = "3"
	coreData.Settings.FlarumTagsMaxSecondaryTags = "3"
	//     "flarum-tags.min_primary_tags": "1",
	//     "flarum-tags.min_secondary_tags": "0",

	// coreData.Settings.ExtensionsEnabled = "[\"flarum-flags\",\"flarum-mentions\",\"flarum-bbcode\",\"flarum-markdown\",\"flarum-approval\",\"flarum-statistics\",\"flarum-sticky\",\"flarum-emoji\",\"flarum-tags\",\"flarum-suspend\",\"flarum-subscriptions\",\"flarum-lock\",\"flarum-likes\",\"flarum-lang-english\"]"
	enabledExtensions := []string{
		"flarum-flags",
		"flarum-mentions",
		"flarum-bbcode",
		"flarum-markdown",
		"flarum-approval",
		"flarum-statistics",
		"flarum-sticky",
		"flarum-emoji",
		"flarum-tags",
		"flarum-suspend",
		"flarum-subscriptions",
		"flarum-lock",
		"flarum-likes",
		"flarum-lang-english",
		"v17development-flarum-blog",
	}

	enabledExtensionsRaw, _ := json.Marshal(&enabledExtensions)
	coreData.Settings.ExtensionsEnabled = string(enabledExtensionsRaw)

	coreData.PhpVersion = "8.0.6"
	coreData.MysqlVersion = "10.4.8-MariaDB-1:10.4.8+maria~bionic"
	// coreData.SlugDrivers = map[string]interface{}{
	// 	"Flarum\\Discussion\\Discussion": []string{
	// 		"default",
	// 	},
	// 	"Flarum\\User\\User": []string{
	// 		"default",
	// 		"id",
	// 	},
	// }

	coreData.DisplayNameDrivers = []string{"username"}
	coreData.ModelStatistics.Users.Total = 3

	return coreData, err
}

// ArticleHomeList 文章主页
func AdminHome(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h
	tpl := h.CurrentTpl(r)
	evn := InitPageData(r)
	scf := h.App.Cf.Site

	redisDB := h.App.RedisDB
	gormDB := h.App.GormDB
	// logger := ctx.GetLogger()
	fmt.Println(h.App.Cf.Main.ExtensionsDir)

	coreData, err := createFlarumAdminAPIDoc(
		ctx, gormDB, redisDB, *h.App.Cf, model.GetSiteInfo(redisDB), scf.TimeZone)
	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("Get api doc error"+err.Error()))
		return
	}
	evn.FlarumInfo = coreData

	h.Render(w, tpl, evn, "layout.html", "admin.html")
}
