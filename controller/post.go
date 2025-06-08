package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/corvofeng/go-flarum/model"
	"github.com/corvofeng/go-flarum/model/flarum"
	"github.com/corvofeng/go-flarum/util"

	"github.com/go-redis/redis/v7"
	"goji.io/pat"
	"gorm.io/gorm"
)

type replyFilter struct {
	FT    filterType
	AID   uint64 // 一个帖子的评论
	CID   uint64 // 单个评论的信息
	UID   uint64 // 某个用户创建的评论
	Page  uint64
	Limit uint64
	IDS   []uint64

	RenderLimit uint64 // 当前页面会显示的评论数量, 一般只显示几条

	LastReadPostNumber uint64
	NearNumber         uint64
	StartNumber        uint64
}

// 获取评论的信息
// eArticle: 获取一条帖子下方的评论信息
// eUserPost: 获取用户的最新评论
// ePost: 获取一条评论信息
func createFlarumPostAPIDoc(
	reqctx *ReqContext,
	gormDB *gorm.DB, redisDB *redis.Client,
	appConf model.AppConf,
	rf replyFilter,
	tz int,
) (flarum.CoreData, error) {
	var err error
	coreData := flarum.NewCoreData()
	apiDoc := &coreData.APIDocument
	inAPI := reqctx.inAPI
	currentUser := reqctx.currentUser
	logger := reqctx.GetLogger()
	siteInfo := model.GetSiteInfo(redisDB)

	rf.RenderLimit = 20
	// 当前全部的评论资源: 数据库中得到
	// var comments []model.CommentListItem
	var comments []model.Comment
	// 当前全部的评论资源: API返回
	var flarumPosts []flarum.Resource

	// 所有分类的信息, 用于整个站点的信息
	var flarumTags []flarum.Resource

	// 当前的话题信息
	var curDisscussion *flarum.Resource

	// 使用startNumber时, 多加载一部分数据
	if rf.StartNumber < 10 {
		rf.StartNumber = 1
	} else {
		rf.StartNumber = rf.StartNumber - 10
	}
	logger.Debugf("Get comments with filter: %+v", rf)

	if rf.FT == eArticle { // 获取一个帖子的所有评论
		comments, err = model.SQLCommentListByTopic(gormDB, redisDB, rf.AID, rf.Limit, tz)
	} else if rf.FT == ePost {
		comments, err = model.SQLCommentListByCID(gormDB, redisDB, rf.CID, rf.Limit, tz)
	} else if rf.FT == eUserPost {
		comments, err = model.SQLCommentListByUser(gormDB, redisDB, rf.UID, rf.Limit, tz)
	} else if rf.FT == ePosts { // 根据post列表获取评论
		comments, err = model.SQLCommentListByList(gormDB, redisDB, rf.IDS, tz)
		rf.RenderLimit = uint64(len(rf.IDS))
	} else {
		logger.Warningf("Can't process filter: `%s`", rf.FT)
		return coreData, fmt.Errorf("can't process filter: `%s`", rf.FT)
	}
	if err != nil {
		logger.Error("Get comments error with filter %+v, %s", rf, err)
	}

	commentsLen := uint64(len(comments))
	logger.Debugf("Get %d comments for %d", commentsLen, rf.AID)
	if commentsLen == 0 {
		logger.Errorf("Can't get any comment for %d", rf.AID)
	}

	// 获取恰当的commentsLen值
	if commentsLen < rf.RenderLimit {
		// logger.Warning("Can't get proper comments for", rf.AID)
		rf.RenderLimit = commentsLen
	}

	if rf.AID == 0 && commentsLen != 0 { // 没有AID时, 进行补充
		rf.AID = comments[0].AID
	}

	allUsers := make(map[uint64]bool)       // 用于保存已经添加的用户, 进行去重
	allDiscussions := make(map[uint64]bool) // 用于保存已经添加的帖子, 进行去重

	// 添加当前用户, 以及session信息
	if currentUser != nil {
		user := model.FlarumCreateCurrentUser(*currentUser)
		allUsers[user.GetID()] = true
		coreData.AddCurrentUser(user)
		if !inAPI { // 做API请求时, 不更新csrf信息, 反之则进行更新
			coreData.AddSessionData(user, currentUser.RefreshCSRF(redisDB))
		}
	}

	hasUpdateComments := make(chan bool)

	// 针对某个话题时, 这里直接进行添加
	if rf.FT == eArticle || rf.FT == ePost || rf.FT == ePosts {
		article, err := model.SQLArticleGetByID(gormDB, redisDB, rf.AID)
		if err != nil {
			logger.Warning("Can't get article: ", rf.AID, err)
		} else {
			diss := model.FlarumCreateDiscussion(article)
			curDisscussion = &diss
			apiDoc.AppendResources(*curDisscussion)
		}
		allDiscussions[rf.AID] = true
		if rf.FT == eArticle || rf.FT == ePost { // 查询当前帖子的信息时, 更新redis中的帖子的评论信息, ePost为刚刚添加帖子的操作
			go article.CacheCommentList(redisDB, comments, hasUpdateComments)
		}
	}

	for _, comment := range comments {
		// lastReadPostNumber只用于记录读取到的位置, 不需要返回评论信息
		if rf.LastReadPostNumber != 0 {
			break
		}

		// 使用lastReadPostNumber来标记起始位置
		if rf.StartNumber != 0 && comment.Number < rf.StartNumber {
			continue
		}

		// 使用renderlimit 标记结束位置
		if rf.RenderLimit == 0 {
			break
		}

		if _, ok := allUsers[comment.UID]; !ok {
			u, err := model.SQLUserGetByID(gormDB, comment.UID)
			if err != nil {
				logger.Warningf("Get user %d error: %s", comment.UID, err)
			} else {
				user := model.FlarumCreateUser(u)
				allUsers[comment.UID] = true
				coreData.AppendResources(user)
			}
		}

		if _, ok := allDiscussions[comment.AID]; !ok {
			article, err := model.SQLArticleGetByID(gormDB, redisDB, comment.AID)
			if err != nil {
				logger.Warning("Can't get article: ", comment.AID, err)
			} else {
				apiDoc.AppendResources(model.FlarumCreateDiscussion(article))
			}
			allDiscussions[comment.AID] = true
		}

		// 处理用户的like信息
		for _, userID := range comment.Likes {
			if _, ok := allUsers[userID]; !ok {
				u, err := model.SQLUserGetByID(gormDB, userID)
				if err != nil {
					logger.Warningf("Get user %d error: %s", userID, err)
				} else {
					user := model.FlarumCreateUser(u)
					allUsers[user.GetID()] = true
					coreData.AppendResources(user)
				}
			}
		}

		post := model.FlarumCreatePost(comment, currentUser)
		logger.Debugf("Create comment post for %+v", post)
		apiDoc.AppendResources(post)
		flarumPosts = append(flarumPosts, post)

		rf.RenderLimit--
	}

	// 针对当前的话题, 补全其关系信息
	if curDisscussion != nil {
		if rf.FT == eArticle || rf.FT == ePost { // 如果是查询全部评论, 等待一下
			<-hasUpdateComments
		}
		article, _ := model.SQLArticleGetByID(gormDB, redisDB, rf.AID)
		postRelation := model.FlarumCreatePostRelations([]flarum.Resource{}, article.GetCommentIDList(redisDB))
		curDisscussion.BindRelations("Posts", postRelation)
	}

	// 添加当前站点信息
	tags, err := model.SQLGetTags(gormDB)
	if err != nil {
		logger.Error("Get all categories error", err)
	}
	for _, category := range tags {
		tag := model.FlarumCreateTag(category)
		coreData.AppendResources(tag)
		flarumTags = append(flarumTags, model.FlarumCreateTag(category))
	}
	coreData.AppendResources(model.FlarumCreateForumInfo(
		currentUser,
		appConf, siteInfo,
		flarumTags,
	))

	if rf.FT == eArticle {
		// apiDoc.SetData(flarumPosts) // 主要信息为全部评论
		if rf.NearNumber != 0 {
			apiDoc.SetData(flarumPosts) // 主要信息为全部评论
		} else {
			// if inAPI {
			// 	apiDoc.SetData(flarumPosts) // 主要信息为当前帖子
			// } else {
			// }
			apiDoc.SetData(*curDisscussion) // 主要信息为当前帖子
		}
	} else if rf.FT == ePost {
		// comment, err := model.SQLGetCommentByID(   redisDB, rf.CID, tz)
		// if err != nil {
		// 	logger.Error("Get comment error:", err)
		// }
		// commentListItem := model.CommentListItem{Comment: comment}
		// post := model.FlarumCreatePost(commentListItem, currentUser)
		if len(flarumPosts) >= 0 {
			apiDoc.SetData(flarumPosts[0]) // 主要信息为这条评论
		}
	} else if rf.FT == eUserPost || rf.FT == ePosts {
		apiDoc.SetData(flarumPosts) // 主要信息为全部评论
	}
	logger.Debugf("Update the api doc: %+v", apiDoc)
	// apiDoc.Links["first"] = "https://flarum.yjzq.fun/api/v1/flarum/discussions?sort=&page%5Blimit%5D=20"
	// apiDoc.Links["next"] = "https://flarum.yjzq.fun/api/v1/flarum/discussions?sort=&page%5Blimit%5D=20"
	model.FlarumCreateLocale(&coreData, reqctx.locale)

	return coreData, nil
}

// FlarumAPICreatePost flarum进行评论的接口
func FlarumAPICreatePost(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h

	redisDB := h.App.RedisDB
	scf := h.App.Cf.Site
	// logger := ctx.GetLogger()

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

	now := util.TimeNow()
	comment := model.Comment{
		Reply: model.Reply{
			AID:      aid,
			UID:      ctx.currentUser.ID,
			Content:  reply.Data.Attributes.Content,
			Number:   1,
			ClientIP: ctx.realIP,
			AddTime:  now,
		},
	}
	comment.Content = model.PreProcessUserMention(h.App.GormDB, redisDB, scf.TimeZone, comment.Content)

	if ok, err := comment.CreateFlarumComment(h.App.GormDB); !ok {
		h.flarumErrorMsg(w, "创建评论出现错误:"+err.Error())
		return
	}

	rf := replyFilter{
		FT:    ePost,
		AID:   comment.AID,
		CID:   comment.ID,
		Limit: comment.Number,
	}

	coreData, err := createFlarumPostAPIDoc(ctx, h.App.GormDB, redisDB, *h.App.Cf, rf, scf.TimeZone)
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
	//
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
	// comment := model.SQLGetCommentByID(   redisDB, postID, scf.TimeZone)
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

// FlarumPosts 获取评论
func FlarumPosts(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	logger := ctx.GetLogger()
	h := ctx.h

	parm := r.URL.Query()
	_userID := parm.Get("filter[user]")
	_authorID := parm.Get("filter[author]")
	_disscussionID := parm.Get("filter[discussion]")
	// _type := parm.Get("filter[type]")
	_limit := parm.Get("page[limit]")
	_ids := parm.Get("filter[id]")
	// _sort := parm.Get("sort")
	_near := parm.Get("page[near]")
	_postID := ""

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
		user, err = model.SQLUserGet(h.App.GormDB, _userID)
		if user.ID == 0 || err != nil {
			h.flarumErrorJsonify(w, createSimpleFlarumError("Can't get the user for: "+_userID+err.Error()))
			return
		}

		rf = replyFilter{
			FT:    eUserPost,
			UID:   user.ID,
			Limit: limit,
		}

	} else if _authorID != "" {
		user, err = model.SQLUserGetByName(h.App.GormDB, _authorID)
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
		article, err := model.SQLArticleGetByID(h.App.GormDB, redisDB, aid)
		if err != nil {
			logger.Error("Can't get discussion id for ", aid)
			h.flarumErrorJsonify(w, createSimpleFlarumError("Can't get discussion for: "+_disscussionID+err.Error()))
			return
		}

		rf = replyFilter{
			FT:    eArticle,
			AID:   aid,
			Limit: article.CommentCount,
		}

		if _near != "" {
			near, err := strconv.ParseUint(_near, 10, 64)
			if err != nil {
				logger.Error("Can't get discussion id for ", _near)
				h.flarumErrorJsonify(w, createSimpleFlarumError("Can't get the page[near] for: "+_near+err.Error()))
				return
			}
			rf.NearNumber = near
		}

	} else if _ids != "" {
		postIds := strings.Split(_ids, ",")
		var _ids64 []uint64
		for _, _id := range postIds {
			_id64, err := strconv.ParseUint(_id, 10, 64)
			if err != nil {
				logger.Error("Can't get post id for", _id)
				continue
			}
			_ids64 = append(_ids64, _id64)
		}
		rf = replyFilter{
			FT:  ePosts,
			IDS: _ids64,
		}
	} else if _postID, err = h.safeGetParm(r, "cid"); _postID != "" && err == nil {
		postID, err := strconv.ParseUint(_postID, 10, 64)
		if err != nil {
			logger.Error("Can't get post id for ", _postID)
			h.flarumErrorJsonify(w, createSimpleFlarumError("Can't get postId for: "+err.Error()))
		}
		rf = replyFilter{
			FT:  ePost,
			CID: postID,
		}
	} else {
		logger.Warning("Can't process post api")
	}

	coreData, err = createFlarumPostAPIDoc(ctx, h.App.GormDB, redisDB, *h.App.Cf, rf, ctx.h.App.Cf.Site.TimeZone)
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
}

// FlarumPostsUtils 对于评论的一些操作
func FlarumPostsUtils(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := GetRetContext(r)
	logger := ctx.GetLogger()
	h := ctx.h
	_cid := pat.Param(r, "cid")
	cid, err := strconv.ParseUint(_cid, 10, 64)
	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("cid type error"))
		return
	}

	redisDB := h.App.RedisDB
	cobj, err := model.SQLCommentByID(h.App.GormDB, redisDB, cid, h.App.Cf.Site.TimeZone)
	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("无法获取评论"))
		return
	}

	// 用户所做的操作
	type CommentUtils struct {
		Data struct {
			ID         string `json:"id"`
			Type       string `json:"type"`
			Attributes map[string]interface{}
		} `json:"data"`
	}

	commentUtils := CommentUtils{}
	err = json.NewDecoder(r.Body).Decode(&commentUtils)
	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("json Decode err:"+err.Error()))
		return
	}

	if val, ok := commentUtils.Data.Attributes["isLiked"]; ok {
		cobj.DoLike(h.App.GormDB, redisDB, ctx.currentUser, val.(bool))
	}
	if val, ok := commentUtils.Data.Attributes["content"]; ok {
		logger.Errorf("Didn't apply the method to update the comment, Update content to ", val.(string))
		h.flarumErrorJsonify(w, createSimpleFlarumError("Didn't support modify comment"))
		return
	}

	rf := replyFilter{
		FT:  ePost,
		CID: cobj.ID,
		AID: cobj.AID,
	}

	coreData, err := createFlarumPostAPIDoc(ctx, h.App.GormDB, redisDB, *h.App.Cf, rf, ctx.h.App.Cf.Site.TimeZone)
	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("Get api doc error"+err.Error()))
		return
	}

	if ctx.inAPI {
		h.jsonify(w, coreData.APIDocument)
		return
	}

	h.flarumErrorJsonify(w, createSimpleFlarumError("此接口仅在API中使用"))
}
