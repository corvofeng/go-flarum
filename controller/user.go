package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"goyoubbs/model"
	"goyoubbs/model/flarum"
	"goyoubbs/util"

	"github.com/dchest/captcha"
	"github.com/go-redis/redis/v7"
	"github.com/op/go-logging"
	"github.com/rs/xid"
	"goji.io/pat"
)

// UserLogin 用户登录与注册页面
func (h *BaseHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	type pageData struct {
		BasePageData
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
		h.jsonify(w, rsp)
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
		h.jsonify(w, rsp)
		return
	}
	defer r.Body.Close()

	if len(rec.Name) == 0 || len(rec.Password) == 0 {
		rsp = normalRsp{400, "name or pw is empty"}
		h.jsonify(w, rsp)
		return
	}
	nameLow := strings.ToLower(rec.Name)
	if !util.IsUserName(nameLow) {
		rsp = response{400, "name fmt err"}
		h.jsonify(w, rsp)
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
		h.jsonify(w, respCaptcha)
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
			h.jsonify(w, respCaptcha)
			return
		}
		if uobj.Password != rec.Password {
			respCaptcha = captchaData{
				response{405, "登录失败, 请检查用户名与密码"},
				model.NewCaptcha(),
			}
			h.jsonify(w, respCaptcha)
			return
		}
		sessionid := xid.New().String()
		// uobj.LastLoginTime = timeStamp
		uobj.Session = sessionid
		uobj.CachedToRedis(redisDB)
		h.SetCookie(w, "SessionID", strconv.FormatUint(uobj.ID, 10)+":"+sessionid, 365)

	} else {
		// register
		siteCf := h.App.Cf.Site
		if siteCf.QQClientID > 0 || siteCf.WeiboClientID > 0 {
			rsp = response{400, "请用QQ 或 微博一键登录"}
			h.jsonify(w, rsp)
			return
		}
		if siteCf.CloseReg {
			rsp = response{400, "已经停用用户注册"}
			h.jsonify(w, rsp)
			return
		}
		if _, err := model.SQLUserGetByName(sqlDB, nameLow); err == nil {
			respCaptcha = captchaData{
				response{405, "用户名已经存在"},
				model.NewCaptcha(),
			}
			h.jsonify(w, respCaptcha)
			return
		}

		uobj := model.User{
			Name:     rec.Name,
			Password: rec.Password,
			RegTime:  timeStamp,
			// LastLoginTime: timeStamp,
			Session: xid.New().String(),
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
	h.jsonify(w, rsp)
}

// UserNotification 用户消息
func (h *BaseHandler) UserNotification(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := h.CurrentUser(w, r)
	if currentUser.ID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	type pageData struct {
		BasePageData
		PageInfo model.ArticlePageInfo
	}

	// db := h.App.Db
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
	// evn.PageInfo = model.ArticleNotificationList(db, currentUser.Notice, scf.TimeZone)

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
	fmt.Println(btn)

	// db := h.App.Db
	redisDB := h.App.RedisDB
	sqlDB := h.App.MySQLdb
	scf := h.App.Cf.Site

	uid := pat.Param(r, "uid")
	uidi, err := strconv.ParseUint(uid, 10, 64)
	if err != nil {
		// u, _ := model.SQLUserGetByName(sqlDB, strings.ToLower(uid))
		// uid := u.ID
		// if uid == "" {
		// 	w.Write([]byte(`{"retcode":400,"retmsg":"uid type err"}`))
		// 	return
		// }
		// http.Redirect(w, r, "/member/"+uid, 301)
		// return
	}

	// cmd := "rscan"
	// if btn == "prev" {
	// 	cmd = "scan"
	// }

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
		// tb := "user_article_reply:" + uid
		// pageInfo = model.UserArticleList(db, cmd, tb, key, h.App.Cf.Site.PageShowNum)
		// pageInfo = model.ArticleList(db, "z"+cmd, tb, key, score, scf.PageShowNum, scf.TimeZone)
	} else {
		act = "post"
		// tb := "user_article_timeline:" + uid
		// pageInfo = model.UserArticleList(db, "h"+cmd, tb, key, scf.PageShowNum, scf.TimeZone)
	}

	type userDetail struct {
		model.User
		RegTimeFmt string
	}
	type pageData struct {
		BasePageData
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

// NewCaptcha 获取新的验证码
func NewCaptcha(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h
	// 返回并且携带新的验证码
	type captchaData struct {
		response
		NewCaptchaID string `json:"newCaptchaID"`
	}
	respCaptcha := captchaData{
		response{200, "success"},
		model.NewCaptcha(),
	}
	h.jsonify(w, respCaptcha)
}

// FlarumUserRegister 用户注册
func FlarumUserRegister(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h
	rsp := response{}

	type recForm struct {
		Name            string `json:"username"`
		Password        string `json:"password"`
		CaptchaID       string `json:"captcha-id"`
		CaptchaSolution string `json:"captcha-solution"`
		Email           string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	var rec recForm
	err := decoder.Decode(&rec)
	if err != nil {
		rsp = normalRsp{
			400,
			"表单解析错误:" + err.Error(),
		}
		h.jsonify(w, rsp)
		return
	}
	defer r.Body.Close()

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
		h.jsonify(w, respCaptcha)
		return
	}
	rsp.Retcode = 200
	rsp.Retmsg = "注册成功"

	h.jsonify(w, rsp)
}

// FlarumUserLogin flarum用户登录
func FlarumUserLogin(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h
	rsp := response{}
	type recForm struct {
		Identification  string `json:"identification"`
		Password        string `json:"password"`
		CaptchaID       string `json:"captcha-id"`
		CaptchaSolution string `json:"captcha-solution"`
	}
	decoder := json.NewDecoder(r.Body)
	var rec recForm
	err := decoder.Decode(&rec)
	if err != nil || rec.Identification == "" || rec.Password == "" {
		rsp = normalRsp{400, "数据填写错误:" + err.Error()}
		h.jsonify(w, rsp)
		return
	}
	defer r.Body.Close()

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
		h.jsonify(w, respCaptcha)
		return
	}

	sqlDB := h.App.MySQLdb
	redisDB := h.App.RedisDB
	// timeStamp := uint64(time.Now().UTC().Unix())

	uobj, err := model.SQLUserGetByName(sqlDB, rec.Identification)
	if err != nil {
		respCaptcha = captchaData{
			response{405, "登录失败, 请检查用户名与密码"},
			model.NewCaptcha(),
		}
		h.jsonify(w, respCaptcha)
		return
	}
	if uobj.Password != rec.Password {
		respCaptcha = captchaData{
			response{405, "登录失败, 请检查用户名与密码"},
			model.NewCaptcha(),
		}
		h.jsonify(w, respCaptcha)
		return
	}
	sessionid := xid.New().String()
	// uobj.LastLoginTime = timeStamp
	uobj.Session = sessionid

	uobj.CachedToRedis(redisDB)
	h.SetCookie(w, "SessionID", uobj.StrID()+":"+sessionid, 365)

	rsp.Retcode = 200
	rsp.Retmsg = "登录成功"
	h.jsonify(w, rsp)
}

func userLogout(user model.User, h *BaseHandler, w http.ResponseWriter, r *http.Request) {
	redisDB := h.App.RedisDB
	cks := []string{"SessionID", "QQURLState", "WeiboURLState", "token"}
	for _, k := range cks {
		h.DelCookie(w, k)
	}
	user.CleareRedisCache(redisDB)
}

func createFlarumUserAPIDoc(
	logger *logging.Logger, sqlDB *sql.DB, redisDB *redis.Client,
	appConf model.AppConf,
	siteInfo model.SiteInfo,
	currentUser *model.User,
	inAPI bool,
	tz int,
) (flarum.CoreData, error) {
	var err error
	coreData := flarum.NewCoreData()

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
	categories, err := model.SQLGetNotEmptyCategory(sqlDB, redisDB)
	if err != nil {
		logger.Error("Get all categories error", err)
	}
	for _, category := range categories {
		flarumTags = append(flarumTags, model.FlarumCreateTag(category))
	}
	coreData.AppendResourcs(model.FlarumCreateForumInfo(appConf, siteInfo, flarumTags))
	model.FlarumCreateLocale(&coreData, currentUser)

	return coreData, err
}

// FlarumUserLogout flarum用户注销
func FlarumUserLogout(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h
	rsp := response{}
	redisDB := h.App.RedisDB

	token := r.FormValue("token")
	if token == "" {
		rsp = normalRsp{400, "表单参数解析错误"}
		h.jsonify(w, rsp)
		return
	}
	user, err := h.CurrentUser(w, r)
	if err != nil {
		rsp = normalRsp{400, "用户未登录:" + err.Error()}
		h.jsonify(w, rsp)
		return
	}

	if !user.VerifyCSRFToken(redisDB, token) {
		rsp = normalRsp{400, "csrf错误"}
		h.jsonify(w, rsp)
		return
	}

	userLogout(user, h, w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	rsp.Retcode = 200
	rsp.Retmsg = "登出成功"
	h.jsonify(w, rsp)
	return
}

// FlarumUser flarum用户查询
func FlarumUser(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h
	sqlDB := h.App.MySQLdb
	inAPI := ctx.inAPI

	_userID := pat.Param(r, "uid")
	user, err := model.SQLUserGet(sqlDB, _userID)

	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("获取用户信息错误: "+err.Error()))
		return
	}

	coreData := flarum.NewCoreData()
	apiDoc := &coreData.APIDocument
	apiDoc.SetData(model.FlarumCreateCurrentUser(user))

	if inAPI {
		h.jsonify(w, apiDoc)
		return
	}
	tpl := h.CurrentTpl(r)
	evn := InitPageData(r)
	evn.FlarumInfo = coreData

	h.Render(w, tpl, evn, "layout.html", "index.html")
	return
}

// FlarumUserSettings flarum用户查询
func FlarumUserSettings(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h
	sqlDB := h.App.MySQLdb
	redisDB := h.App.RedisDB
	scf := h.App.Cf.Site
	si := model.GetSiteInfo(redisDB)
	// inAPI := ctx.inAPI
	// coreData := flarum.NewCoreData()
	// apiDoc := &coreData.APIDocument
	// apiDoc.SetData(model.FlarumCreateCurrentUser(*ctx.currentUser))
	logger := ctx.GetLogger()

	tpl := h.CurrentTpl(r)

	coreData, err := createFlarumUserAPIDoc(logger, sqlDB, redisDB, *h.App.Cf, si, ctx.currentUser, ctx.inAPI, scf.TimeZone)
	if err != nil {
		h.flarumErrorMsg(w, "查询用户信息错误:"+err.Error())
	}
	evn := InitPageData(r)
	evn.FlarumInfo = coreData

	h.Render(w, tpl, evn, "layout.html", "index.html")
	return
}

// FlarumUserPage flarum用户查询
func FlarumUserPage(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	// currentUser := ctx.currentUser
	h := ctx.h
	sqlDB := h.App.MySQLdb
	inAPI := ctx.inAPI

	username := pat.Param(r, "username")
	user, err := model.SQLUserGetByName(sqlDB, username)

	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("获取用户信息错误"+err.Error()))
		return
	}

	coreData := flarum.NewCoreData()

	// 添加主站点信息
	redisDB := h.App.RedisDB
	si := model.GetSiteInfo(redisDB)

	coreData.AppendResourcs(model.FlarumCreateForumInfo(*h.App.Cf, si, []flarum.Resource{}))

	apiDoc := &coreData.APIDocument

	// if user.ID == currentUser.ID {
	u := model.FlarumCreateCurrentUser(user)
	coreData.AppendResourcs(u)
	apiDoc.SetData(u)
	currentUser := ctx.currentUser
	// coreData.AppendResourcs(model.FlarumCreateCurrentUser(*ctx.currentUser))
	// 添加当前用户的session信息
	if currentUser != nil {
		user := model.FlarumCreateCurrentUser(*currentUser)
		coreData.AddCurrentUser(user)
		if !inAPI { // 做API请求时, 不更新csrf信息
			coreData.AddSessionData(user, currentUser.RefreshCSRF(redisDB))
		}
	}

	// } else {
	// u := model.FlarumCreateCurrentUser(*currentUser)
	// apiDoc.SetData(u)
	// coreData.AppendResourcs(u)
	// }
	apiDoc.Links["first"] = ""
	apiDoc.Links["next"] = ""

	tpl := h.CurrentTpl(r)
	evn := InitPageData(r)
	evn.FlarumInfo = coreData

	h.Render(w, tpl, evn, "layout.html", "index.html")
	return
}

// FlarumUserUpdate flarum用户更新配置信息
func FlarumUserUpdate(w http.ResponseWriter, r *http.Request) {
	_uid := pat.Param(r, "uid")
	ctx := GetRetContext(r)
	h := ctx.h
	sqlDB := h.App.MySQLdb
	redisDB := h.App.RedisDB
	if ctx.currentUser.StrID() != _uid {
		h.flarumErrorMsg(w, "当期仅允许修改自己的配置")
		return
	}

	type UserUpdate struct {
		Data struct {
			Type       string `json:"type"`
			ID         string `json:"id"`
			Attributes struct {
				Preferences flarum.Preferences `json:"preferences"`
			} `json:"attributes"`
		} `json:"data"`
	}
	userUpdateInfo := UserUpdate{}
	err := json.NewDecoder(r.Body).Decode(&userUpdateInfo)
	if err != nil {
		h.flarumErrorMsg(w, "解析json错误:"+err.Error())
		return
	}
	ctx.currentUser.SetPreference(
		sqlDB, redisDB,
		userUpdateInfo.Data.Attributes.Preferences,
	)
	coreData := flarum.NewCoreData()
	apiDoc := &coreData.APIDocument
	apiDoc.SetData(model.FlarumCreateCurrentUser(*ctx.currentUser))

	if ctx.inAPI {
		h.jsonify(w, apiDoc)
		return
	}
}
