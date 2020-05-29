package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"goyoubbs/model"
	"goyoubbs/util"

	"github.com/dchest/captcha"
	"github.com/rs/xid"
	"goji.io/pat"
)

// UserLogin 用户登录与注册页面
func (h *BaseHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	type pageData struct {
		PageData
		Act       string
		Token     string
		CaptchaID string
	}
	act := strings.TrimLeft(r.RequestURI, "/")
	title := "登录"
	if act == "register" {
		title = "注册"
	}

	tpl := h.CurrentTpl(r)
	evn := &pageData{}
	evn.SiteCf = h.App.Cf.Site
	evn.Title = title
	evn.Keywords = ""
	evn.Description = ""
	evn.IsMobile = tpl == "mobile"

	evn.ShowSideAd = true
	evn.PageName = "user_login_register"

	evn.Act = act
	evn.CaptchaID = model.NewCaptcha()

	token := h.GetCookie(r, "token")
	if len(token) == 0 {
		token := xid.New().String()
		h.SetCookie(w, "token", token, 1)
	}

	h.Render(w, tpl, evn, "layout.html", "userlogin.html")
}

// UserLoginPost 用于用户登录及注册接口
// 保存密码时, 用户前端传来的密码为md5值, 因此我们也不需要保存明文密码, 也就不需要token了
func (h *BaseHandler) UserLoginPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8") // .. deprecated: 2020-05-29

	rsp := response{}
	token := h.GetCookie(r, "token")
	if len(token) == 0 {
		rsp = response{400, "token cookie missed"}
		h.Jsonify(w, rsp)
		return
	}

	act := strings.TrimLeft(r.RequestURI, "/")

	type recForm struct {
		Name            string `json:"name"`
		Password        string `json:"password"`
		CaptchaID       string `json:"captchaID"`
		CaptchaSolution string `json:"captchaSolution"`
	}

	decoder := json.NewDecoder(r.Body)
	var rec recForm
	err := decoder.Decode(&rec)
	if err != nil {
		rsp = normalRsp{
			400,
			"表单解析错误:" + err.Error(),
		}
		h.Jsonify(w, rsp)
		return
	}
	defer r.Body.Close()

	if len(rec.Name) == 0 || len(rec.Password) == 0 {
		rsp = normalRsp{400, "name or pw is empty"}
		h.Jsonify(w, rsp)
		return
	}
	nameLow := strings.ToLower(rec.Name)
	if !util.IsUserName(nameLow) {
		rsp = response{400, "name fmt err"}
		h.Jsonify(w, rsp)
		return
	}
	// 返回并且携带新的验证码
	type captchaData struct {
		response
		NewCaptchaID string `json:"newCaptchaID"`
	}

	var respCaptcha captchaData
	if !captcha.VerifyString(rec.CaptchaID, rec.CaptchaSolution) {
		respCaptcha = captchaData{
			response{405, "验证码错误"},
			model.NewCaptcha(),
		}
		h.Jsonify(w, respCaptcha)
		return
	}

	sqlDB := h.App.MySQLdb
	redisDB := h.App.RedisDB
	timeStamp := uint64(time.Now().UTC().Unix())

	if act == "login" {
		uobj, err := model.SQLUserGetByName(sqlDB, nameLow)

		if err != nil {
			respCaptcha = captchaData{
				response{405, "登录失败, 请检查用户名与密码"},
				model.NewCaptcha(),
			}
			h.Jsonify(w, respCaptcha)
			return
		}
		if uobj.Password != rec.Password {
			respCaptcha = captchaData{
				response{405, "登录失败, 请检查用户名与密码"},
				model.NewCaptcha(),
			}
			h.Jsonify(w, respCaptcha)
			return
		}
		sessionid := xid.New().String()
		uobj.LastLoginTime = timeStamp
		uobj.Session = sessionid
		jb, _ := json.Marshal(uobj)
		redisDB.HSet("user", fmt.Sprintf("%d", uobj.ID), jb)
		h.SetCookie(w, "SessionID", strconv.FormatUint(uobj.ID, 10)+":"+sessionid, 365)
	} else {
		// register
		siteCf := h.App.Cf.Site
		if siteCf.QQClientID > 0 || siteCf.WeiboClientID > 0 {
			rsp = response{400, "请用QQ 或 微博一键登录"}
			h.Jsonify(w, rsp)
			return
		}
		if siteCf.CloseReg {
			rsp = response{400, "已经停用用户注册"}
			h.Jsonify(w, rsp)
			return
		}
		if _, err := model.SQLUserGetByName(sqlDB, nameLow); err == nil {
			respCaptcha = captchaData{
				response{405, "用户名已经存在"},
				model.NewCaptcha(),
			}
			h.Jsonify(w, respCaptcha)
			return
		}

		uobj := model.User{
			Name:          rec.Name,
			Password:      rec.Password,
			RegTime:       timeStamp,
			LastLoginTime: timeStamp,
			Session:       xid.New().String(),
		}
		uobj.SQLRegister(sqlDB)

		// uidStr := strconv.FormatUint(uobj.ID, 10)
		// err = util.GenerateAvatar("male", rec.Name, 73, 73, "static/avatar/"+uidStr+".jpg")
		// if err != nil {
		// 	uobj.Avatar = "0"
		// } else {
		// 	uobj.Avatar = uidStr
		// }

		// jb, _ := json.Marshal(uobj)
		// db.Hset("user", youdb.I2b(uobj.ID), jb)
		// db.Hset("user_name2uid", []byte(nameLow), youdb.I2b(uobj.ID))
		// db.Hset("user_flag:"+strconv.Itoa(flag), youdb.I2b(uobj.ID), []byte(""))

		h.SetCookie(w, "SessionID", strconv.FormatUint(uobj.ID, 10)+":"+uobj.Session, 365)
	}

	h.DelCookie(w, "token")

	rsp.Retcode = 200
	h.Jsonify(w, rsp)
}

// UserNotification 用户消息
func (h *BaseHandler) UserNotification(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := h.CurrentUser(w, r)
	if currentUser.ID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	type pageData struct {
		PageData
		PageInfo model.ArticlePageInfo
	}

	db := h.App.Db
	scf := h.App.Cf.Site

	tpl := h.CurrentTpl(r)

	evn := &pageData{}
	evn.SiteCf = scf
	evn.Title = "站内提醒 - " + scf.Name
	evn.IsMobile = tpl == "mobile"

	evn.CurrentUser = currentUser
	evn.ShowSideAd = true
	evn.PageName = "user_notification"
	// evn.HotNodes = model.CategoryHot(db, scf.CategoryShowNum)
	// evn.NewestNodes = model.CategoryNewest(db, scf.CategoryShowNum)
	evn.PageInfo = model.ArticleNotificationList(db, currentUser.Notice, scf.TimeZone)

	h.Render(w, tpl, evn, "layout.html", "notification.html")
}

// UserLogout 用户退出登录
func (h *BaseHandler) UserLogout(w http.ResponseWriter, r *http.Request) {
	cks := []string{"SessionID", "QQURLState", "WeiboURLState", "token"}
	for _, k := range cks {
		h.DelCookie(w, k)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// UserDetail 用户详情页
func (h *BaseHandler) UserDetail(w http.ResponseWriter, r *http.Request) {
	act, btn, key, score := r.FormValue("act"), r.FormValue("btn"), r.FormValue("key"), r.FormValue("score")
	if len(key) > 0 {
		_, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			w.Write([]byte(`{"retcode":400,"retmsg":"key type err"}`))
			return
		}
	}
	if len(score) > 0 {
		_, err := strconv.ParseUint(score, 10, 64)
		if err != nil {
			w.Write([]byte(`{"retcode":400,"retmsg":"score type err"}`))
			return
		}
	}

	db := h.App.Db
	redisDB := h.App.RedisDB
	sqlDB := h.App.MySQLdb
	scf := h.App.Cf.Site

	uid := pat.Param(r, "uid")
	uidi, err := strconv.ParseUint(uid, 10, 64)
	if err != nil {
		uid = model.UserGetIDByName(db, strings.ToLower(uid))
		if uid == "" {
			w.Write([]byte(`{"retcode":400,"retmsg":"uid type err"}`))
			return
		}
		http.Redirect(w, r, "/member/"+uid, 301)
		return
	}

	cmd := "rscan"
	if btn == "prev" {
		cmd = "scan"
	}

	uobj, err := model.SQLUserGetByID(sqlDB, uidi)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	currentUser, _ := h.CurrentUser(w, r)

	if uobj.Hidden && !currentUser.IsAdmin() {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"retcode":404,"retmsg":"not found"}`))
		return
	}

	var pageInfo model.ArticlePageInfo

	if act == "reply" {
		tb := "user_article_reply:" + uid
		// pageInfo = model.UserArticleList(db, cmd, tb, key, h.App.Cf.Site.PageShowNum)
		pageInfo = model.ArticleList(db, "z"+cmd, tb, key, score, scf.PageShowNum, scf.TimeZone)
	} else {
		act = "post"
		tb := "user_article_timeline:" + uid
		pageInfo = model.UserArticleList(db, "h"+cmd, tb, key, scf.PageShowNum, scf.TimeZone)
	}

	type userDetail struct {
		model.User
		RegTimeFmt string
	}
	type pageData struct {
		PageData
		Act      string
		Uobj     userDetail
		PageInfo model.ArticlePageInfo
	}

	tpl := h.CurrentTpl(r)

	evn := &pageData{}
	evn.SiteCf = scf
	evn.Title = uobj.Name + " - " + scf.Name
	evn.Keywords = uobj.Name
	evn.Description = uobj.About
	evn.IsMobile = tpl == "mobile"

	evn.CurrentUser = currentUser
	evn.ShowSideAd = true
	evn.PageName = "category_detail"
	// evn.HotNodes = model.CategoryHot(db, scf.CategoryShowNum)
	// evn.NewestNodes = model.CategoryNewest(db, scf.CategoryShowNum)

	evn.Act = act
	evn.Uobj = userDetail{
		User:       uobj,
		RegTimeFmt: util.TimeFmt(uobj.RegTime, "2006-01-02 15:04", scf.TimeZone),
	}
	evn.PageInfo = pageInfo
	evn.SiteInfo = model.GetSiteInfo(redisDB)

	h.Render(w, tpl, evn, "layout.html", "user.html")
}
