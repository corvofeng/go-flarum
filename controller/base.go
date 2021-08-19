package controller

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"

	"zoe/model"
	"zoe/system"
	"zoe/util"

	"github.com/op/go-logging"
	"goji.io/pat"
)

var mobileRegexp = regexp.MustCompile(`Mobile|iP(hone|od|ad)|Android|BlackBerry|IEMobile|Kindle|NetFront|Silk-Accelerated|(hpw|web)OS|Fennec|Minimo|Opera M(obi|ini)|Blazer|Dolfin|Dolphin|Skyfire|Zune`)

type (
	// BaseHandler 基础handler
	BaseHandler struct { // .. deprecated: 2020-06-11 Please don't use it
		App *system.Application
	}

	// BasePageData 每个页面中的基础信息
	BasePageData struct {
		SiteCf        *model.SiteConf
		Title         string
		Keywords      string
		Description   string
		IsMobile      bool
		IsInAdmin     bool       // 是否在admin页面中
		CurrentUser   model.User // 历史原因: 这里切换成指针有太多的报错, 暂时不处理
		PageName      string     // index/post_add/post_detail/...
		ShowPostTopAd bool
		ShowPostBotAd bool
		ShowSideAd    bool
		HotNodes      []model.CategoryMini
		NewestNodes   []model.Category
		SiteInfo      model.SiteInfo
		PrimaryColor  string
	}
	// response  返回信息
	response struct {
		Retcode int `json:"retcode"`

		Retmsg string `json:"retmsg"`
	}

	normalRsp = response // .. deprecated: 2020-05-29 Please don't use it

	flarumError struct {
		Detail string `json:"detail"`
	}
	// FlarumErrorResponse  flarum API调用时出现的错误
	FlarumErrorResponse struct {
		Errors []flarumError `json:"errors"`
	}

	// ContextKey 记录context的value
	ContextKey int64

	// ReqContext 请求时将会携带的contex信息
	ReqContext struct {
		currentUser *model.User
		inAPI       bool
		inAdmin     bool
		h           *BaseHandler
		realIP      string
		locale      string
		err         error
	}

	// PageData 每个页面中的全部信息
	PageData struct {
		BasePageData

		// 此处的变量为兼容旧的youbbs变量
		PageInfo   model.ArticlePageInfo
		Cobj       model.Category
		Aobj       model.Article
		MainNodes  []model.CategoryMini
		FlarumInfo interface{}

		FlarumJSPrefix string

		PluginHTML map[string]string
	}
)

const (
	ckRequest ContextKey = iota
)

// InitPageData 初始化返回页面
func InitPageData(r *http.Request) *PageData {
	ctx := GetRetContext(r)
	h := ctx.h
	redisDB := h.App.RedisDB

	pd := PageData{
		BasePageData: BasePageData{
			SiteCf:      h.App.Cf.Site,
			Title:       h.App.Cf.Site.Name,
			Description: h.App.Cf.Site.Desc,
			IsInAdmin:   ctx.inAdmin,
			SiteInfo:    model.GetSiteInfo(redisDB),
		},
		PluginHTML:     make(map[string]string),
		FlarumJSPrefix: "forum",
		// SiteInfo: model.GetSiteInfo(redisDB),
	}
	if pd.IsInAdmin {
		pd.FlarumJSPrefix = "admin"
	}

	if ctx.currentUser != nil {
		pd.CurrentUser = *ctx.currentUser
	}

	return &pd
}

// GetRetContext 获取当前上线信息中的自有的context
func GetRetContext(r *http.Request) *ReqContext {
	return r.Context().Value(ckRequest).(*ReqContext)
}

// createSimpleFlarumError 初始化一个最简单的错误值
func createSimpleFlarumError(errMsg string) FlarumErrorResponse {
	return FlarumErrorResponse{[]flarumError{initFlarumError(errMsg)}}
}

// initFlarumError 初始化一个错误值
func initFlarumError(err string) flarumError {
	return flarumError{Detail: err}
}

// Render 渲染html
/**
 * .. version_changed: 2020-05-28 增加了对flaru主题的支持, 将会渲染不同的模板
 */
func (h *BaseHandler) Render(w http.ResponseWriter, tpl string, data interface{}, tplPath ...string) error {
	if len(tplPath) == 0 {
		return errors.New("File path can not be empty ")
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Server", h.App.Cf.Main.ServerName)

	tplDir := path.Join(h.App.Cf.Main.ViewDir, tpl)
	tmpl := template.New(h.App.Cf.Main.ServerStyle).Funcs(template.FuncMap{
		"marshal": func(v interface{}) template.JS {
			a, _ := json.Marshal(v)
			return template.JS(a)
		},
		// 将html中的内容直接进行渲染
		// https://stackoverflow.com/a/44222211
		// https://stackoverflow.com/a/42055211
		"safeHTML": func(str string) template.HTML {
			return template.HTML(str)
		},
	})
	for _, tpath := range tplPath {
		tmpl = template.Must(tmpl.ParseFiles(
			path.Join(tplDir, tpath),
		))
	}
	err := tmpl.Execute(w, data)
	if err != nil {
		h.App.Logger.Error("Can't render template with err", err)
	}

	return err
}

// sp.HandleFunc(pat.Get("/d/:aid/:lrn"), ct.FlarumArticleDetail) // lastReadNumber
// goji的Param会在没有变量时抛出异常, 因此进行catch处理
func (h *BaseHandler) safeGetParm(r *http.Request, parm string) (data string, err error) {
	defer func() {
		if r := recover(); r != nil {
			data = ""
			err = errors.New("Can't get parm")
		}
	}()
	data = pat.Param(r, parm)
	return data, nil
}

// jsonify 序列化结构体并进行返回
func (h *BaseHandler) jsonify(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return json.NewEncoder(w).Encode(data)
}

// flarumErrorJsonify flarum错误需要此函数进行返回
// h.flarumErrorJsonify(w, createSimpleFlarumError("这是其中的错误"))
func (h *BaseHandler) flarumErrorJsonify(w http.ResponseWriter, data FlarumErrorResponse) error {
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return json.NewEncoder(w).Encode(data)
}

// flarumErrorJsonify flarum错误需要此函数进行返回
// h.flarumErrorJsonify(w, createSimpleFlarumError("这是其中的错误"))
func (h *BaseHandler) flarumErrorMsg(w http.ResponseWriter, errMsg string) error {
	return h.flarumErrorJsonify(w, createSimpleFlarumError(errMsg))
}

// CurrentUser 当前用户
// 原有的策略是保存用户到文件中, 现在经过重新改写, 将从数据库中获取用户,
func (h *BaseHandler) CurrentUser(w http.ResponseWriter, r *http.Request) (user model.User, err error) {
	sqlDB := h.App.MySQLdb
	redisDB := h.App.RedisDB
	ssValue := h.GetCookie(r, "SessionID")
	if len(ssValue) == 0 {
		return user, errors.New("SessionID cookie not found ")
	}
	z := strings.Split(ssValue, ":")
	rawUID := z[0]

	user, err = model.RedisGetUserByID(redisDB, rawUID)
	// user, err = model.SQLUserGetByID(sqlDB, user.ID)
	if util.CheckError(err, "获取用户") {
		// 程序运行到这里, 表明redis中已经将其删掉, 用户需要重新登录
		h.DelCookie(w, "SessionID")
		return user, err
	}
	user.GetPreference(sqlDB, redisDB)

	return user, nil
}

// SetCookie 浏览器设置cookie
func (h *BaseHandler) SetCookie(w http.ResponseWriter, name, value string, days int) error {
	encoded, err := h.App.Sc.Encode(name, value)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    encoded,
		Path:     "/",
		Secure:   h.App.Cf.Main.CookieSecure,
		HttpOnly: h.App.Cf.Main.CookieHttpOnly,
		Expires:  time.Now().UTC().AddDate(0, 0, days),
	})
	return err
}

// GetCookie 根据name获取当前所存的cookie值
func (h *BaseHandler) GetCookie(r *http.Request, name string) string {
	if cookie, err := r.Cookie(name); err == nil {
		var value string
		if err = h.App.Sc.Decode(name, cookie.Value, &value); err == nil {
			return value
		}
	}
	return ""
}

// DelCookie 删除Cookie, 用户下线
func (h *BaseHandler) DelCookie(w http.ResponseWriter, name string) {
	if len(name) > 0 {
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			Secure:   h.App.Cf.Main.CookieSecure,
			HttpOnly: h.App.Cf.Main.CookieHttpOnly,
			Expires:  time.Unix(0, 0),
		})
	}
}

// CurrentTpl 当前使用的模板类型
func (h *BaseHandler) CurrentTpl(r *http.Request) string {
	// 如果使用其他主题, 那么直接返回该主题
	serverStyle := h.App.Cf.Main.ServerStyle
	if serverStyle != "youbbs" {
		return serverStyle
	}

	tpl := "desktop"
	//tpl := "mobile"

	cookieTpl := h.GetCookie(r, "tpl")
	if len(cookieTpl) > 0 {
		if cookieTpl == "desktop" || cookieTpl == "mobile" {
			return cookieTpl
		}
	}

	ua := r.Header.Get("User-Agent")
	if len(ua) < 6 {
		return tpl
	}
	if mobileRegexp.MatchString(ua) {
		return "mobile"
	}
	return tpl
}

// GetLogger 获取当前的logger
// TODO: 期望未来能按照用户进行日志打印
func (ctx *ReqContext) GetLogger() *logging.Logger {
	return ctx.h.App.Logger
}
