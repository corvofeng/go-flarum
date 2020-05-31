package controller

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"goyoubbs/model"
	"goyoubbs/util"

	"github.com/ego008/youdb"
	"github.com/rs/xid"
	"goji.io/pat"
)

// ArticleAdd 添加新的帖子
func (h *BaseHandler) ArticleAdd(w http.ResponseWriter, r *http.Request) {
	cid := pat.Param(r, "cid")
	_, err := strconv.Atoi(cid)
	if err != nil {
		w.Write([]byte(`{"retcode":400,"retmsg":"cid type err"}`))
		return
	}

	currentUser, _ := h.CurrentUser(w, r)
	if currentUser.ID == 0 {
		w.Write([]byte(`{"retcode":401,"retmsg":"authored err"}`))
		return
	}
	if !currentUser.CanCreateTopic() {
		var msg string
		if currentUser.Flag == 1 {
			msg = "注册验证中，等待管理员通过"
		} else {
			msg = "您已被禁用"
		}
		w.Write([]byte(`{"retcode":401,"retmsg":"` + msg + `"}`))
		return
	}

	sqlDB := h.App.MySQLdb
	// db := h.App.Db

	cobj, err := model.SQLCategoryGetByID(sqlDB, cid)
	if err != nil {
		w.Write([]byte(`{"retcode":404,"retmsg":"` + err.Error() + `"}`))
		return
	}

	if cobj.Hidden {
		w.Write([]byte(`{"retcode":403,"retmsg":"category is Hidden"}`))
		return
	}

	type pageData struct {
		PageData
		Cobj model.Category

		// 可能是想获取当前节点的父节点下的所有子节点
		MainNodes []model.CategoryMini
	}

	tpl := h.CurrentTpl(r)
	evn := &pageData{}
	evn.SiteCf = h.App.Cf.Site
	evn.Title = "发表文章"
	evn.IsMobile = tpl == "mobile"
	evn.CurrentUser = currentUser
	evn.ShowSideAd = true
	evn.PageName = "article_add"

	evn.Cobj = cobj

	// 当前的主节点就直接从数据库中读取所有节点了
	evn.MainNodes, _ = model.SQLGetAllCategory(sqlDB)

	h.SetCookie(w, "token", xid.New().String(), 1)
	h.Render(w, tpl, evn, "layout.html", "articlecreate.html")
}

// ArticleAddPost 添加新帖, 文章预览接口
func (h *BaseHandler) ArticleAddPost(w http.ResponseWriter, r *http.Request) {
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
	if !currentUser.CanCreateTopic() || !currentUser.CanReply() {
		w.Write([]byte(`{"retcode":403,"retmsg":"user flag err"}`))
		return
	}

	type recForm struct {
		Act     string `json:"act"`
		CID     uint64 `json:"cid"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	decoder := json.NewDecoder(r.Body)
	var rec recForm
	err := decoder.Decode(&rec)
	if err != nil {
		w.Write([]byte(`{"retcode":400,"retmsg":"json Decode err:` + err.Error() + `"}`))
		return
	}
	defer r.Body.Close()

	rec.Title = strings.TrimSpace(rec.Title)
	rec.Content = strings.TrimSpace(rec.Content)

	sqlDB := h.App.MySQLdb
	db := h.App.Db
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
	if db.Hget("title_md5", []byte(titleMd5)).State == "ok" {
		w.Write([]byte(`{"retcode":403,"retmsg":"title has existed"}`))
		return
	}

	now := uint64(time.Now().UTC().Unix())
	scf := h.App.Cf.Site

	// 控制发帖间隔, 不能太短
	if !currentUser.IsAdmin() && currentUser.LastPostTime > 0 {
		if (now - currentUser.LastPostTime) < uint64(scf.PostInterval) {
			w.Write([]byte(`{"retcode":403,"retmsg":"PostInterval limited"}`))
			return
		}
	}

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

	cobj, err := model.SQLCategoryGetByID(sqlDB, strconv.FormatUint(rec.CID, 10))
	if err != nil {
		w.Write([]byte(`{"retcode":404,"retmsg":"` + err.Error() + `"}`))
		return
	}

	if cobj.Hidden {
		w.Write([]byte(`{"retcode":403,"retmsg":"category is Hidden"}`))
		return
	}

	newAid, _ := db.HnextSequence("article")
	aobj := model.Article{
		ArticleBase: model.ArticleBase{
			ID:  newAid,
			UID: currentUser.ID,
			CID: rec.CID,
		},
		Title:    rec.Title,
		Content:  rec.Content,
		AddTime:  now,
		EditTime: now,
		ClientIP: r.Header.Get("X-REAL-IP"),

		Active:        1, // 帖子为激活状态
		FatherTopicID: 0, // 没有原始主题
	}
	aobj.SQLCreateTopic(sqlDB)

	jb, _ := json.Marshal(aobj)
	aidB := youdb.I2b(newAid)
	db.Hset("article", aidB, jb)
	// 总文章列表
	db.Zset("article_timeline", aidB, aobj.EditTime)
	// 分类文章列表
	db.Zset("category_article_timeline:"+strconv.FormatUint(aobj.CID, 10), aidB, aobj.EditTime)
	// 用户文章列表
	db.Hset("user_article_timeline:"+strconv.FormatUint(aobj.UID, 10), youdb.I2b(aobj.ID), []byte(""))
	// 分类下文章数
	db.Zincr("category_article_num", youdb.I2b(aobj.CID), 1)

	currentUser.LastPostTime = now
	currentUser.Articles++

	jb, _ = json.Marshal(currentUser)
	db.Hset("user", youdb.I2b(aobj.UID), jb)

	// title md5
	db.Hset("title_md5", []byte(titleMd5), aidB)

	// send task work
	// get tag from title
	if scf.AutoGetTag && len(scf.GetTagApi) > 0 {
		db.Hset("task_to_get_tag", aidB, []byte(rec.Title))
	}

	// @ somebody in content
	sbs := util.GetMention(rec.Content,
		[]string{currentUser.Name, strconv.FormatUint(currentUser.ID, 10)})

	aid := strconv.FormatUint(newAid, 10)
	for _, sb := range sbs {
		var sbObj model.User
		sbu, err := strconv.ParseUint(sb, 10, 64)
		if err != nil {
			// @ user name
			sbObj, err = model.UserGetByName(db, strings.ToLower(sb))
		} else {
			// @ user id
			sbObj, err = model.UserGetByID(db, sbu)
		}

		if err == nil {
			if len(sbObj.Notice) > 0 {
				aidList := util.SliceUniqStr(strings.Split(aid+","+sbObj.Notice, ","))
				if len(aidList) > 100 {
					aidList = aidList[:100]
				}
				sbObj.Notice = strings.Join(aidList, ",")
				sbObj.NoticeNum = len(aidList)
			} else {
				sbObj.Notice = aid
				sbObj.NoticeNum = 1
			}
			jb, _ := json.Marshal(sbObj)
			db.Hset("user", youdb.I2b(sbObj.ID), jb)
		}
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

// IFeelLucky 随机的抽取一些帖子
func (h *BaseHandler) IFeelLucky(w http.ResponseWriter, r *http.Request) {

	var err error
	db := h.App.Db
	scf := h.App.Cf.Site
	redisDB := h.App.RedisDB
	logger := h.App.Logger

	type pageData struct {
		PageData
		SiteInfo model.SiteInfo
		PageInfo model.ArticlePageInfo
		Links    []model.Link
	}

	var count uint64
	sqlDB := h.App.MySQLdb
	// 获取全部的帖子数目
	err = sqlDB.QueryRow("SELECT COUNT(*) FROM topic").Scan(&count)
	if err != nil {
		logger.Debugf("Error %s", err)
		return
	}

	luckNum := uint64(scf.HomeShowNum)
	if luckNum > count {
		luckNum = count
	}

	articleList := make([]uint64, luckNum)

	func() {
		// 获得一些不重复的随机数
		m := make(map[uint64]uint64)
		for uint64(len(m)) < luckNum {
			n := uint64(h.App.Rand.Intn(int(count + 1)))
			if m[n] == n {
				continue
			} else {
				m[n] = n
			}
		}
		i := 0
		for _, v := range m {
			articleList[i] = v
			i++
		}
		sort.Slice(articleList, func(i, j int) bool { return articleList[i] < articleList[j] })
	}()
	logger.Debug("Get Article List", articleList)

	pageInfo := model.SQLArticleGetByList(sqlDB, db, redisDB, articleList, scf.TimeZone)
	categories, err := model.SQLGetAllCategory(sqlDB)

	tpl := h.CurrentTpl(r)
	evn := &pageData{}
	evn.SiteCf = scf
	evn.Title = scf.Name
	evn.Keywords = evn.Title
	evn.Description = scf.Desc
	evn.IsMobile = tpl == "mobile"
	currentUser, _ := h.CurrentUser(w, r)
	evn.CurrentUser = currentUser
	evn.ShowSideAd = false
	evn.PageName = "I feel lucky"
	evn.NewestNodes = categories
	// evn.HotNodes = model.CategoryHot(db, scf.CategoryShowNum)

	evn.SiteInfo = model.GetSiteInfo(redisDB)
	evn.PageInfo = pageInfo

	// 右侧的链接
	evn.Links = model.RedisLinkList(redisDB, false)

	h.Render(w, tpl, evn, "layout.html", "index.html")
}

// ArticleDetail 帖子的详情
func (h *BaseHandler) ArticleDetail(w http.ResponseWriter, r *http.Request) {
	var start uint64
	var err error

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
	aobj, err := model.SQLArticleGetByID(sqlDB, redisDB, aid)
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

	/*
		if len(currentUser.Notice) > 0 && len(currentUser.Notice) >= len(aid) {
			if len(aid) == len(currentUser.Notice) && aid == currentUser.Notice {
				currentUser.Notice = ""
				currentUser.NoticeNum = 0
				jb, _ := json.Marshal(currentUser)
				db.Hset("user", youdb.I2b(currentUser.ID), jb)
			} else {
				subStr := "," + aid + ","
				newNotice := "," + currentUser.Notice + ","
				if strings.Index(newNotice, subStr) >= 0 {
					currentUser.Notice = strings.Trim(strings.Replace(newNotice, subStr, "", 1), ",")
					currentUser.NoticeNum--
					jb, _ := json.Marshal(currentUser)
					db.Hset("user", youdb.I2b(currentUser.ID), jb)
				}
			}
		}
	*/

	if aobj.Hidden && !currentUser.IsAdmin() {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"retcode":404,"retmsg":"not found"}`))
		return
	}

	// 获取帖子所在的节点
	cobj, err := model.SQLCategoryGetByID(sqlDB, strconv.FormatUint(aobj.CID, 10))

	err = nil
	cobj.Hidden = false
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	if cobj.Hidden && !currentUser.IsAdmin() {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"retcode":404,"retmsg":"not found"}`))
		return
	}

	// Authorized
	if scf.Authorized && currentUser.Flag < 5 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"retcode":401,"retmsg":"Unauthorized"}`))
		return
	}

	// if btn == "prev" {
	// 	start = start - uint64(scf.HomeShowNum) - 1
	// }

	// cobj.Articles = db.Zget("category_article_num", youdb.I2b(cobj.ID)).Uint64()
	pageInfo := model.SQLCommentList(
		sqlDB,
		redisDB,
		aobj.ID,
		start,
		btn,
		scf.CommentListNum,
		scf.TimeZone,
	)

	type articleForDetail struct {
		model.Article
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
		PageData
		Aobj     articleForDetail
		Author   model.User
		Cobj     model.Category
		Relative model.ArticleRelative
		PageInfo model.CommentPageInfo
		Views    uint64
		SiteInfo model.SiteInfo
	}

	tpl := h.CurrentTpl(r)
	evn := &pageData{}
	evn.SiteCf = scf
	evn.Title = aobj.Title + " - " + cobj.Name + " - " + scf.Name
	// evn.Keywords = aobj.Tags
	// evn.Description = cobj.Name + " - " + aobj.Title + " - " + aobj.Tags
	evn.IsMobile = tpl == "mobile"

	evn.CurrentUser = currentUser
	evn.ShowSideAd = true
	evn.PageName = "article_detail"
	// evn.HotNodes = model.CategoryHot(db, scf.CategoryShowNum)
	// evn.NewestNodes = model.CategoryNewest(db, scf.CategoryShowNum)

	author, _ := model.SQLUserGetByID(sqlDB, aobj.UID)

	if author.ID == 2 {
		// 这部分的网页是转载而来的, 所以需要保持原样式, 这里要牺牲XSS的安全性了
		evn.Aobj = articleForDetail{
			Article:     aobj,
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
			Article:     aobj,
			ContentFmt:  template.HTML(util.ContentFmt(aobj.Content)),
			CommentsCnt: commentsCnt,
			Name:        author.Name,
			Avatar:      author.Avatar,
			Views:       aobj.ClickCnt,
			AddTimeFmt:  util.TimeFmt(aobj.AddTime, "2006-01-02 15:04", scf.TimeZone),
			EditTimeFmt: util.TimeFmt(aobj.EditTime, "2006-01-02 15:04", scf.TimeZone),
		}
	}

	if len(aobj.Tags) > 0 {
		var tags []string
		for _, v := range strings.Split(aobj.Tags, ",") {
			tags = append(tags, `<a href="/tag/`+v+`">`+v+`</a>`)
		}
		evn.Aobj.TagStr = template.HTML(strings.Join(tags, ", "))
	}

	evn.Cobj = cobj
	evn.Author = author
	evn.Relative = model.ArticleGetRelative(aobj.ID, aobj.Tags)
	evn.PageInfo = pageInfo
	evn.SiteInfo = model.GetSiteInfo(redisDB)

	token := h.GetCookie(r, "token")
	if len(token) == 0 {
		token := xid.New().String()
		h.SetCookie(w, "token", token, 1)
	}

	if h.InAPI {
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

// ArticleDetailPost 负责处理用户的评论
// 原有的程序中帖子与评论一起保存在文件中, 经过修改后
// 当前的所有帖子与评论均保存在数据库中, 评论数量可以保存在文件中
func (h *BaseHandler) ArticleDetailPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	token := h.GetCookie(r, "token")
	if len(token) == 0 {
		w.Write([]byte(`{"retcode":400,"retmsg":"token cookie missed"}`))
		return
	}

	_aid := pat.Param(r, "aid")
	aid, err := strconv.ParseUint(_aid, 10, 64)
	if err != nil {
		w.Write([]byte(`{"retcode":400,"retmsg":"aid type err:` + err.Error() + `"}`))
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
	err = decoder.Decode(&rec)
	if err != nil {
		w.Write([]byte(`{"retcode":400,"retmsg":"json Decode err:` + err.Error() + `"}`))
		return
	}
	defer r.Body.Close()

	sqlDB := h.App.MySQLdb
	cntDB := h.App.Db
	redisDB := h.App.RedisDB
	rsp := response{}

	if rec.Act == "link_click" {
		rsp.Retcode = 200
		if len(rec.Link) > 0 {
			hash := md5.Sum([]byte(rec.Link))
			urlMd5 := hex.EncodeToString(hash[:])
			bn := "article_detail_token"
			clickKey := []byte(token + ":click:" + urlMd5)
			if cntDB.Zget(bn, clickKey).State == "ok" {
				w.Write([]byte(`{"retcode":403,"retmsg":"err"}`))
				return
			}
			cntDB.Zset(bn, clickKey, uint64(time.Now().UTC().Unix()))
			cntDB.Hincr("url_md5_click", []byte(urlMd5), 1)

			w.Write([]byte(`{"retcode":200,"retmsg":"ok"}`))
			return
		}
	} else if rec.Act == "comment_preview" {
		rsp.Retcode = 200
		rsp.Html = template.HTML(util.ContentFmt(rec.Content))
	} else if rec.Act == "comment_submit" {
		timeStamp := uint64(time.Now().UTC().Unix())
		currentUser, _ := h.CurrentUser(w, r)
		if !currentUser.CanReply() {
			w.Write([]byte(`{"retcode":403,"retmsg":"forbidden"}`))
			return
		}
		if (timeStamp - currentUser.LastReplyTime) < uint64(h.App.Cf.Site.CommentInterval) {
			w.Write([]byte(`{"retcode":403,"retmsg":"out off comment interval"}`))
			return
		}
		// 获取当前的话题
		aobj, err := model.SQLArticleGetByID(sqlDB, redisDB, aid)
		if err != nil {
			w.Write([]byte(`{"retcode":404,"retmsg":"not found"}`))
			return
		}
		if aobj.CloseComment {
			w.Write([]byte(`{"retcode":403,"retmsg":"comment forbidden"}`))
			return
		}
		// commentID, _ := cntDB.HnextSequence("article_comment:" + aid)
		obj := model.Comment{
			Aid:      aobj.ID,
			UID:      currentUser.ID,
			Content:  rec.Content,
			AddTime:  timeStamp,
			ClientIp: r.Header.Get("X-REAL-IP"),
		}

		obj.SQLSaveComment(sqlDB)
		// jb, _ := json.Marshal(obj)

		// cntDB.Hset("article_comment:"+aid, youdb.I2b(obj.ID), jb) // 文章评论bucket
		// cntDB.Hincr("count", []byte("comment_num"), 1)            // 评论总数
		// // 用户回复文章列表
		// cntDB.Zset("user_article_reply:"+strconv.FormatUint(obj.UID, 10), youdb.I2b(obj.Aid), obj.AddTime)

		// 更新文章列表时间

		// aobj.Comments = commentID
		aobj.RUID = currentUser.ID
		aobj.EditTime = timeStamp
		jb2, _ := json.Marshal(aobj)
		cntDB.Hset("article", youdb.I2b(aobj.ID), jb2)

		currentUser.LastReplyTime = timeStamp
		currentUser.Replies++
		jb3, _ := json.Marshal(currentUser)
		cntDB.Hset("user", youdb.I2b(currentUser.ID), jb3)

		// 总文章列表
		cntDB.Zset("article_timeline", youdb.I2b(aobj.ID), timeStamp)
		// 分类文章列表
		cntDB.Zset("category_article_timeline:"+strconv.FormatUint(aobj.CID, 10), youdb.I2b(aobj.ID), timeStamp)

		// // @ somebody in comment & topic author
		// sbs := util.GetMention("@"+strconv.FormatUint(aobj.UID, 10)+" "+rec.Content,
		// 	[]string{currentUser.Name, strconv.FormatUint(currentUser.ID, 10)})
		// for _, sb := range sbs {
		// 	var sbObj model.User
		// 	sbu, err := strconv.ParseUint(sb, 10, 64)
		// 	if err != nil {
		// 		// @ user name
		// 		sbObj, err = model.UserGetByName(cntDB, strings.ToLower(sb))
		// 	} else {
		// 		// @ user id
		// 		sbObj, err = model.UserGetByID(cntDB, sbu)
		// 	}

		// 	if err == nil {
		// 		if len(sbObj.Notice) > 0 {
		// 			aidList := util.SliceUniqStr(strings.Split(aid+","+sbObj.Notice, ","))
		// 			if len(aidList) > 100 {
		// 				aidList = aidList[:100]
		// 			}
		// 			sbObj.Notice = strings.Join(aidList, ",")
		// 			sbObj.NoticeNum = len(aidList)
		// 		} else {
		// 			sbObj.Notice = aid
		// 			sbObj.NoticeNum = 1
		// 		}
		// 		jb, _ := json.Marshal(sbObj)
		// 		cntDB.Hset("user", youdb.I2b(sbObj.ID), jb)
		// 	}
		// }

		rsp.Retcode = 200
	}

	json.NewEncoder(w).Encode(rsp)
}

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
