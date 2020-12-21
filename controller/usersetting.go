package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"zoe/model"
	"zoe/util"

	"github.com/rs/xid"
)

// UserSetting 用户配置修改
func (h *BaseHandler) UserSetting(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := h.CurrentUser(w, r)
	if currentUser.ID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	type pageData struct {
		BasePageData
		Uobj model.User
		Now  int64
	}
	redisDB := h.App.RedisDB

	tpl := h.CurrentTpl(r)
	evn := &pageData{}
	evn.SiteCf = h.App.Cf.Site
	evn.Title = "设置"
	evn.Keywords = ""
	evn.Description = ""
	evn.IsMobile = tpl == "mobile"
	evn.CurrentUser = currentUser

	evn.ShowSideAd = true
	evn.PageName = "user_setting"
	evn.SiteInfo = model.GetSiteInfo(redisDB)

	evn.Uobj = currentUser
	evn.Now = time.Now().UTC().Unix()

	h.SetCookie(w, "token", xid.New().String(), 1)
	h.Render(w, tpl, evn, "layout.html", "usersetting.html")
}

// UserSettingPost 用户修改资料
func (h *BaseHandler) UserSettingPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	token := h.GetCookie(r, "token")
	if len(token) == 0 {
		w.Write([]byte(`{"retcode":400,"retmsg":"token cookie missed"}`))
		return
	}

	currentUser, _ := h.CurrentUser(w, r)
	if currentUser.ID == 0 {
		w.Write([]byte(`{"retcode":401,"retmsg":"authored require"}`))
		return
	}

	logger := h.App.Logger
	sqlDB := h.App.MySQLdb

	// r.ParseForm() // don't use ParseForm !important
	act := r.FormValue("act")
	if act == "avatar" {

		r.ParseMultipartForm(32 << 20)

		file, _, err := r.FormFile("avatar")
		defer file.Close()

		buff := make([]byte, 512)
		file.Read(buff)
		if len(util.CheckImageType(buff)) == 0 {
			w.Write([]byte(`{"retcode":400,"retmsg":"unknown image format"}`))
			return
		}

		var imgData bytes.Buffer
		file.Seek(0, 0)
		if fileSize, err := io.Copy(&imgData, file); err != nil {
			w.Write([]byte(`{"retcode":400,"retmsg":"read image data err ` + err.Error() + `"}`))
			return
		} else {
			if fileSize > 5360690 {
				w.Write([]byte(`{"retcode":400,"retmsg":"image size too much"}`))
				return
			}
		}

		img, err := util.GetImageObj(&imgData)
		if err != nil {
			w.Write([]byte(`{"retcode":400,"retmsg":"fail to get image obj ` + err.Error() + `"}`))
			return
		}

		uid := strconv.FormatUint(currentUser.ID, 10)
		avatarPath := fmt.Sprintf("static/avatar/%s.jpg", uid)
		logger.Debug("Save avatar", avatarPath)

		err = util.AvatarResize(img, 73, 73, avatarPath)
		if err != nil {
			w.Write([]byte(`{"retcode":400,"retmsg":"fail to resize avatar ` + err.Error() + `"}`))
			return
		}

		// 存储到数据库时, 需要带上前面的'/', 保证为绝对路径
		currentUser.SaveAvatar(sqlDB, h.App.RedisDB, "/"+avatarPath)

		http.Redirect(w, r, "/setting#2", http.StatusSeeOther)
		return
	}

	type recForm struct {
		Act       string `json:"act"`
		Email     string `json:"email"`
		URL       string `json:"url"`
		About     string `json:"about"`
		Password0 string `json:"password0"`
		Password  string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	var rec recForm
	err := decoder.Decode(&rec)
	if err != nil {
		w.Write([]byte(`{"retcode":400,"retmsg":"json Decode err:` + err.Error() + `"}`))
		return
	}
	defer r.Body.Close()

	recAct := rec.Act
	if len(recAct) == 0 {
		w.Write([]byte(`{"retcode":400,"retmsg":"missed act "}`))
		return
	}

	isChanged := false
	if recAct == "info" {
		currentUser.Email = rec.Email
		currentUser.URL = rec.URL
		currentUser.About = rec.About
		isChanged = true
	} else if recAct == "change_pw" {
		if len(rec.Password0) == 0 || len(rec.Password) == 0 {
			w.Write([]byte(`{"retcode":400,"retmsg":"missed args"}`))
			return
		}
		if currentUser.Password != rec.Password0 {
			w.Write([]byte(`{"retcode":400,"retmsg":"当前密码不正确"}`))
			return
		}
		currentUser.Password = rec.Password
		isChanged = true
	} else if recAct == "set_pw" {
		if len(rec.Password) == 0 {
			w.Write([]byte(`{"retcode":400,"retmsg":"missed args"}`))
			return
		}
		currentUser.Password = rec.Password
		isChanged = true
	}

	rlt := true
	if isChanged {
		rlt = currentUser.SQLUserUpdate(sqlDB)
	}

	type response struct {
		normalRsp
	}

	rsp := response{}
	if rlt {
		rsp.Retcode = 200
		rsp.Retmsg = "修改成功"
	} else {
		rsp.Retcode = 400
		rsp.Retmsg = "修改失败"
	}
	json.NewEncoder(w).Encode(rsp)
}
