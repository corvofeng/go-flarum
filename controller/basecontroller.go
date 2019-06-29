package controller

import (
	"errors"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"../model"
	"../system"
)

var mobileRegexp = regexp.MustCompile(`Mobile|iP(hone|od|ad)|Android|BlackBerry|IEMobile|Kindle|NetFront|Silk-Accelerated|(hpw|web)OS|Fennec|Minimo|Opera M(obi|ini)|Blazer|Dolfin|Dolphin|Skyfire|Zune`)

type (
	BaseHandler struct {
		App *system.Application
	}

	PageData struct {
		SiteCf        *system.SiteConf
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
		NewestNodes   []model.CategoryMini
	}
	normalRsp struct {
		Retcode int    `json:"retcode"`
		Retmsg  string `json:"retmsg"`
	}
)

func (h *BaseHandler) Render(w http.ResponseWriter, tpl string, data interface{}, tplPath ...string) error {
	if len(tplPath) == 0 {
		return errors.New("File path can not be empty ")
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Server", "GoYouBBS")

	tplDir := h.App.Cf.Main.ViewDir + "/" + tpl + "/"
	tmpl := template.New("youbbs")
	for _, tpath := range tplPath {
		tmpl = template.Must(tmpl.ParseFiles(tplDir + tpath))
	}
	err := tmpl.Execute(w, data)

	return err
}

// CurrentUser 当前用户
// 原有的策略是保存用户到文件中, 现在经过重新改写, 将从数据库中获取用户,
// 但session的使用与原来一致, 仍然从文件中加载, 为了减轻数据库的负担.
func (h *BaseHandler) CurrentUser(w http.ResponseWriter, r *http.Request) (model.User, error) {
	var user model.User
	var uid uint64
	var err error

	db := h.App.Db
	sqlDB := h.App.MySQLdb

	ssValue := h.GetCookie(r, "SessionID")
	if len(ssValue) == 0 {
		return user, errors.New("SessionID cookie not found ")
	}
	z := strings.Split(ssValue, ":")
	rawUID := z[0]
	sessionID := z[1]

	if len(rawUID) > 0 {
		uid, err = strconv.ParseUint(rawUID, 10, 64)
		if err != nil {
			return user, nil
		}
	}
	// 首先通过数据库获取当前用户
	user, err = model.SQLUserGetByID(sqlDB, uid)

	if err != nil {
		return user, err
	}

	// 但是session仍然使用原文件, 是通过用户名来获取session值
	uobj, err := model.UserGetByName(db, user.Name)
	if sessionID == uobj.Session {
		h.SetCookie(w, "SessionID", ssValue, 365)
		return user, nil
	}

	return user, errors.New("user not found")
}

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

func (h *BaseHandler) CurrentTpl(r *http.Request) string {
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
