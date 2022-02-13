package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"zoe/model"
	"zoe/model/flarum"
	"zoe/util"

	"github.com/rs/xid"
	"goji.io/pat"
)

// ArticleDetail 帖子的详情
func (h *BaseHandler) ArticleDetail(w http.ResponseWriter, r *http.Request) {
	var start uint64
	var err error

	ctx := GetRetContext(r)
	inAPI := ctx.inAPI

	btn, key, score := r.FormValue("btn"), r.FormValue("key"), r.FormValue("score")
	if len(key) > 0 {
		start, err = strconv.ParseUint(key, 10, 64)
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

	_aid := pat.Param(r, "aid")
	aid, err := strconv.ParseUint(_aid, 10, 64)
	if err != nil {
		w.Write([]byte(`{"retcode":400,"retmsg":"aid type err"}`))
		return
	}

	var commentsCnt uint64
	scf := h.App.Cf.Site
	logger := h.App.Logger
	redisDB := h.App.RedisDB
	sqlDB := h.App.MySQLdb

	// 获取帖子详情
	aobj, err := model.SQLArticleGetByID(h.App.GormDB, sqlDB, redisDB, aid)
	if util.CheckError(err, fmt.Sprintf("获取帖子 %d 失败", aid)) {
		w.Write([]byte(err.Error()))
		return
	}
	aobj.IncrArticleCntFromRedisDB(sqlDB, redisDB)

	// 获取帖子评论数目
	err = sqlDB.QueryRow(
		"SELECT COUNT(*) FROM reply where topic_id = ?", aid,
	).Scan(&commentsCnt)
	if util.CheckError(err, "帖子评论数") {
		w.Write([]byte(err.Error()))
		return
	}

	currentUser, _ := h.CurrentUser(w, r)

	// 获取帖子所在的节点
	// cobj, err := model.SQLCategoryGetByID(sqlDB, strconv.FormatUint(aobj.CID, 10))

	err = nil
	// cobj.Hidden = false
	// if err != nil {
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }

	// if cobj.Hidden && !currentUser.IsAdmin() {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	w.Write([]byte(`{"retcode":404,"retmsg":"not found"}`))
	// 	return
	// }

	// Authorized
	if scf.Authorized && currentUser.Flag < 5 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"retcode":401,"retmsg":"Unauthorized"}`))
		return
	}

	pageInfo := model.SQLCommentList(
		h.App.GormDB,
		sqlDB,
		redisDB,
		aobj.ID,
		start,
		btn,
		scf.CommentListNum,
		scf.TimeZone,
	)

	type articleForDetail struct {
		model.Topic
		ContentFmt  template.HTML
		TagStr      template.HTML
		Name        string
		Avatar      string
		Views       uint64
		AddTimeFmt  string
		EditTimeFmt string
		CommentsCnt uint64
	}

	type pageData struct {
		BasePageData
		Aobj       articleForDetail
		Author     model.User
		Cobj       model.Tag
		Relative   model.ArticleRelative
		PageInfo   model.CommentPageInfo
		Views      uint64
		SiteInfo   model.SiteInfo
		FlarumInfo interface{}
	}

	tpl := h.CurrentTpl(r)
	evn := &pageData{}
	evn.SiteCf = scf
	// evn.Title = aobj.Title + " - " + cobj.Name + " - " + scf.Name
	// evn.Keywords = aobj.Tags
	// evn.Description = cobj.Name + " - " + aobj.Title + " - " + aobj.Tags
	evn.IsMobile = tpl == "mobile"

	evn.CurrentUser = currentUser
	evn.ShowSideAd = true
	evn.PageName = "article_detail"
	// evn.HotNodes = model.CategoryHot(db, scf.CategoryShowNum)
	// evn.NewestNodes = model.CategoryNewest(db, scf.CategoryShowNum)

	author, _ := model.SQLUserGetByID(h.App.GormDB, aobj.UserID)

	if author.ID == 2 {
		// 这部分的网页是转载而来的, 所以需要保持原样式, 这里要牺牲XSS的安全性了
		evn.Aobj = articleForDetail{
			Topic:       aobj,
			ContentFmt:  template.HTML(aobj.Content),
			CommentsCnt: commentsCnt,
			Name:        author.Name,
			Avatar:      author.Avatar,
			Views:       aobj.ClickCnt,
			AddTimeFmt:  util.TimeFmt(aobj.AddTime, "2006-01-02 15:04", scf.TimeZone),
			EditTimeFmt: util.TimeFmt(aobj.EditTime, "2006-01-02 15:04", scf.TimeZone),
		}
	} else {
		evn.Aobj = articleForDetail{
			Topic:       aobj,
			ContentFmt:  template.HTML(model.ContentFmt(aobj.Content)),
			CommentsCnt: commentsCnt,
			Name:        author.Name,
			Avatar:      author.Avatar,
			Views:       aobj.ClickCnt,
			AddTimeFmt:  util.TimeFmt(aobj.AddTime, "2006-01-02 15:04", scf.TimeZone),
			EditTimeFmt: util.TimeFmt(aobj.EditTime, "2006-01-02 15:04", scf.TimeZone),
		}
	}

	// if len(aobj.Tags) > 0 {
	// 	var tags []string
	// 	for _, v := range strings.Split(aobj.Tags, ",") {
	// 		tags = append(tags, `<a href="/tag/`+v+`">`+v+`</a>`)
	// 	}
	// 	evn.Aobj.TagStr = template.HTML(strings.Join(tags, ", "))
	// }

	// evn.Cobj = cobj
	evn.Author = author
	// evn.Relative = model.ArticleGetRelative(aobj.ID, aobj.Tags)
	evn.PageInfo = pageInfo
	evn.SiteInfo = model.GetSiteInfo(redisDB)

	token := h.GetCookie(r, "token")
	if len(token) == 0 {
		token := xid.New().String()
		h.SetCookie(w, "token", token, 1)
	}

	if inAPI {
		type TopicData struct {
			model.RestfulAPIBase
			Data model.RestfulTopic `json:"data"`
		}
		topicData := TopicData{
			RestfulAPIBase: model.RestfulAPIBase{
				State: true,
			},
			Data: model.RestfulTopic{
				ID:      evn.Aobj.ID,
				UID:     evn.Author.ID,
				Content: evn.Aobj.ContentFmt,
				Title:   evn.Aobj.Title,
				Author: model.RestfulUser{
					Name:   evn.Author.Name,
					Avatar: evn.Author.Avatar,
				},
				CreateAt:   evn.Aobj.EditTimeFmt,
				VisitCount: evn.Aobj.Views,
				Replies:    []model.RestfulReply{},
			},
		}
		logger.Debug("This is in api version")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(topicData)
	} else {
		h.Render(w, tpl, evn, "layout.html", "article.html")
	}
}

// FlarumArticleDetail 获取flarum中的某篇帖子
// TODO: #12
func FlarumArticleDetail(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h
	inAPI := ctx.inAPI
	scf := h.App.Cf.Site
	sqlDB := h.App.MySQLdb
	redisDB := h.App.RedisDB
	logger := ctx.GetLogger()

	// _filter := r.FormValue("filter")
	// fmt.Println(_filter)
	// _near := r.FormValue("page[near]")
	// fmt.Println(_near)
	type QueryFilter struct {
		Data struct {
			Type       string `json:"type"`
			ID         string `json:"id"`
			Attributes struct {
				LastReadPostNumber uint64 `json:"lastReadPostNumber"`
			} `json:"attributes"`
		} `json:"data"`
	}
	getLastReadPostNumber := false

	qf := QueryFilter{}
	if inAPI {
		if err := json.NewDecoder(r.Body).Decode(&qf); err == nil {
			getLastReadPostNumber = true
		}
	}

	_aid := pat.Param(r, "aid")
	aid, err := strconv.ParseUint(_aid, 10, 64)
	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("aid type error"))
		return
	}

	// if !getLastReadPostNumber {
	// 	if near, err := strconv.ParseUint(_near, 10, 64); err == nil {
	// 		getLastReadPostNumber = true
	// 		qf.Data.Attributes.LastReadPostNumber = near
	// 	}
	// }

	article, err := model.SQLArticleGetByID(h.App.GormDB, sqlDB, redisDB, aid)
	if err != nil {
		logger.Error("Can't get discussion id for ", aid)
		h.flarumErrorJsonify(w, createSimpleFlarumError("Can't get discussion for: "+_aid+err.Error()))
		return
	}

	rf := replyFilter{
		FT:    eArticle,
		AID:   aid,
		Limit: article.CommentCount,

		LastReadPostNumber: 0,
	}
	if getLastReadPostNumber {
		rf.LastReadPostNumber = qf.Data.Attributes.LastReadPostNumber
	}
	_sn, err := h.safeGetParm(r, "sn")
	if err == nil {
		if sn, err := strconv.ParseUint(_sn, 10, 64); err == nil {
			getLastReadPostNumber = true
			rf.StartNumber = sn
		}
	}
	coreData, err := createFlarumReplyAPIDoc(
		ctx, h.App.GormDB, sqlDB, redisDB, *h.App.Cf, rf, scf.TimeZone)

	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("Get api doc error"+err.Error()))
		return
	}

	// 如果是API直接进行返回
	if inAPI {
		h.jsonify(w, coreData.APIDocument)
		return
	}

	tpl := h.CurrentTpl(r)

	evn := InitPageData(r)
	evn.FlarumInfo = coreData
	h.Render(w, tpl, evn, "layout.html", "article.html")
}

// FlarumAPICreateDiscussion 用户创建一条话题
func FlarumAPICreateDiscussion(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h
	sqlDB := h.App.MySQLdb
	gormDB := h.App.GormDB
	redisDB := h.App.RedisDB
	logger := ctx.GetLogger()
	scf := h.App.Cf.Site

	// 用户创建的话题
	type PostedDiscussion struct {
		Data struct {
			Type       string `json:"type"`
			Attributes struct {
				Title   string `json:"title"`
				Content string `json:"content"`
			} `json:"attributes"`
			Relationships struct {
				Tags struct {
					Data []struct {
						Type string `json:"type"`
						ID   string `json:"id"`
					} `json:"data"`
				} `json:"tags"`
			} `json:"relationships"`
		} `json:"data"`
	}

	diss := PostedDiscussion{}
	err := json.NewDecoder(r.Body).Decode(&diss)

	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("json Decode err:"+err.Error()))
		return
	}

	now := uint64(time.Now().UTC().Unix())
	aobj := model.Topic{
		UserID:       ctx.currentUser.ID,
		Title:        diss.Data.Attributes.Title,
		Content:      diss.Data.Attributes.Content,
		CommentCount: 1,
		AddTime:      now,
		EditTime:     now,
		ClientIP:     ctx.realIP,
		// Active:        1, // 帖子为激活状态
		// FatherTopicID: 0, // 没有原始主题
	}
	tagsArray := flarum.RelationArray{}
	for _, rela := range diss.Data.Relationships.Tags.Data {
		tagID, err := strconv.Atoi(rela.ID)
		if err != nil {
			logger.Warning("Get wrong tag id", rela.ID)
			continue
		}
		tagsArray.Data = append(tagsArray.Data, flarum.InitBaseResources(uint64(tagID), rela.Type))
	}

	for _, rela := range diss.Data.Relationships.Tags.Data {
		var tag model.Tag
		result := gormDB.First(&tag, rela.ID)
		if result.Error != nil {
			logger.Warning("Get wrong tag id", rela.ID)
			continue
		}
		aobj.Tags = append(aobj.Tags, tag)
	}

	_, err = aobj.CreateFlarumDiscussion(gormDB, tagsArray)
	if err != nil {
		logger.Error("Can't create topic", err)
		h.flarumErrorJsonify(w, createSimpleFlarumError("Can't create topic"+err.Error()))
		return
	}

	rf := replyFilter{
		FT:  eArticle,
		AID: aobj.ID,
	}
	coreData, err := createFlarumReplyAPIDoc(ctx, h.App.GormDB, sqlDB, redisDB, *h.App.Cf, rf, scf.TimeZone)

	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("Get api doc error"+err.Error()))
		return
	}

	// 刷新当前的页面展示
	// TODO: 优化逻辑, 不进行全局处理
	go model.TimelyResort()

	h.jsonify(w, coreData.APIDocument)
}
