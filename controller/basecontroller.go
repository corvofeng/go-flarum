package controller

import (
	"encoding/json"
	"errors"
	"goyoubbs/util"
	"html/template"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"goyoubbs/model"
	"goyoubbs/system"
)

var mobileRegexp = regexp.MustCompile(`Mobile|iP(hone|od|ad)|Android|BlackBerry|IEMobile|Kindle|NetFront|Silk-Accelerated|(hpw|web)OS|Fennec|Minimo|Opera M(obi|ini)|Blazer|Dolfin|Dolphin|Skyfire|Zune`)

type (
	// BaseHandler 基础handler
	BaseHandler struct {
		App   *system.Application
		InAPI bool
	}

	// PageData 每个页面中的基础信息
	PageData struct {
		SiteCf        *model.SiteConf
		Title         string
		Keywords      string
		Description   string
		IsMobile      bool
		CurrentUser   model.User
		PageName      string // index/post_add/post_detail/...
		ShowPostTopAd bool
		ShowPostBotAd bool
		ShowSideAd    bool
		HotNodes      []model.CategoryMini
		NewestNodes   []model.Category
		SiteInfo      model.SiteInfo
		PrimaryColor  string
	}
	response struct {
		Retcode int    `json:"retcode"`
		Retmsg  string `json:"retmsg"`
	}
	normalRsp = response // .. deprecated: 2020-05-29 Please don't use it
)

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

// Jsonify 序列化结构体并进行返回
func (h *BaseHandler) Jsonify(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return json.NewEncoder(w).Encode(data)
}

// CurrentUser 当前用户
// 原有的策略是保存用户到文件中, 现在经过重新改写, 将从数据库中获取用户,
func (h *BaseHandler) CurrentUser(w http.ResponseWriter, r *http.Request) (model.User, error) {
	var (
		user model.User
		uid  uint64
		err  error
	)
	sqlDB := h.App.MySQLdb

	ssValue := h.GetCookie(r, "SessionID")
	if len(ssValue) == 0 {
		return user, errors.New("SessionID cookie not found ")
	}
	z := strings.Split(ssValue, ":")
	rawUID := z[0]

	if len(rawUID) > 0 {
		uid, err = strconv.ParseUint(rawUID, 10, 64)
		if err != nil {
			return user, nil
		}
	}
	// TODO: 直接通过数据库获取当前用户, 性能瓶颈了再说
	user, err = model.SQLUserGetByID(sqlDB, uid)
	if util.CheckError(err, "获取用户") {
		return user, err
	}

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
