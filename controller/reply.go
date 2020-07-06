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
// eUser: 获取用户的最新评论
// eReply: 获取一条评论信息
func createFlarumReplyAPIDoc(
	logger *logging.Logger, sqlDB *sql.DB, redisDB *redis.Client,
	appConf model.AppConf,
	siteInfo model.SiteInfo,
	currentUser *model.User,
	inAPI bool,
	rf replyFilter,
	tz int,
) (flarum.CoreData, error) {
	// var err error
	coreData := flarum.NewCoreData()
	apiDoc := &coreData.APIDocument
	// var pageInfo model.CommentPageInfo

	if rf.FT == eArticle { // 获取一个帖子的所有评论
		article, err := model.SQLArticleGetByID(sqlDB, redisDB, rf.AID)
		if err != nil {
			logger.Error("Get article error", err)
			return coreData, err
		}

		diss := model.FlarumCreateDiscussionFromArticle(article)
		pageInfo := model.SQLCommentListByPage(sqlDB, redisDB, article.ID, tz)

		// 获取该文章下面所有的评论信息
		postArr := []flarum.Resource{}
		allUsers := make(map[uint64]bool)
		for _, comment := range pageInfo.Items {
			post := model.FlarumCreatePost(comment)
			apiDoc.AppendResourcs(post)
			postArr = append(postArr, post)

			// 当前用户会在后面统一添加
			if currentUser != nil && currentUser.ID == comment.UID {
				continue
			}

			// 用户不存在则添加, 已经存在的用户不会考虑
			if _, ok := allUsers[comment.UID]; !ok {
				user := model.FlarumCreateUserFromComments(comment)
				apiDoc.AppendResourcs(user)
				allUsers[comment.UID] = true
			}
		}

		// 获取评论的作者
		if len(pageInfo.Items) == 0 {
			logger.Errorf("Can't get any comment for %d", article.ID)
		}

		// 文章当前的分类
		categories, err := model.SQLGetNotEmptyCategory(sqlDB, redisDB)

		if err != nil {
			logger.Error("Get all categories error", err)
			return coreData, err
		}

		// 添加所有分类的信息
		var flarumTags []flarum.Resource
		for _, category := range categories {
			tag := model.FlarumCreateTag(category)
			flarumTags = append(flarumTags, tag)
		}

		postRelation := model.FlarumCreatePostRelations(postArr)
		diss.BindRelations("Posts", postRelation)
		apiDoc.SetData(diss)

		apiDoc.Links["first"] = "https://flarum.yjzq.fun/api/v1/flarum/discussions?sort=&page%5Blimit%5D=20"
		apiDoc.Links["next"] = "https://flarum.yjzq.fun/api/v1/flarum/discussions?sort=&page%5Blimit%5D=20"

		// coreData.APIDocument = apiDoc
		// 添加主站点信息
		coreData.AppendResourcs(model.FlarumCreateForumInfo(appConf, siteInfo, flarumTags))

		// 添加当前用户的session信息
		if currentUser != nil {
			user := model.FlarumCreateCurrentUser(*currentUser)
			coreData.AddCurrentUser(user)
			if !inAPI { // 做API请求时, 不更新csrf信息
				coreData.AddSessionData(user, currentUser.RefreshCSRF(redisDB))
			}
		}

		return coreData, nil

	} else if rf.FT == eReply { // 获取其中一条评论
		comment, err := model.SQLGetCommentByID(sqlDB, redisDB, rf.CID, tz)
		if err != nil {
			logger.Error("Get comment error:", err)
			return coreData, err
		}
		commentListItem := model.CommentListItem{Comment: comment}

		article, err := model.SQLArticleGetByID(sqlDB, redisDB, comment.AID)
		if err != nil {
			logger.Error("Get article error:", err)
			return coreData, err
		}

		diss := model.FlarumCreateDiscussionFromArticle(article)
		post := model.FlarumCreatePost(commentListItem)
		apiDoc.SetData(post)

		if currentUser != nil && comment.UID == currentUser.ID { // 当前用户单独进行添加
			user := model.FlarumCreateCurrentUser(*currentUser)
			apiDoc.AppendResourcs(user)
		} else {
			user := model.FlarumCreateUserFromComments(commentListItem)
			apiDoc.AppendResourcs(user)
		}

		pageInfo := model.SQLCommentListByPage(sqlDB, redisDB, article.ID, tz)

		postArr := []flarum.Resource{}
		for _, comment := range pageInfo.Items {
			post := model.FlarumCreatePost(comment)
			postArr = append(postArr, post)
		}

		postRelation := model.FlarumCreatePostRelations(postArr)
		diss.BindRelations("Posts", postRelation)
		apiDoc.AppendResourcs(diss)

		return coreData, nil
	} else if rf.FT == eUser {
		var err error
		coreData := flarum.NewCoreData()
		apiDoc := &coreData.APIDocument
		postArr := []flarum.Resource{}

		allUsers := make(map[uint64]bool)       // 用于保存已经添加的用户, 进行去重
		allDiscussions := make(map[uint64]bool) // 用于保存已经添加的帖子, 进行去重
		pageInfo := model.SQLCommentListByUser(sqlDB, redisDB, rf.UID, rf.Limit, tz)
		if err != nil {
			logger.Warning("Can't get comments for  user", rf.UID)
		}
		for _, comment := range pageInfo.Items {
			post := model.FlarumCreatePost(comment)
			apiDoc.AppendResourcs(post)
			postArr = append(postArr, post)

			// 当前用户会在后面统一添加
			if currentUser == nil || currentUser.ID != comment.UID {
				if _, ok := allUsers[comment.UID]; !ok {
					user := model.FlarumCreateUserFromComments(comment)
					apiDoc.AppendResourcs(user)
					allUsers[comment.UID] = true
				}
			}

			if _, ok := allDiscussions[comment.AID]; !ok {
				article, err := model.SQLArticleGetByID(sqlDB, redisDB, comment.AID)
				if err != nil {
					logger.Warning("Can't get article: ", comment.AID, err)
				} else {
					apiDoc.AppendResourcs(model.FlarumCreateDiscussionFromArticle(article))
				}
				allDiscussions[comment.AID] = true
			}
		}
		apiDoc.SetData(postArr)

		// 添加当前用户的session信息
		if currentUser != nil {
			user := model.FlarumCreateCurrentUser(*currentUser)
			coreData.AddCurrentUser(user)
			if !inAPI { // 做API请求时, 不更新csrf信息
				coreData.AddSessionData(user, currentUser.RefreshCSRF(redisDB))
			}
		}

		return coreData, err
	}

	return coreData, fmt.Errorf("Can't process filter: %s", rf.FT)
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
		CID: comment.ID,
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

// FlarumComments 获取用户的评论
func FlarumComments(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	logger := ctx.GetLogger()
	h := ctx.h

	parm := r.URL.Query()
	_userID := parm.Get("filter[user]")
	// _type := parm.Get("filter[type]")
	_limit := parm.Get("page[limit]")
	// _sort := parm.Get("sort")
	sqlDB := h.App.MySQLdb
	redisDB := h.App.RedisDB
	inAPI := ctx.inAPI

	var limit uint64
	var userID uint64
	var err error
	var user model.User

	if len(_limit) > 0 {
		limit, err = strconv.ParseUint(_limit, 10, 64)
		if err != nil {
			return
		}
	}
	limit = 20

	// 尝试获取用户
	for true {
		if user, err = model.SQLUserGetByName(sqlDB, _userID); err == nil {
			break
		}
		if userID, err = strconv.ParseUint(_userID, 10, 64); err != nil {
			logger.Error("Can't get user id for ", _userID)
			break
		}
		if user, err = model.SQLUserGetByID(sqlDB, userID); err != nil {
			logger.Error("Can't get user by err: ", err)
			break
		}
		break
	}

	if user.ID == 0 {
		h.flarumErrorJsonify(w, createSimpleFlarumError("Can't get the user for: "+_userID))
		return
	}

	rf := replyFilter{
		FT:    eUser,
		UID:   user.ID,
		Limit: limit,
	}

	coreData, err := createFlarumReplyAPIDoc(logger, sqlDB, redisDB, *h.App.Cf, model.GetSiteInfo(redisDB), ctx.currentUser, ctx.inAPI, rf, ctx.h.App.Cf.Site.TimeZone)
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
