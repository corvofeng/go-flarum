package router

import (
	"net/http"
	"os"

	ct "goyoubbs/controller"
	"goyoubbs/model"
	"goyoubbs/system"

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
		NewGoYouBBSRouter(app, sp)
	}
	return sp
}

// NewGoYouBBSRouter goyoubbs的router
func NewGoYouBBSRouter(app *system.Application, sp *goji.Mux) *goji.Mux {
	h := ct.BaseHandler{App: app}
	sp.Use(h.InitMiddlewareContext)
	sp.Use(h.AuthMiddleware)

	sp.HandleFunc(pat.Get("/"), h.ArticleHomeList)
	sp.HandleFunc(pat.Get("/luck"), h.IFeelLucky)
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
	sp.HandleFunc(pat.Post("/topic/:aid"), h.ArticleDetailPost)

	sp.HandleFunc(pat.Get("/setting"), h.UserSetting)
	sp.HandleFunc(pat.Post("/setting"), h.UserSettingPost)

	sp.HandleFunc(pat.Get("/newpost/:cid"), h.ArticleAdd)
	sp.HandleFunc(pat.Post("/newpost/:cid"), h.ArticleAddPost)

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

	sp.HandleFunc(pat.Get("/admin/post/edit/:aid"), h.ArticleEdit)
	sp.HandleFunc(pat.Post("/admin/post/edit/:aid"), h.ArticleEditPost)

	if os.Getenv("type") == "admin" {
		sp.HandleFunc(pat.Get("/admin/comment/edit/:aid/:cid"), h.CommentEdit)
		sp.HandleFunc(pat.Post("/admin/comment/edit/:aid/:cid"), h.CommentEditPost)
		sp.HandleFunc(pat.Get("/admin/user/edit/:uid"), h.UserEdit)
		sp.HandleFunc(pat.Post("/admin/user/edit/:uid"), h.UserEditPost)
		sp.HandleFunc(pat.Get("/admin/user/list"), h.AdminUserList)
		sp.HandleFunc(pat.Post("/admin/user/list"), h.AdminUserListPost)
		sp.HandleFunc(pat.Get("/admin/category/list"), h.AdminCategoryList)
		sp.HandleFunc(pat.Post("/admin/category/list"), h.AdminCategoryListPost)
		sp.HandleFunc(pat.Get("/admin/link/list"), h.AdminLinkList)
		sp.HandleFunc(pat.Post("/admin/link/list"), h.AdminLinkListPost)
	}

	return sp
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

	// 用户相关
	sp.HandleFunc(pat.Post("/register"), ct.FlarumUserRegister)
	sp.HandleFunc(pat.Post("/login"), ct.FlarumUserLogin)
	sp.HandleFunc(pat.Get("/logout"), ct.FlarumUserLogout)

	// 语言包支持
	sp.HandleFunc(pat.Get("/locale/:locale/flarum-lang.js"), h.GetLocaleData)

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
