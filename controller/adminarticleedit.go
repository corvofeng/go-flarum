package controller

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"goyoubbs/model"
	"goyoubbs/util"

	"github.com/ego008/youdb"
	"github.com/rs/xid"
	"goji.io/pat"
)

// ArticleEdit 超级管理员可以编辑帖子
func (h *BaseHandler) ArticleEdit(w http.ResponseWriter, r *http.Request) {
	_aid := pat.Param(r, "aid")
	aid, err := strconv.Atoi(_aid)
	if err != nil {
		w.Write([]byte(`{"retcode":400,"retmsg":"cid type err"}`))
		return
	}

	currentUser, _ := h.CurrentUser(w, r)
	if currentUser.ID == 0 {
		w.Write([]byte(`{"retcode":401,"retmsg":"authored err"}`))
		return
	}
	db := h.App.Db
	sqlDB := h.App.MySQLdb
	redisDB := h.App.RedisDB

	aobj, err := model.SQLArticleGetByID(sqlDB, redisDB, uint64(aid))
	if err != nil {
		w.Write([]byte(`{"retcode":403,"retmsg":"aid not found"}`))
		return
	}

	if !currentUser.CanEdit(&aobj.ArticleBase) {
		w.Write([]byte(`{"retcode":403,"retmsg":"flag forbidden}`))
		return
	}

	cobj, err := model.SQLCategoryGetByID(sqlDB, strconv.FormatUint(aobj.CID, 10))
	// cobj, err := model.CategoryGetByID(db, strconv.FormatUint(aobj.CID, 10))
	if err != nil {
		w.Write([]byte(`{"retcode":404,"retmsg":"` + err.Error() + `"}`))
		return
	}

	act := r.FormValue("act")

	if act == "del" {
		aidB := youdb.I2b(aobj.ID)
		// remove
		// 总文章列表
		db.Zdel("article_timeline", aidB)
		// 分类文章列表
		db.Zdel("category_article_timeline:"+strconv.FormatUint(aobj.CID, 10), aidB)
		// 用户文章列表
		db.Hdel("user_article_timeline:"+strconv.FormatUint(aobj.UID, 10), aidB)
		// 分类下文章数
		db.Zincr("category_article_num", youdb.I2b(aobj.CID), -1)
		// 删除标题记录
		hash := md5.Sum([]byte(aobj.Title))
		titleMd5 := hex.EncodeToString(hash[:])
		db.Hdel("title_md5", []byte(titleMd5))

		// set
		db.Hset("article_hidden", aidB, []byte(""))
		aobj.Hidden = true
		jb, _ := json.Marshal(aobj)
		db.Hset("article", aidB, jb)
		uobj, _ := model.UserGetByID(db, aobj.UID)
		if uobj.Articles > 0 {
			uobj.Articles--
		}
		jb, _ = json.Marshal(uobj)
		db.Hset("user", youdb.I2b(uobj.ID), jb)

		// tag send task work，自动处理tag与文章id
		at := model.ArticleTag{
			ID:      aobj.ID,
			OldTags: aobj.Tags,
			NewTags: "",
		}
		jb, _ = json.Marshal(at)
		db.Hset("task_to_set_tag", youdb.I2b(at.ID), jb)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	type pageData struct {
		PageData
		Cobj      model.Category
		MainNodes []model.CategoryMini
		Aobj      model.Article
	}

	tpl := h.CurrentTpl(r)
	evn := &pageData{}
	evn.SiteCf = h.App.Cf.Site
	evn.Title = "修改文章"
	evn.IsMobile = tpl == "mobile"
	evn.CurrentUser = currentUser
	evn.ShowSideAd = true
	evn.PageName = "article_edit"
	evn.SiteInfo = model.GetSiteInfo(redisDB)

	evn.Cobj = cobj
	// evn.MainNodes = model.CategoryGetMain(db, cobj)
	evn.Aobj = aobj

	h.SetCookie(w, "token", xid.New().String(), 1)
	h.Render(w, tpl, evn, "layout.html", "adminarticleedit.html")
}

// ArticleEditPost 超级用户修改帖子
func (h *BaseHandler) ArticleEditPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	_aid := pat.Param(r, "aid")
	aid, err := strconv.Atoi(_aid)
	if err != nil {
		w.Write([]byte(`{"retcode":400,"retmsg":"cid type err"}`))
		return
	}

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

	db := h.App.Db
	sqlDB := h.App.MySQLdb
	redisDB := h.App.RedisDB

	aobj, err := model.SQLArticleGetByID(sqlDB, redisDB, uint64(aid))
	if err != nil {
		w.Write([]byte(`{"retcode":403,"retmsg":"aid not found"}`))
		return
	}

	if !currentUser.CanEdit(&aobj.ArticleBase) {
		w.Write([]byte(`{"retcode":403,"retmsg":"flag forbidden}`))
		return
	}

	// 提交的表单
	type recForm struct {
		Aid          uint64 `json:"aid"`
		Act          string `json:"act"`
		CID          uint64 `json:"cid"`
		Title        string `json:"title"`
		Content      string `json:"content"`
		Tags         string `json:"tags"`
		CloseComment string `json:"closecomment"`
	}

	decoder := json.NewDecoder(r.Body)
	var rec recForm
	err = decoder.Decode(&rec)
	if err != nil {
		w.Write([]byte(`{"retcode":400,"retmsg":"json Decode err:` + err.Error() + `"}`))
		return
	}
	defer r.Body.Close()

	rec.Aid = uint64(aid)

	// aidS := strconv.FormatUint(rec.Aid, 10)
	aidB := youdb.I2b(rec.Aid)

	rec.Title = strings.TrimSpace(rec.Title)
	rec.Content = strings.TrimSpace(rec.Content)
	rec.Tags = util.CheckTags(rec.Tags)

	if rec.Act == "preview" {
		tmp := struct {
			normalRsp
			Html string `json:"html"`
		}{
			normalRsp{200, ""},
			util.ContentFmt(rec.Content),
		}
		json.NewEncoder(w).Encode(tmp)
		return
	}

	// check title
	hash := md5.Sum([]byte(rec.Title))
	titleMd5 := hex.EncodeToString(hash[:])
	rs0 := db.Hget("title_md5", []byte(titleMd5))
	if rs0.State == "ok" && !bytes.Equal(rs0.Data[0], aidB) {
		w.Write([]byte(`{"retcode":403,"retmsg":"title has existed"}`))
		return
	}

	scf := h.App.Cf.Site

	if rec.CID == 0 || len(rec.Title) == 0 {
		w.Write([]byte(`{"retcode":400,"retmsg":"missed args"}`))
		return
	}
	if len(rec.Title) > scf.TitleMaxLen {
		w.Write([]byte(`{"retcode":403,"retmsg":"TitleMaxLen limited"}`))
		return
	}
	if len(rec.Content) > scf.ContentMaxLen {
		w.Write([]byte(`{"retcode":403,"retmsg":"ContentMaxLen limited"}`))
		return
	}

	// 获取当前分类
	_, err = model.SQLCategoryGetByID(sqlDB, strconv.FormatUint(rec.CID, 10))
	if err != nil {
		w.Write([]byte(`{"retcode":404,"retmsg":"` + err.Error() + `"}`))
		return
	}

	var closeComment bool
	if rec.CloseComment == "1" {
		closeComment = true
	}

	if aobj.CID == rec.CID && aobj.Title == rec.Title && aobj.Content == rec.Content && aobj.Tags == rec.Tags && aobj.CloseComment == closeComment {
		w.Write([]byte(`{"retcode":201,"retmsg":"nothing changed"}`))
		return
	}

	oldCID := aobj.CID
	oldTitle := aobj.Title
	oldTags := aobj.Tags

	aobj.CID = rec.CID
	aobj.Title = rec.Title
	aobj.Content = rec.Content
	aobj.Tags = rec.Tags
	aobj.CloseComment = closeComment
	aobj.ClientIP = r.Header.Get("X-REAL-IP")
	aobj.EditTime = uint64(time.Now().UTC().Unix())
	aobj.SQLArticleUpdate(sqlDB, db, redisDB)

	jb, _ := json.Marshal(aobj)
	db.Hset("article", aidB, jb)

	if oldCID != rec.CID {
		db.Zincr("category_article_num", youdb.I2b(rec.CID), 1)
		db.Zincr("category_article_num", youdb.I2b(oldCID), -1)

		db.Zset("category_article_timeline:"+strconv.FormatUint(rec.CID, 10), aidB, aobj.EditTime)
		db.Zdel("category_article_timeline:"+strconv.FormatUint(oldCID, 10), aidB)
	}

	if oldTitle != rec.Title {
		hash0 := md5.Sum([]byte(oldTitle))
		titleMd50 := hex.EncodeToString(hash0[:])
		db.Hdel("title_md5", []byte(titleMd50))
		db.Hset("title_md5", []byte(titleMd5), aidB)
	}

	if oldTags != rec.Tags {
		// tag send task work ，自动处理tag与文章id
		at := model.ArticleTag{
			ID:      aobj.ID,
			OldTags: oldTags,
			NewTags: aobj.Tags,
		}
		jb, _ = json.Marshal(at)
		db.Hset("task_to_set_tag", youdb.I2b(at.ID), jb)
	}

	h.DelCookie(w, "token")

	tmp := struct {
		normalRsp
		Aid uint64 `json:"aid"`
	}{
		normalRsp{200, "ok"},
		aobj.ID,
	}
	json.NewEncoder(w).Encode(tmp)
}
