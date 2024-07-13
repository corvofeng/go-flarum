package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/corvofeng/go-flarum/model"

	"goji.io/pat"
)

func FlarumDiscussionEdit(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h
	// gormDB := h.App.GormDB
	// inAPI := ctx.inAPI
	// scf := h.App.Cf.Site
	// redisDB := h.App.RedisDB
	logger := ctx.GetLogger()

	type PostedDiscussion struct {
		Data struct {
			Type       string `json:"type"`
			ID         string `json:"id"`
			Attributes struct {
				IsSticky           bool   `json:"isSticky"`
				LastReadPostNumber uint64 `json:"lastReadPostNumber"`
			} `json:"attributes"`
		} `json:"data"`
	}
	qf := PostedDiscussion{}
	bytedata, err := io.ReadAll(r.Body)
	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("Read body error:"+err.Error()))
		return
	}
	err = json.Unmarshal(bytedata, &qf)
	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("json Decode err:"+err.Error()))
		return
	}

	ctx.actionRecords = string(bytedata)
	logger.Debugf("Update %s,%s with: %s", qf.Data.Type, qf.Data.ID, string(bytedata))
}

// FlarumDiscussionDetail 获取flarum中的某篇帖子
// TODO: #12
func FlarumDiscussionDetail(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h
	inAPI := ctx.inAPI
	scf := h.App.Cf.Site
	redisDB := h.App.RedisDB
	logger := ctx.GetLogger()

	type PostedDiscussion struct {
		Data struct {
			Type       string `json:"type"`
			ID         string `json:"id"`
			Attributes struct {
				IsSticky           bool   `json:"isSticky"`
				LastReadPostNumber uint64 `json:"lastReadPostNumber"`
			} `json:"attributes"`
		} `json:"data"`
	}
	getLastReadPostNumber := false

	qf := PostedDiscussion{}
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

	article, err := model.SQLArticleGetByID(h.App.GormDB, redisDB, aid)
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
	coreData, err := createFlarumPostAPIDoc(
		ctx, h.App.GormDB, redisDB, *h.App.Cf, rf, scf.TimeZone)

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

	gormDB := h.App.GormDB
	redisDB := h.App.RedisDB
	logger := ctx.GetLogger()
	scf := h.App.Cf.Site

	// 用户创建的话题
	type PostedDiscussion struct {
		Data struct {
			Type       string `json:"type"`
			Attributes struct {
				Title        string `json:"title"`
				Content      string `json:"content"`
				Subscription string `json:"subscription"`
				IsSticky     bool   `json:"isSticky"`
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
	bytedata, err := io.ReadAll(r.Body)
	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("Read body error:"+err.Error()))
		return
	}
	logger.Debugf("Upate discussion with: %s", string(bytedata))

	err = json.Unmarshal(bytedata, &diss)
	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("json Decode err:"+err.Error()))
		return
	}
	err = model.CreateActionRecord(gormDB, ctx.currentUser.ID, string(bytedata))
	if err != nil {
		logger.Error("Can't create action record", err)
		h.flarumErrorJsonify(w, createSimpleFlarumError("Can't create action record"+err.Error()))
		return
	}

	aobj := model.Topic{
		UserID:       ctx.currentUser.ID,
		Title:        diss.Data.Attributes.Title,
		Content:      diss.Data.Attributes.Content,
		CommentCount: 1,
		ClientIP:     ctx.realIP,
		IsSticky:     diss.Data.Attributes.IsSticky,
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

	_, err = aobj.CreateFlarumTopic(gormDB)
	if err != nil {
		logger.Error("Can't create topic", err)
		h.flarumErrorJsonify(w, createSimpleFlarumError("Can't create topic"+err.Error()))
		return
	}

	rf := replyFilter{
		FT:  eArticle,
		AID: aobj.ID,
	}
	coreData, err := createFlarumPostAPIDoc(ctx, h.App.GormDB, redisDB, *h.App.Cf, rf, scf.TimeZone)

	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("Get api doc error"+err.Error()))
		return
	}

	// 刷新当前的页面展示
	// TODO: 优化逻辑, 不进行全局处理
	go model.TimelyResort()

	h.jsonify(w, coreData.APIDocument)
}
