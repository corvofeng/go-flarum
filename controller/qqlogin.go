package controller

import (
	"encoding/json"
	"goyoubbs/lib/qqOAuth"
	"goyoubbs/model"
	"goyoubbs/util"
	"github.com/ego008/youdb"
	"github.com/rs/xid"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (h *BaseHandler) QQOauthHandler(w http.ResponseWriter, r *http.Request) {
	scf := h.App.Cf.Site
	qq, err := qqOAuth.NewQQOAuth(strconv.Itoa(scf.QQClientID), scf.QQClientSecret, scf.MainDomain+"/oauth/qq/callback")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	// qqOAuth.Logging = true

	now := time.Now().UTC().Unix()
	qqURLState := strconv.FormatInt(now, 10)[6:]

	urlStr, err := qq.GetAuthorizationURL(qqURLState)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	h.SetCookie(w, "QQURLState", qqURLState, 1)
	http.Redirect(w, r, urlStr, http.StatusSeeOther)
}

func (h *BaseHandler) QQOauthCallback(w http.ResponseWriter, r *http.Request) {
	qqURLState := h.GetCookie(r, "QQURLState")
	if len(qqURLState) == 0 {
		w.Write([]byte(`qqURLState cookie missed`))
		return
	}

	scf := h.App.Cf.Site
	qq, err := qqOAuth.NewQQOAuth(strconv.Itoa(scf.QQClientID), scf.QQClientSecret, scf.MainDomain+"/oauth/qq/callback")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	// qqOAuth.Logging = true

	code := r.FormValue("code")
	if code == "" {
		w.Write([]byte("Invalid code"))
		return
	}

	state := r.FormValue("state")
	if state != qqURLState {
		w.Write([]byte("Invalid state"))
		return
	}

	token, err := qq.GetAccessToken(code)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	openid, err := qq.GetOpenID(token.AccessToken)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	timeStamp := uint64(time.Now().UTC().Unix())

	db := h.App.Db
	rs := db.Hget("oauth_qq", []byte(openid.OpenID))
	if rs.State == "ok" {
		// login
		obj := model.QQ{}
		json.Unmarshal(rs.Data[0], &obj)
		uobj, err := model.UserGetByID(db, obj.UID)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		sessionid := xid.New().String()
		uobj.LastLoginTime = timeStamp
		uobj.Session = sessionid
		jb, _ := json.Marshal(uobj)
		db.Hset("user", youdb.I2b(uobj.ID), jb)
		h.SetCookie(w, "SessionID", strconv.FormatUint(uobj.ID, 10)+":"+sessionid, 365)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	profile, err := qq.GetUserInfo(token.AccessToken, openid.OpenID)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	if profile.Ret != 0 {
		w.Write([]byte(profile.Message))
		return
	}

	// register

	siteCf := h.App.Cf.Site
	if siteCf.CloseReg {
		w.Write([]byte(`{"retcode":400,"retmsg":"stop to new register"}`))
		return
	}

	name := util.RemoveCharacter(profile.Nickname)
	name = strings.TrimSpace(strings.Replace(name, " ", "", -1))
	if len(name) == 0 {
		name = "qq"
	}
	nameLow := strings.ToLower(name)
	i := 1
	for {
		if db.Hget("user_name2uid", []byte(nameLow)).State == "ok" {
			i++
			nameLow = name + strconv.Itoa(i)
		} else {
			name = nameLow
			break
		}
	}

	userID, _ := db.HnextSequence("user")
	flag := 5
	if siteCf.RegReview {
		flag = 1
	}
	if userID == 1 {
		flag = 99
	}

	gender := "female"
	if profile.Gender == "男" {
		gender = "male"
	}

	uobj := model.User{
		ID:            userID,
		Name:          name,
		Flag:          flag,
		Gender:        gender,
		RegTime:       timeStamp,
		LastLoginTime: timeStamp,
		Session:       xid.New().String(),
	}

	uidStr := strconv.FormatUint(userID, 10)
	savePath := "static/avatar/" + uidStr + ".jpg"
	err = util.FetchAvatar(profile.Avatar, savePath, r.UserAgent())
	if err != nil {
		err = util.GenerateAvatar(gender, name, 73, 73, savePath)
	}
	if err != nil {
		uobj.Avatar = "0"
	} else {
		uobj.Avatar = uidStr
	}

	jb, _ := json.Marshal(uobj)
	db.Hset("user", youdb.I2b(uobj.ID), jb)
	db.Hset("user_name2uid", []byte(nameLow), youdb.I2b(userID))
	db.Hset("user_flag:"+strconv.Itoa(flag), youdb.I2b(uobj.ID), []byte(""))

	obj := model.QQ{
		UID:    userID,
		Name:   name,
		Openid: openid.OpenID,
	}
	jb, _ = json.Marshal(obj)
	db.Hset("oauth_qq", []byte(openid.OpenID), jb)

	h.SetCookie(w, "SessionID", strconv.FormatUint(uobj.ID, 10)+":"+uobj.Session, 365)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
