package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"zoe/model"
	"zoe/model/flarum"

	"goji.io/pat"
)

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

	_, err = aobj.CreateFlarumTopic(gormDB, tagsArray)
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
