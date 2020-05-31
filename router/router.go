package router

import (
	"net/http"
	"os"

	// "fmt"
	"goyoubbs/controller"
	"goyoubbs/model"
	"goyoubbs/system"

	// "github.com/dchest/captcha"

	"goji.io"
	"goji.io/pat"
)

// NewRouter create the router
func NewRouter(app *system.Application) *goji.Mux {
	sp := goji.SubMux()
	if app.IsFlarum() {
		NewFlarumAPIRouter(app, sp)
		NewFlarumRouter(app, sp)
	} else {
		NewGoYouBBSRouter(app, sp)
	}
	return sp
}

// NewAPIRouter create api router
func NewAPIRouter(app *system.Application) *goji.Mux {
	sp := goji.SubMux()
	h := controller.BaseHandler{App: app, InAPI: true}

	sp.HandleFunc(pat.Get("/node/:cid"), h.CategoryDetailNew)
	sp.HandleFunc(pat.Get("/topic/:aid"), h.ArticleDetail)
	sp.HandleFunc(pat.Get("/topics"), h.CategoryDetailNew)

	return sp
}

// NewGoYouBBSRouter goyoubbs的router
func NewGoYouBBSRouter(app *system.Application, sp *goji.Mux) *goji.Mux {
	h := controller.BaseHandler{App: app}

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
	h := controller.BaseHandler{App: app}

	sp.HandleFunc(pat.Get("/"), h.ArticleHomeList)

	//	discussion
	// sp.HandleFunc(pat.Get("/"), h.ArticleHomeList)
	sp.HandleFunc(pat.Get("/d/:aid"), h.ArticleDetail)
	sp.HandleFunc(pat.Post("/d/:aid"), h.ArticleDetailPost)

	// user
	sp.HandleFunc(pat.Get("/u/:username"), h.UserDetail)

	return sp
}

// NewFlarumAPIRouter flarum的API
func NewFlarumAPIRouter(app *system.Application, sp *goji.Mux) *goji.Mux {
	app.Logger.Notice("Init flarum api router")
	h := controller.BaseHandler{App: app, InAPI: true}

	sp.HandleFunc(pat.Get(model.FlarumAPIPath+"/discussions"), h.FlarumAPIDiscussions)
	sp.HandleFunc(pat.Get(model.FlarumAPIPath+"/discussions/:aid"), h.FlarumArticleDetail)

	return sp
}
