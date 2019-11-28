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
	evn.CaptchaID = captcha.New()
	model.SaveImage(evn.CaptchaID)
	// captcha.WriteImage()

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
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	token := h.GetCookie(r, "token")
	if len(token) == 0 {
		w.Write([]byte(`{"retcode":400,"retmsg":"token cookie missed"}`))
		return
	}

	act := strings.TrimLeft(r.RequestURI, "/")

	type recForm struct {
		Name            string `json:"name"`
		Password        string `json:"password"`
		CaptchaID       string `json:"captchaID"`
		CaptchaSolution string `json:"captchaSolution"`
	}

	type response struct {
		normalRsp
	}

	decoder := json.NewDecoder(r.Body)
	var rec recForm
	err := decoder.Decode(&rec)
	if err != nil {
		w.Write([]byte(`{"retcode":400,"retmsg":"json Decode err:` + err.Error() + `"}`))
		return
	}
	defer r.Body.Close()

	if len(rec.Name) == 0 || len(rec.Password) == 0 {
		w.Write([]byte(`{"retcode":400,"retmsg":"name or pw is empty"}`))
		return
	}
	nameLow := strings.ToLower(rec.Name)
	if !util.IsUserName(nameLow) {
		w.Write([]byte(`{"retcode":400,"retmsg":"name fmt err"}`))
		return
	}

	if !captcha.VerifyString(rec.CaptchaID, rec.CaptchaSolution) {
		w.Write([]byte(`{"retcode":405,"retmsg":"验证码错误","newCaptchaID":"` + captcha.New() + `"}`))
		return
	}

	db := h.App.Db
	sqlDB := h.App.MySQLdb
	redisDB := h.App.RedisDB
	timeStamp := uint64(time.Now().UTC().Unix())

	if act == "login" {
		bn := "user_login_token"
		key := []byte(token + ":loginerr")
		if db.Zget(bn, key).State == "ok" {
			// 	// todo
			// 	//w.Write([]byte(`{"retcode":400,"retmsg":"name and pw not match"}`))
			// 	//return
		}
		uobj, err := model.SQLUserGetByName(sqlDB, nameLow)
		fmt.Println(uobj)

		if err != nil {
			w.Write([]byte(`{"retcode":405,"retmsg":"json Decode err:` + err.Error() + `","newCaptchaID":"` + captcha.New() + `"}`))
			return
		}
		if uobj.Password != rec.Password {
			w.Write([]byte(`{"retcode":405,"retmsg":"name and pw not match","newCaptchaID":"` + captcha.New() + `"}`))
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
			w.Write([]byte(`{"retcode":400,"retmsg":"请用QQ 或 微博一键登录"}`))
			return
		}
		if siteCf.CloseReg {
			w.Write([]byte(`{"retcode":400,"retmsg":"stop to new register"}`))
			return
		}
		if _, err := model.SQLUserGetByName(sqlDB, nameLow); err == nil {
			w.Write([]byte(`{"retcode":405,"retmsg":"name is exist","newCaptchaID":"` + captcha.New() + `"}`))
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

	rsp := response{}
	rsp.Retcode = 200
	json.NewEncoder(w).Encode(rsp)
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
	evn.HotNodes = model.CategoryHot(db, scf.CategoryShowNum)
	evn.NewestNodes = model.CategoryNewest(db, scf.CategoryShowNum)
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

	if uobj.Hidden && currentUser.Flag < 99 {
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
	evn.HotNodes = model.CategoryHot(db, scf.CategoryShowNum)
	evn.NewestNodes = model.CategoryNewest(db, scf.CategoryShowNum)

	evn.Act = act
	evn.Uobj = userDetail{
		User:       uobj,
		RegTimeFmt: util.TimeFmt(uobj.RegTime, "2006-01-02 15:04", scf.TimeZone),
	}
	evn.PageInfo = pageInfo
	evn.SiteInfo = model.GetSiteInfo(redisDB, db)

	h.Render(w, tpl, evn, "layout.html", "user.html")
}
