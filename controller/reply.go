package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"goyoubbs/model"
	"goyoubbs/model/flarum"
	"goyoubbs/util"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/op/go-logging"
)

// ContentPreviewPost 预览主题以及评论
func (h *BaseHandler) ContentPreviewPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	token := h.GetCookie(r, "token")
	if len(token) == 0 {
		w.Write([]byte(`{"retcode":400,"retmsg":"token cookie missed"}`))
		return
	}

	currentUser, _ := h.CurrentUser(w, r)
	if !currentUser.CanCreateTopic() || !currentUser.CanReply() {
		w.Write([]byte(`{"retcode":403,"retmsg":"forbidden"}`))
		return
	}

	type recForm struct {
		Act     string `json:"act"`
		Link    string `json:"link"`
		Content string `json:"content"`
	}

	type response struct {
		normalRsp
		Content string        `json:"content"`
		Html    template.HTML `json:"html"`
	}

	decoder := json.NewDecoder(r.Body)
	var rec recForm
	err := decoder.Decode(&rec)
	if err != nil {
		w.Write([]byte(`{"retcode":400,"retmsg":"json Decode err:` + err.Error() + `"}`))
		return
	}
	defer r.Body.Close()

	// db := h.App.Db
	rsp := response{}

	if rec.Act == "preview" && len(rec.Content) > 0 {
		rsp.Retcode = 200
		rsp.Html = template.HTML(util.ContentFmt(rec.Content))
	}
	json.NewEncoder(w).Encode(rsp)

}

type replyFilter struct {
	FT    filterType
	AID   uint64 // 一个帖子的评论
	CID   uint64 // 单个评论的信息
	UID   uint64 // 某个用户创建的评论
	Page  uint64
	Limit uint64
}

// 获取评论的信息
// eArticle: 获取一条帖子下方的评论信息
// eUserPost: 获取用户的最新评论
// ePost: 获取一条评论信息
func createFlarumReplyAPIDoc(
	logger *logging.Logger, sqlDB *sql.DB, redisDB *redis.Client,
	appConf model.AppConf,
	siteInfo model.SiteInfo,
	currentUser *model.User,
	inAPI bool,
	rf replyFilter,
	tz int,
) (flarum.CoreData, error) {
	var err error
	coreData := flarum.NewCoreData()
	apiDoc := &coreData.APIDocument

	// 当前全部的评论资源: 数据库中得到
	var comments []model.CommentListItem
	// 当前全部的评论资源: API返回
	var flarumPosts []flarum.Resource

	// 所有分类的信息, 用于整个站点的信息
	var flarumTags []flarum.Resource

	// 当前的话题信息
	var curDisscussion *flarum.Resource

	if rf.FT == eArticle { // 获取一个帖子的所有评论
		pageInfo := model.SQLCommentListByPage(sqlDB, redisDB, rf.AID, rf.Limit, tz)
		comments = pageInfo.Items
	} else if rf.FT == ePost {
		pageInfo := model.SQLCommentListByPage(sqlDB, redisDB, rf.AID, rf.Limit, tz)
		comments = pageInfo.Items
	} else if rf.FT == eUserPost {
		pageInfo := model.SQLCommentListByUser(sqlDB, redisDB, rf.UID, rf.Limit, tz)
		comments = pageInfo.Items
	} else {
		return coreData, fmt.Errorf("Can't process filter: %s", rf.FT)
	}

	if len(comments) == 0 {
		logger.Errorf("Can't get any comment for %d", rf.AID)
	}

	allUsers := make(map[uint64]bool)       // 用于保存已经添加的用户, 进行去重
	allDiscussions := make(map[uint64]bool) // 用于保存已经添加的帖子, 进行去重

	// 添加当前用户, 以及session信息
	if currentUser != nil {
		user := model.FlarumCreateCurrentUser(*currentUser)
		allUsers[user.GetID()] = true
		coreData.AddCurrentUser(user)
		if !inAPI { // 做API请求时, 不更新csrf信息
			coreData.AddSessionData(user, currentUser.RefreshCSRF(redisDB))
		}
	}

	// 针对某个话题时, 这里直接进行添加
	if rf.FT == eArticle || rf.FT == ePost {
		article, err := model.SQLArticleGetByID(sqlDB, redisDB, rf.AID)
		if err != nil {
			logger.Warning("Can't get article: ", rf.AID, err)
		} else {
			diss := model.FlarumCreateDiscussion(article.ToArticleListItem(sqlDB, redisDB, tz))
			curDisscussion = &diss
			apiDoc.AppendResourcs(*curDisscussion)
		}
		allDiscussions[rf.AID] = true
	}

	for _, comment := range comments {
		post := model.FlarumCreatePost(comment)
		apiDoc.AppendResourcs(post)
		flarumPosts = append(flarumPosts, post)

		if _, ok := allUsers[comment.UID]; !ok {
			u, err := model.SQLUserGetByID(sqlDB, comment.UID)
			if err != nil {
				logger.Warningf("Get user %d error: %s", comment.UID, err)
			} else {
				user := model.FlarumCreateUser(u)
				allUsers[comment.UID] = true
				coreData.AppendResourcs(user)
			}
		}

		if _, ok := allDiscussions[comment.AID]; !ok {
			article, err := model.SQLArticleGetByID(sqlDB, redisDB, comment.AID)
			if err != nil {
				logger.Warning("Can't get article: ", comment.AID, err)
			} else {
				apiDoc.AppendResourcs(model.FlarumCreateDiscussion(article.ToArticleListItem(sqlDB, redisDB, tz)))
			}
			allDiscussions[comment.AID] = true
		}
	}
	// 针对当前的话题, 补全其关系信息
	if curDisscussion != nil {
		postRelation := model.FlarumCreatePostRelations(flarumPosts)
		curDisscussion.BindRelations("Posts", postRelation)
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

	if rf.FT == eArticle {
		apiDoc.SetData(*curDisscussion) // 主要信息为当前帖子
	} else if rf.FT == ePost {
		comment, err := model.SQLGetCommentByID(sqlDB, redisDB, rf.CID, tz)
		if err != nil {
			logger.Error("Get comment error:", err)
		}
		commentListItem := model.CommentListItem{Comment: comment}
		post := model.FlarumCreatePost(commentListItem)
		apiDoc.SetData(post) // 主要信息为这条评论
	} else if rf.FT == eUserPost {
		apiDoc.SetData(flarumPosts) // 主要信息为全部评论
	}
	// apiDoc.Links["first"] = "https://flarum.yjzq.fun/api/v1/flarum/discussions?sort=&page%5Blimit%5D=20"
	// apiDoc.Links["next"] = "https://flarum.yjzq.fun/api/v1/flarum/discussions?sort=&page%5Blimit%5D=20"

	return coreData, nil
}

// FlarumAPICreatePost flarum进行评论的接口
func FlarumAPICreatePost(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h
	sqlDB := h.App.MySQLdb
	redisDB := h.App.RedisDB
	scf := h.App.Cf.Site
	si := model.GetSiteInfo(redisDB)
	logger := ctx.GetLogger()

	type PostedReply struct {
		Data struct {
			Type       string `json:"type"`
			Attributes struct {
				Content string `json:"content"`
			} `json:"attributes"`
			Relationships struct {
				Discussion struct {
					Data struct {
						Type string `json:"type"`
						ID   string `json:"id"`
					} `json:"data"`
				} `json:"discussion"`
			} `json:"relationships"`
		} `json:"data"`
	}

	reply := PostedReply{}
	err := json.NewDecoder(r.Body).Decode(&reply)
	if err != nil {
		h.flarumErrorMsg(w, "解析json错误:"+err.Error())
		return
	}
	aid, err := strconv.ParseUint(reply.Data.Relationships.Discussion.Data.ID, 10, 64)
	if err != nil {
		h.flarumErrorMsg(w, "无法获取正确的帖子信息:"+err.Error())
		return
	}

	now := uint64(time.Now().UTC().Unix())
	comment := model.Comment{
		CommentBase: model.CommentBase{
			AID:      aid,
			UID:      ctx.currentUser.ID,
			Content:  reply.Data.Attributes.Content,
			Number:   1,
			ClientIP: ctx.realIP,
			AddTime:  now,
		},
	}

	if ok, err := comment.CreateFlarumComment(sqlDB); !ok {
		h.flarumErrorMsg(w, "创建评论出现错误:"+err.Error())
		return
	}

	rf := replyFilter{
		FT:    ePost,
		AID:   comment.AID,
		CID:   comment.ID,
		Limit: comment.Number,
	}

	coreData, err := createFlarumReplyAPIDoc(logger, sqlDB, redisDB, *h.App.Cf, si, ctx.currentUser, ctx.inAPI, rf, scf.TimeZone)
	if err != nil {
		h.flarumErrorMsg(w, "查询评论出现错误:"+err.Error())
		return
	}

	h.jsonify(w, coreData.APIDocument)
}

// FlarumConfirmUserAndPost 确认当前的用户的评论信息
// FIXME: 这个函数我只知道是在评论时, @其他用户时会调用这个接口, 但是接口具体的行为还不太了解
func FlarumConfirmUserAndPost(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h
	scf := h.App.Cf.Site
	// sqlDB := h.App.MySQLdb
	// redisDB := h.App.RedisDB
	// logger := ctx.GetLogger()

	_filter := strings.TrimSpace(r.FormValue("filter[q]"))
	_pageLimit := r.FormValue("page[limit]")

	// filterData := strings.Split(_filter, "#")
	// if len(filterData) != 2 {
	// 	h.flarumErrorJsonify(w, createSimpleFlarumError("给定的回复信息有误"))
	// 	return
	// }

	// pageLimit, err := strconv.ParseUint(_pageLimit, 10, 64)
	// if err != nil {
	// 	logger.Error(err)
	// 	h.flarumErrorJsonify(w, createSimpleFlarumError("页面限制信息给定错误"))
	// 	return
	// }

	// postID, err := strconv.ParseUint(filterData[1], 10, 64)
	// if err != nil {
	// 	logger.Error(err)
	// 	h.flarumErrorJsonify(w, createSimpleFlarumError(fmt.Sprintf("无法解析评论信息: %s", filterData)))
	// 	return
	// }
	// comment := model.SQLGetCommentByID(sqlDB, redisDB, postID, scf.TimeZone)
	// if comment.UserName != filterData[0] {
	// 	logger.Warningf("用户与评论信息不符合: %s", filterData)
	// }
	coreData := flarum.NewCoreData()
	apiDoc := &coreData.APIDocument // 注意, 获取到的是指针

	apiDoc.Links["first"] = scf.MainDomain + model.FlarumAPIPath + "/users?" +
		fmt.Sprintf("filter%%5Bq%%5D=%s&page%%5Blimit%%5D=%s", url.QueryEscape(_filter), _pageLimit)
		// fmt.Sprintf("filter%%5Bq%%5D=%s%%23%d+&page%%5Blimit%%5D=%d", comment.UserName, comment.ID, pageLimit)

	h.jsonify(w, apiDoc)
	return
}

// FlarumComments 获取评论
func FlarumComments(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	logger := ctx.GetLogger()
	h := ctx.h

	parm := r.URL.Query()
	_userID := parm.Get("filter[user]")
	_disscussionID := parm.Get("filter[discussion]")
	// _type := parm.Get("filter[type]")
	_limit := parm.Get("page[limit]")
	// _sort := parm.Get("sort")
	sqlDB := h.App.MySQLdb
	redisDB := h.App.RedisDB
	inAPI := ctx.inAPI

	var limit uint64
	var err error
	var user model.User
	var coreData flarum.CoreData

	if len(_limit) > 0 {
		limit, err = strconv.ParseUint(_limit, 10, 64)
		if err != nil {
			return
		}
	}
	limit = 20
	var rf replyFilter

	if _userID != "" {
		user, err = model.SQLUserGet(sqlDB, _userID)
		if user.ID == 0 || err != nil {
			h.flarumErrorJsonify(w, createSimpleFlarumError("Can't get the user for: "+_userID+err.Error()))
			return
		}

		rf = replyFilter{
			FT:    eUserPost,
			UID:   user.ID,
			Limit: limit,
		}
	} else if _disscussionID != "" {
		aid, err := strconv.ParseUint(_disscussionID, 10, 64)
		if err != nil {
			logger.Error("Can't get discussion id for ", _disscussionID)
			h.flarumErrorJsonify(w, createSimpleFlarumError("Can't get the article for: "+_disscussionID+err.Error()))
			return
		}
		article, err := model.SQLArticleGetByID(sqlDB, redisDB, aid)
		if err != nil {
			logger.Error("Can't get discussion id for ", aid)
			h.flarumErrorJsonify(w, createSimpleFlarumError("Can't get discussion for: "+_disscussionID+err.Error()))
			return
		}

		rf = replyFilter{
			FT:    eArticle,
			AID:   aid,
			Limit: article.Comments,
		}
	}
	coreData, err = createFlarumReplyAPIDoc(logger, sqlDB, redisDB, *h.App.Cf, model.GetSiteInfo(redisDB), ctx.currentUser, ctx.inAPI, rf, ctx.h.App.Cf.Site.TimeZone)
	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("Get api doc error"+err.Error()))
		return
	}

	// fmt.Println(userID, _type, _limit, _sort, limit, user, comments)
	// 如果是API直接进行返回
	if inAPI {
		h.jsonify(w, coreData.APIDocument)
		return
	}

	return
}
