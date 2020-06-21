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

// 获取一条评论的信息
func createFlarumCommentAPIDoc(
	logger *logging.Logger, sqlDB *sql.DB, redisDB *redis.Client,
	appConf model.AppConf,
	siteInfo model.SiteInfo,
	currentUser *model.User,
	inAPI bool,
	cid uint64, tz int,
) (flarum.CoreData, error) {
	var err error
	coreData := flarum.CoreData{}
	apiDoc := &coreData.APIDocument

	comment, err := model.SQLGetCommentByID(sqlDB, redisDB, cid, tz)
	if err != nil {
		logger.Error("Get comment error", err)
		return coreData, err
	}
	commentListItem := model.CommentListItem{Comment: comment}

	article, err := model.SQLArticleGetByID(sqlDB, redisDB, comment.AID)
	if err != nil {
		logger.Error("Get article error", err)
		return coreData, err
	}

	diss := model.FlarumCreateDiscussionFromArticle(article)
	post := model.FlarumCreatePost(commentListItem)
	apiDoc.SetData(post)
	// apiDoc.AppendResourcs(post)

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
		h.flarumErrorJsonify(w, createSimpleFlarumError("解析json错误:"+err.Error()))
		return
	}
	aid, err := strconv.ParseUint(reply.Data.Relationships.Discussion.Data.ID, 10, 64)
	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("无法获取正确的帖子信息:"+err.Error()))
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
		h.flarumErrorJsonify(w, createSimpleFlarumError("创建评论出现错误:"+err.Error()))
		return

	}

	coreData, err := createFlarumCommentAPIDoc(logger, sqlDB, redisDB, *h.App.Cf, si, ctx.currentUser, ctx.inAPI, comment.ID, scf.TimeZone)
	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("查询评论出现错误:"+err.Error()))
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
