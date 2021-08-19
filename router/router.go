package router

import (
	"net/http"

	ct "zoe/controller"
	"zoe/model"
	"zoe/system"

	"goji.io"
	"goji.io/pat"
)

// NewRouter create the router
func NewRouter(app *system.Application) *goji.Mux {
	sp := goji.SubMux()
	sp.Use(ct.TrackerMiddleware)

	if app.IsFlarum() {
		NewFlarumRouter(app, sp)
	} else {
		NewzoeRouter(app, sp)
	}
	return sp
}

// NewzoeRouter zoe的router
func NewzoeRouter(app *system.Application, sp *goji.Mux) *goji.Mux {
	h := ct.BaseHandler{App: app}
	sp.Use(h.InitMiddlewareContext)
	sp.Use(h.AuthMiddleware)

	sp.HandleFunc(pat.Get("/"), h.ArticleHomeList)
	sp.HandleFunc(pat.Get("/view"), h.ViewAtTpl)
	// sp.HandleFunc(pat.Get("/feed"), h.FeedHandler)
	sp.HandleFunc(pat.Get("/robots.txt"), h.Robots)

	fs := http.FileServer(http.Dir("static/captcha"))
	sp.Handle(pat.Get("/captcha/*"), http.StripPrefix("/captcha/", fs))

	sp.HandleFunc(pat.Get("/node/:cid"), h.CategoryDetailNew)
	sp.HandleFunc(pat.Get("/member/:uid"), h.UserDetail)
	// sp.HandleFunc(pat.Get("/tag/:tag"), h.TagDetail)
	sp.HandleFunc(pat.Get("/search"), h.SearchDetail)

	sp.HandleFunc(pat.Get("/logout"), h.UserLogout)
	// sp.HandleFunc(pat.Get("/notification"), h.UserNotification)

	sp.HandleFunc(pat.Get("/topic/:aid"), h.ArticleDetail)
	// sp.HandleFunc(pat.Post("/topic/:aid"), h.ArticleDetailPost)

	sp.HandleFunc(pat.Get("/setting"), h.UserSetting)
	sp.HandleFunc(pat.Post("/setting"), h.UserSettingPost)

	// sp.HandleFunc(pat.Get("/newpost/:cid"), h.ArticleAdd)
	// sp.HandleFunc(pat.Post("/newpost/:cid"), h.ArticleAddPost)

	sp.HandleFunc(pat.Get("/login"), h.UserLogin)
	sp.HandleFunc(pat.Post("/login"), h.UserLoginPost)
	sp.HandleFunc(pat.Get("/register"), h.UserLogin)
	sp.HandleFunc(pat.Post("/register"), h.UserLoginPost)

	// sp.HandleFunc(pat.Get("/qqlogin"), h.QQOauthHandler)
	// sp.HandleFunc(pat.Get("/oauth/qq/callback"), h.QQOauthCallback)
	// sp.HandleFunc(pat.Get("/wblogin"), h.WeiboOauthHandler)
	// sp.HandleFunc(pat.Get("/oauth/wb/callback"), h.WeiboOauthCallback)

	sp.HandleFunc(pat.Post("/content/preview"), h.ContentPreviewPost)
	// sp.HandleFunc(pat.Post("/file/upload"), h.FileUpload)

	return sp
}

func NewFlarumAdminRouter(app *system.Application, sp *goji.Mux) *goji.Mux {
	app.Logger.Notice("Init flarum admin router")

	// https://flarum.yjzq.fun/admin#/basics
	// 管理员页面使用的是不同的路由导向的, 在goji中, 它们都会到/admin这个路径
	sp.HandleFunc(pat.Get(model.FlarumAdminPath), ct.MiddlewareArrayToChains(
		[]ct.HTTPMiddleWareFunc{
			ct.MustAuthMiddleware,
			ct.MustAdminUser,
			ct.IsInAdmin,
		},
		ct.AdminHome,
	))

	extAPISP := goji.SubMux()
	sp.Handle(pat.New(model.FlarumExtensionAPI+"/*"), extAPISP)
	// 修改扩展的配置调用类似如下的api
	// https://flarum.yjzq.fun/api/extensions/flarum-mentions

	// adminSP.HandleFunc(pat.Get("/"), ct.MiddlewareArrayToChains(
	// 	[]ct.HTTPMiddleWareFunc{
	// 		ct.MustAuthMiddleware,
	// 		ct.MustAdminUser,
	// 	},
	// 	ct.AdminHome,
	// ))

	return extAPISP
}

// NewFlarumRouter flarum的router
func NewFlarumRouter(app *system.Application, sp *goji.Mux) *goji.Mux {
	app.Logger.Notice("Init flarum router")
	h := ct.BaseHandler{App: app}

	sp.Use(h.InitMiddlewareContext)
	sp.Use(h.AuthMiddleware)
	sp.Use(ct.RealIPMiddleware)
	sp.Use(ct.AdjustLocaleMiddleware)

	sp.HandleFunc(pat.Get("/"), ct.MiddlewareArrayToChains(
		[]ct.HTTPMiddleWareFunc{
			// ct.TestMiddleware,
			// ct.TestMiddleware2,
		},
		ct.FlarumIndex,
	))

	sp.HandleFunc(pat.Get("/tags"), ct.FlarumIndex)

	// robots.txt
	sp.HandleFunc(pat.Get("/robots.txt"), h.Robots)

	// 用户相关
	sp.HandleFunc(pat.Post("/register"), ct.FlarumUserRegister)
	sp.HandleFunc(pat.Post("/login"), ct.FlarumUserLogin)
	sp.HandleFunc(pat.Get("/logout"), ct.FlarumUserLogout)

	// 语言包支持
	sp.HandleFunc(pat.Get("/locale/:locale/flarum-lang.js"), h.GetLocaleData)
	sp.HandleFunc(pat.Get("/locale/:locale/admin-lang.js"), h.GetLocaleData)

	fs := http.FileServer(http.Dir("static/captcha"))
	sp.Handle(pat.Get("/captcha/*"), http.StripPrefix("/captcha/", fs))

	sp.HandleFunc(pat.Get("/auth/github"), ct.GithubOauthHandler)
	sp.HandleFunc(pat.Get("/auth/github/callback"), ct.GithubOauthCallbackHandler)

	//	discussion
	sp.HandleFunc(pat.Get("/d/:aid"), ct.FlarumArticleDetail)
	sp.HandleFunc(pat.Get("/d/:aid/:sn"), ct.FlarumArticleDetail) // startNumber
	sp.HandleFunc(pat.Post("/d/:aid"), ct.FlarumArticleDetail)

	sp.HandleFunc(pat.Get("/t/:tag"), ct.FlarumIndex)

	// user
	sp.HandleFunc(pat.Get("/u/:username"), ct.FlarumUserPage)
	// 获取用户的设置 GET请求
	sp.HandleFunc(pat.Get("/settings"), ct.MustAuthMiddleware(ct.FlarumUserSettings))

	NewFlarumAPIRouter(app, sp)
	if app.CanServeAdmin() {
		NewFlarumAdminRouter(app, sp)
	}

	return sp
}

func NewFlarumAPIRouter(app *system.Application, sp *goji.Mux) *goji.Mux {
	apiSP := goji.SubMux()
	sp.Handle(pat.New(model.FlarumAPIPath+"/*"), apiSP)
	apiSP.Use(ct.InAPIMiddleware)
	apiSP.HandleFunc(pat.Get("/users/:uid"), ct.MiddlewareArrayToChains(
		[]ct.HTTPMiddleWareFunc{
			ct.MustAuthMiddleware,
		},
		ct.FlarumUser,
	))

	// API handler
	// 获取全部的帖子信息
	apiSP.HandleFunc(pat.Get("/discussions"), ct.FlarumAPIDiscussions)

	// 获取某个帖子的详细信息 GET请求
	apiSP.HandleFunc(pat.Get("/discussions/:aid"), ct.FlarumArticleDetail)

	// 获取帖子的详细信息, POST请求
	// 与上面不同的是, 这里的请求中可能携带有当前登录用户阅读到的位置, 将其进行记录
	apiSP.HandleFunc(pat.Post("/discussions/:aid"), ct.MiddlewareArrayToChains(
		[]ct.HTTPMiddleWareFunc{
			ct.MustAuthMiddleware,
		},
		ct.FlarumArticleDetail,
	))

	// 创建一篇帖子
	apiSP.HandleFunc(pat.Post("/discussions"), ct.MiddlewareArrayToChains(
		[]ct.HTTPMiddleWareFunc{
			ct.MustAuthMiddleware,
			ct.MustCSRFMiddleware,
		},
		ct.FlarumAPICreateDiscussion,
	))

	apiSP.HandleFunc(pat.Post("/posts/:cid"), ct.MiddlewareArrayToChains(
		[]ct.HTTPMiddleWareFunc{
			ct.MustAuthMiddleware,
			ct.MustCSRFMiddleware,
		},
		ct.FlarumCommentsUtils,
	))

	apiSP.HandleFunc(pat.Get("/new_captcha"), ct.NewCaptcha)

	apiSP.HandleFunc(pat.Get("/posts"), ct.FlarumComments)

	apiSP.HandleFunc(pat.Get("/posts/:cid"), ct.FlarumComments)

	// 创建一篇评论
	apiSP.HandleFunc(pat.Post("/posts"), ct.MiddlewareArrayToChains(
		[]ct.HTTPMiddleWareFunc{
			ct.MustAuthMiddleware,
			ct.MustCSRFMiddleware,
		},
		ct.FlarumAPICreatePost,
	))

	apiSP.HandleFunc(pat.Post("/users/:uid"), ct.MiddlewareArrayToChains(
		[]ct.HTTPMiddleWareFunc{
			ct.MustAuthMiddleware,
		},
		ct.FlarumUserUpdate,
	))

	apiSP.HandleFunc(pat.Get("/users"), ct.FlarumConfirmUserAndPost)

	return sp
}
