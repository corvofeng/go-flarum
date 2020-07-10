package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"math"
	"strconv"
	"strings"
	"time"

	"goyoubbs/model/flarum"
	"goyoubbs/util"
	"html/template"

	"github.com/go-redis/redis/v7"

	"github.com/ego008/youdb"
)

// ArticleBase 基础的文档类, 在数据库表中的字段
type ArticleBase struct {
	ID  uint64 `json:"id"`
	UID uint64 `json:"uid"`
	CID uint64 `json:"cid"`

	Title   string `json:"title"`
	Content string `json:"content"`

	FirstPostID uint64
	LastPostID  uint64

	ClickCnt uint64 `json:"clickcnt"` // 不保证精确性
	Comments uint64 `json:"comments"`

	AddTime  uint64 `json:"addtime"`
	EditTime uint64 `json:"edittime"`

	ClientIP string `json:"clientip"`
}

// Article store in database
type Article struct {
	ArticleBase
	RUID              uint64 `json:"ruid"`
	Tags              string `json:"tags"`
	AnonymousComments bool   `json:"annoymous_comments"` // 是否允许匿名评论
	CloseComment      bool   `json:"closecomment"`
	Hidden            bool   `json:"hidden"`    // Depreacte, do not use it.
	StickTop          bool   `json:"stick_top"` // 是否置顶

	// 帖子被管理员修改后, 已经保存的旧的帖子ID
	FatherTopicID uint64 `json:"fathertopicid"`

	// 记录当前帖子是否可以被用户看到, 与上面的hidden类似
	Active uint64 `json:"active"`
}

// ArticleMini 缩略版的Article信息
type ArticleMini struct {
	ArticleBase
	Ruid   uint64 `json:"ruid"`
	Hidden bool   `json:"hidden"`
}

// ArticleListItem data strucy only used in page.
type ArticleListItem struct {
	ArticleBase
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	Cname       string `json:"cname"`
	Ruid        uint64 `json:"ruid"`
	Rname       string `json:"rname"`
	EditTimeFmt string `json:"edittimefmt"`
	ClickCnt    uint64 `json:"clickcnt"`

	/**
	 * When in search page, every article item have the highlight content,
	 * we need tell users that we do not return a random list.
	 */
	HighlightContent template.HTML `json:"highlight_content"`

	LastComment *CommentListItem
}

// ArticlePageInfo data in every list page
type ArticlePageInfo struct {
	Items      []ArticleListItem `json:"items"`
	HasPrev    bool              `json:"hasprev"`
	HasNext    bool              `json:"hasnext"`
	FirstKey   uint64            `json:"firstkey"`
	FirstScore uint64            `json:"firstscore"`
	LastKey    uint64            `json:"lastkey"`
	PageNum    uint64            `json:"pagenum"`
	PagePrev   uint64            `json:"pageprev"`
	PageNext   uint64            `json:"pagenext"`
	LastScore  uint64            `json:"lastscore"`
}

// FlarumArticlePageInfo flarum站点的数据信息
type FlarumArticlePageInfo struct {
	Items     []flarum.Discussion
	LinkFirst string
	LinkPrev  string
	LinkNext  string
}

type ArticleLi struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
	Tags  string `json:"tags"`
}

type ArticleRelative struct {
	Articles []ArticleLi
	Tags     []string
}

// ArticleFeedListItem rss资源
type ArticleFeedListItem struct {
	ID          uint64
	UID         uint64
	Name        string
	Cname       string
	Title       string
	AddTimeFmt  string
	EditTimeFmt string
	Des         string
}

// ArticleTag 文章添加、编辑后传给后台任务的信息
// TODO: delete
type ArticleTag struct {
	ID      uint64
	OldTags string
	NewTags string
}

// SQLArticleGetByID 通过 article id获取内容
func SQLArticleGetByID(db *sql.DB, redisDB *redis.Client, aid uint64) (Article, error) {
	articleBaseList := sqlGetArticleBaseByList(db, redisDB, []uint64{aid})
	obj := Article{}
	if len(articleBaseList) == 0 {
		return obj, errors.New("No result")
	}
	obj.ArticleBase = articleBaseList[0]
	obj.GetArticleCntFromRedisDB(db, redisDB)

	return obj, nil
}

// GetCommentsSize 获取评论
/*
 * db (*sql.DB): TODO
 * cntDB (*youdb.DB): TODO
 * redisDB (redis.Client): TODO
 */
func (article *Article) GetCommentsSize(db *sql.DB) uint64 {
	return article.Comments
}

// GetWeight 获取当前帖子的权重
/**
 * (Log10(QView) * 2 + 4 * comments)/ QAge
 *
 * db (*sql.DB): TODO
 * redisDB (redis.Client): TODO
 */
func (article *Article) GetWeight(db *sql.DB, redisDB *redis.Client) float64 {
	var editTime time.Time
	var now = time.Now()
	if article.EditTime == 0 {
		editTime = time.Unix(now.Unix()-2*24*3600, 0)
	} else {
		editTime = time.Unix(int64(article.EditTime), 0)
		sTime := time.Unix(1577836800, 0) // 2020-01-01 08:00:00

		if editTime.Before(sTime) {
			editTime = time.Unix(now.Unix()-2*24*3600, 0)
		}
	}
	if article.ClickCnt == 0 { // 避免出现0的情况
		article.ClickCnt = 1
	}
	qAge := now.Sub(editTime).Hours()
	weight := (math.Log10(float64(article.ClickCnt))*2 + 4*float64(article.GetCommentsSize(db))) / (qAge * 1.0)
	return weight
}

func (article *Article) sqlCreateTopic(tx *sql.Tx) (bool, error) {
	row, err := tx.Exec(
		("INSERT INTO `topic` " +
			" (`node_id`, `user_id`, `title`, `content`, created_at, updated_at, client_ip, reply_count, active)" +
			" VALUES " +
			" (?, ?, ?, ?, ?, ?, ?, 0, ?)"),
		article.CID,
		article.UID,
		article.Title,
		article.Content,
		article.AddTime,
		article.EditTime,
		article.ClientIP,
		article.Active,
	)
	if err != nil {
		return false, err
	}
	aid, err := row.LastInsertId()
	article.ID = uint64(aid)
	return true, nil
}

// SQLCreateTopic 创建主题
func (article *Article) SQLCreateTopic(db *sql.DB) bool {
	tx, err := db.Begin()
	defer clearTransaction(tx)
	if err != nil {
		return false
	}
	article.sqlCreateTopic(tx)
	if err := tx.Commit(); err != nil {
		logger := util.GetLogger()
		logger.Error("Create topic with error", err)
		return false
	}
	return true
}

// CreateFlarumDiscussion 创建flarum的帖子
func (article *Article) CreateFlarumDiscussion(db *sql.DB, tags flarum.RelationArray) (bool, error) {
	tx, err := db.Begin()
	logger := util.GetLogger()
	defer clearTransaction(tx)
	if err != nil {
		return false, err
	}
	if ok, err := article.sqlCreateTopic(tx); !ok {
		return false, err
	}
	logger.Debugf("Create article %d success", article.ID)

	comment := Comment{
		CommentBase: CommentBase{
			AID:      article.ID,
			UID:      article.UID,
			Content:  article.Content,
			Number:   1,
			ClientIP: article.ClientIP,
			AddTime:  article.AddTime,
		},
	}
	if ok, err := comment.sqlSaveComment(tx); !ok {
		return false, err
	}
	article.LastPostID = comment.ID
	article.FirstPostID = comment.ID
	if ok, err := article.updateFlarumPost(tx); !ok {
		return false, err
	}

	if ok, err := article.updateFlarumTag(tx, tags); !ok {
		return false, err
	}
	logger.Debugf("Update article post and tag %d success", article.ID)

	if err := tx.Commit(); err != nil {
		logger := util.GetLogger()
		logger.Error("Create topic with error", err)
		return false, err
	}

	return true, nil
}

// updateFlarumPost 更新评论信息
func (article *Article) updateFlarumPost(tx *sql.Tx) (bool, error) {

	_, err := tx.Exec(
		"UPDATE `topic` SET"+
			" first_post_id=?,"+
			" last_post_id=?"+
			" where id=?",
		article.FirstPostID,
		article.LastPostID,
		article.ID,
	)
	if util.CheckError(err, "更新帖子") {
		return false, err
	}
	return true, nil
}

// updateFlarumTag 更新评论信息
func (article *Article) updateFlarumTag(tx *sql.Tx, tags flarum.RelationArray) (bool, error) {
	for _, rela := range tags.Data {
		_, err := tx.Exec(
			("INSERT INTO `topic_tag` " +
				" (`topic_id`, `tag_id`)" +
				" VALUES " +
				" (?, ?)"),
			article.ID,
			rela.ID,
		)
		if util.CheckError(err, "更新帖子") {
			return false, err
		}
	}
	return true, nil
}

// SQLArticleUpdate 更新当前帖子
func (article *Article) SQLArticleUpdate(db *sql.DB, redisDB *redis.Client) bool {
	// 更新记录必须要被保存, 配合数据库中的father_topic_id来实现
	// 每次更新主题, 会将前帖子复制为一个新帖子(active=0不被看见),
	// 当前帖子的id没有变化, 但是father_topic_id变为这个新的帖子.
	// 通过father_topic_id组成了链表的关系

	// 以当前帖子为模板创建一个新的帖子
	// 对象中只有简单的数据结构, 浅拷贝即可, 需要将其设为不可见
	oldArticle, err := SQLArticleGetByID(db, redisDB, article.ID)
	oldArticle.Active = 0
	if util.CheckError(err, "修改时拷贝") {
		return false
	}
	oldArticle.SQLCreateTopic(db)

	_, err = db.Exec(
		"UPDATE `topic` "+
			"set title=?,"+
			"content = ?,"+
			"node_id=?,"+
			"user_id=?,"+
			"updated_at=?,"+
			"client_ip=?,"+
			"father_topic_id = ?"+
			" where id=?",
		article.Title,
		article.Content,
		article.CID,
		article.UID,
		article.EditTime,
		article.ClientIP,
		oldArticle.ID,
		article.ID,
	)
	if util.CheckError(err, "更新帖子") {
		return false
	}

	return true
}

func (ab *ArticleBase) ToArticleListItem(sqlDB *sql.DB, redisDB *redis.Client, tz int) ArticleListItem {
	item := ArticleListItem{
		ArticleBase: *ab,
	}
	item.EditTimeFmt = util.TimeFmt(item.EditTime, util.TIME_FMT, tz)
	item.Cname = GetCategoryNameByCID(sqlDB, redisDB, item.CID)
	if item.LastPostID != 0 {
		lastComment, err := SQLGetCommentByID(sqlDB, redisDB, item.LastPostID, tz)
		if err != nil {
			util.GetLogger().Errorf("Can't get last comment(%d)for article(%d)", item.LastPostID, item.ID)
		} else {
			lc := lastComment.toCommentListItem(sqlDB, redisDB, tz)
			item.LastComment = &lc
		}
	}

	return item
}

// SQLArticleGetByList 通过id列表获取对应的帖子
func SQLArticleGetByList(db *sql.DB, redisDB *redis.Client, articleList []uint64, tz int) ArticlePageInfo {
	var items []ArticleListItem
	var hasPrev, hasNext bool
	var firstKey, firstScore, lastKey, lastScore uint64
	articleBaseList := sqlGetArticleBaseByList(db, redisDB, articleList)
	m := make(map[uint64]ArticleListItem)

	for _, articleBase := range articleBaseList {
		m[articleBase.ID] = articleBase.ToArticleListItem(db, redisDB, tz)
	}

	for _, id := range articleList {
		if item, ok := m[id]; ok {
			items = append(items, item)
		}
	}
	return ArticlePageInfo{
		Items:      items,
		HasPrev:    hasPrev,
		HasNext:    hasNext,
		FirstKey:   firstKey,
		FirstScore: firstScore,
		LastKey:    lastKey,
		LastScore:  lastScore,
	}
}

// sqlGetArticleBaseByList 获取帖子信息, NOTE: 请尽量调用该函数, 而不是自己去写sql语句
func sqlGetArticleBaseByList(db *sql.DB, redisDB *redis.Client, articleList []uint64) (items []ArticleBase) {
	var err error
	var rows *sql.Rows
	var articleListStr []string
	logger := util.GetLogger()
	defer rowsClose(rows)

	if len(articleList) == 0 {
		logger.Warning("SQLArticleGetByList: Can't process the article list empty")
		return
	}

	for _, v := range articleList {
		articleListStr = append(articleListStr, strconv.FormatInt(int64(v), 10))
	}
	qFieldList := []string{
		"id", "title", "content",
		"node_id", "user_id", "hits", "reply_count",
		"first_post_id", "last_post_id",
		"created_at", "updated_at",
	}
	sql := fmt.Sprintf("select %s from topic where id in (%s)",
		strings.Join(qFieldList, ","),
		strings.Join(articleListStr, ","))

	rows, err = db.Query(sql)
	if err != nil {
		logger.Errorf("Query failed,err:%v", err)
		return
	}
	m := make(map[uint64]ArticleBase)
	for rows.Next() {
		item := ArticleBase{}
		err = rows.Scan(
			&item.ID, &item.Title, &item.Content,
			&item.CID, &item.UID, &item.ClickCnt, &item.Comments,
			&item.FirstPostID, &item.LastPostID,
			&item.AddTime, &item.EditTime)

		if err != nil {
			logger.Errorf("Scan failed,err:%v", err)
			continue
		}
		item.ClickCnt = GetArticleCntFromRedisDB(db, redisDB, item.ID)
		m[item.ID] = item
	}

	for _, id := range articleList {
		if item, ok := m[id]; ok {
			items = append(items, item)
		}
	}

	return
}

// SQLArticleGetByCID 根据页码获取某个分类的列表
func SQLArticleGetByCID(db *sql.DB, redisDB *redis.Client, nodeID, page, limit uint64, tz int) ArticlePageInfo {
	var pageInfo ArticlePageInfo
	articleList := GetTopicListByPageNum(nodeID, page, limit)
	logger := util.GetLogger()
	logger.Debug("Get article list", page, limit, articleList)
	if len(articleList) == 0 {
		// TODO: remove it
		articleIteratorStart := GetCIDArticleMax(nodeID)
		pageInfo = SQLCIDArticleList(db, redisDB, nodeID, articleIteratorStart, "next", limit, tz)
		// 先前没有缓存, 需要加入到rank map中
		var items []ArticleRankItem
		for _, a := range pageInfo.Items {
			items = append(items, ArticleRankItem{
				AID:     a.ID,
				SQLDB:   db,
				RedisDB: redisDB,
				Weight:  a.ClickCnt,
			})
		}
		AddNewArticleList(nodeID, items)
	} else {
		pageInfo = SQLArticleGetByList(db, redisDB, articleList, tz)
	}
	pageInfo.PageNum = page
	pageInfo.PageNext = page + 1
	pageInfo.PagePrev = page - 1
	if len(articleList) == int(limit) {
		pageInfo.HasNext = true
	}
	return pageInfo
}

// SQLArticleGetByUID 根据创建用户获取帖子列表
func SQLArticleGetByUID(db *sql.DB, redisDB *redis.Client, uid, page, limit uint64, tz int) ArticlePageInfo {
	var rows *sql.Rows
	var err error
	var pageInfo ArticlePageInfo
	var articleList []uint64
	logger := util.GetLogger()
	defer rowsClose(rows)

	rows, err = db.Query(
		"SELECT id FROM topic WHERE user_id = ? and active = 1 ORDER BY created_at DESC limit ? offset ?",
		uid, limit, (page-1)*limit,
	)
	if err != nil {
		logger.Errorf("Query failed,err:%v", err)
		return pageInfo
	}
	for rows.Next() {
		var aid uint64
		err = rows.Scan(&aid) //不scan会导致连接不释放
		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			continue
		}
		articleList = append(articleList, aid)
	}

	logger.Debug("Get article list", page, limit, articleList)
	pageInfo = SQLArticleGetByList(db, redisDB, articleList, tz)

	pageInfo.PageNum = page
	pageInfo.PageNext = page + 1
	pageInfo.PagePrev = page - 1
	if len(articleList) == int(limit) {
		pageInfo.HasNext = true
	}
	return pageInfo
}

// SQLArticleSetClickCnt 更新每个帖子的权重, 用于将redis中的数据同步过去
func SQLArticleSetClickCnt(sqlDB *sql.DB, aid uint64, clickCnt uint64) {
	_, err := sqlDB.Exec("UPDATE `topic`"+
		" set hits = ?"+
		" where id = ?",
		clickCnt,
		aid,
	)
	util.CheckError(err, "更新帖子点击次数")
}

// SQLArticleSetCommentCnt 更新每个帖子的权重, 用于将redis中的数据同步过去
func SQLArticleSetCommentCnt(sqlDB *sql.DB, aid uint64, replyCnt uint64) {
	_, err := sqlDB.Exec("UPDATE `topic`"+
		" set reply_count = ?"+
		" where id = ?",
		replyCnt,
		aid,
	)
	util.CheckError(err, "更新帖子评论数目")
}

// SQLCIDArticleList 返回某个节点的主题
// nodeID 为0 表示全部主题
// TODO: delete it
func SQLCIDArticleList(db *sql.DB, redisDB *redis.Client, nodeID, start uint64, btnAct string, limit uint64, tz int) ArticlePageInfo {
	var hasPrev, hasNext bool
	var firstKey, firstScore, lastKey, lastScore uint64
	var rows *sql.Rows
	var err error
	logger := util.GetLogger()
	valueList := "id"
	selectList := " active != 0 "
	var articleList []uint64
	if nodeID == 0 {
		if btnAct == "" || btnAct == "next" {
			rows, err = db.Query(
				"SELECT "+valueList+" FROM topic WHERE id > ? and "+selectList+" ORDER BY id limit ?",
				start, limit,
			)
		} else if btnAct == "prev" {
			rows, err = db.Query(
				"SELECT * FROM (SELECT "+valueList+" FROM topic WHERE id < ? and "+selectList+" ORDER BY id DESC limit ?) as t ORDER BY id",
				start, limit,
			)
		} else {
			logger.Error("Get wrond button action")
		}
	} else {
		if btnAct == "" || btnAct == "next" {
			rows, err = db.Query(
				"SELECT "+valueList+" FROM topic WHERE node_id = ? And id > ? and "+selectList+" ORDER BY id limit ?",
				nodeID, start, limit,
			)
		} else if btnAct == "prev" {
			rows, err = db.Query(
				"SELECT * FROM (SELECT "+valueList+" FROM topic WHERE node_id = ? and id < ? and "+selectList+" ORDER BY id DESC limit ?) as t ORDER BY id",
				nodeID, start, limit,
			)
		} else {
			logger.Error("Get wrond button action")
		}
	}

	defer rowsClose(rows)
	if err != nil {
		logger.Errorf("Query failed,err:%v", err)
		return ArticlePageInfo{}
	}

	for rows.Next() {
		var aid uint64
		err = rows.Scan(&aid) //不scan会导致连接不释放
		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			continue
		}
		articleList = append(articleList, aid)
	}

	pageInfo := SQLArticleGetByList(db, redisDB, articleList, tz)

	if len(articleList) > 0 {
		firstKey = articleList[0]
		lastKey = articleList[len(articleList)-1]
		hasNext = true
		hasPrev = true

		// 前一页, 后一页的判断其实比较复杂的, 这里只针对最容易的情况进行了判断,
		// 因为帖子较多时, 这算是一种近似

		// 如果最开始的帖子ID为1, 那肯定是没有了前一页了
		if articleList[0] == 1 || start < uint64(limit) {
			hasPrev = false
		}

		// 查询出的数量比要求的数量要少, 说明没有下一页
		if uint64(len(articleList)) < limit {
			hasNext = false
		}
	}

	return ArticlePageInfo{
		Items:      pageInfo.Items,
		HasPrev:    hasPrev,
		HasNext:    hasNext,
		FirstKey:   firstKey,
		FirstScore: firstScore,
		LastKey:    lastKey,
		LastScore:  lastScore,
	}
}

// only for rank
func sqlGetAllArticleWithCID(db *sql.DB, cid uint64, active bool) ([]ArticleMini, error) {
	var articles []ArticleMini
	var rows *sql.Rows
	var err error
	logger := util.GetLogger()

	activeData := 0

	if active {
		activeData = 1
	} else {
		activeData = 0
	}

	if cid == 0 {
		// cid为0, 查询所有节点
		rows, err = db.Query(
			"SELECT t_list.topic_id FROM (SELECT topic_id FROM `topic_tag`) as t_list LEFT JOIN topic ON t_list.topic_id = topic.id WHERE active = ?",
			activeData)
	} else {
		rows, err = db.Query(
			"SELECT t_list.topic_id FROM (SELECT topic_id FROM `topic_tag` where tag_id = ?) as t_list LEFT JOIN topic ON t_list.topic_id = topic.id WHERE active = ?",
			cid, activeData)
	}
	defer rowsClose(rows)

	if err != nil {
		logger.Error("Can't get topic_tag info", err)
		return articles, err
	}
	for rows.Next() {
		obj := ArticleMini{}
		err = rows.Scan(&obj.ID)
		if err != nil {
			logger.Error("Scan failed", err)
			return articles, err
		}

		articles = append(articles, obj)
	}

	return articles, nil
}

// IncrArticleCntFromRedisDB 增加点击次数
func (article *Article) IncrArticleCntFromRedisDB(sqlDB *sql.DB, redisDB *redis.Client) uint64 {
	var clickCnt uint64 = 0
	aid := article.ID
	rep := redisDB.HGet("article_views", fmt.Sprintf("%d", aid))
	_, err := rep.Uint64()
	if err == redis.Nil { // 只有当redis中的数据不存在时，才向mysql与内存数据库请求
		if clickCnt == 0 {
			// 首先从sqlDB中查找
			// fmt.Println("Get data from sqlDB", aid)
			rows, err := sqlDB.Query("SELECT hits FROM topic where id = ?", aid)
			defer func() {
				if rows != nil {
					rows.Close()
				}
			}()
			if err != nil {
				fmt.Printf("Query failed,err:%v", err)
				clickCnt = 0
			} else {
				for rows.Next() {
					err = rows.Scan(&clickCnt)
					if err != nil {
						fmt.Printf("Scan failed,err:%v", err)
					}

				}
			}
		}
		if clickCnt == 0 {
			fmt.Println("Get data from cntDB", aid)
			redisDB.HSet("article_views", fmt.Sprintf("%d", aid), clickCnt)
		}
		if clickCnt == 0 {
			rep := redisDB.HIncrBy("article_views", fmt.Sprintf("%d", aid), 1)
			clickCnt = uint64(rep.Val())
		}
	} else {
		rep := redisDB.HIncrBy("article_views", fmt.Sprintf("%d", aid), 1)
		clickCnt = uint64(rep.Val())
	}

	return clickCnt
}

// GetArticleCntFromRedisDB 获取当前帖子的点击次数
/*
 * sqlDB (*sql.DB): TODO
 * redisDB (*redis.Client): TODO
 */
func (article *Article) GetArticleCntFromRedisDB(sqlDB *sql.DB, redisDB *redis.Client) uint64 {
	article.ClickCnt = GetArticleCntFromRedisDB(sqlDB, redisDB, article.ID)
	return article.ClickCnt
}

// GetArticleCntFromRedisDB 从不同的数据库中获取点击数
func GetArticleCntFromRedisDB(sqlDB *sql.DB, redisDB *redis.Client, aid uint64) uint64 {
	rep := redisDB.HGet("article_views", fmt.Sprintf("%d", aid))
	data, err := rep.Uint64()

	if err != nil && err != redis.Nil {
		fmt.Printf("Get %d with error :%v", aid, err)
		data = 0
	}
	return data
}

// SQLArticleList 返回所有节点的主题
func SQLArticleList(db *sql.DB, redisDB *redis.Client, start uint64, btnAct string, limit uint64, tz int) ArticlePageInfo {
	return SQLCIDArticleList(
		db, redisDB, 0, start, btnAct, limit, tz,
	)
}

// ArticleFeedList 旧有函数, TODO: 增加feeds功能
func ArticleFeedList(db *youdb.DB, limit, tz int) []ArticleFeedListItem {
	var items []ArticleFeedListItem
	var keys [][]byte
	keyStart := []byte("")

	for {
		rs := db.Hrscan("article", keyStart, limit)
		if rs.State != "ok" {
			break
		}
		for i := 0; i < (len(rs.Data) - 1); i += 2 {
			keyStart = rs.Data[i]
			keys = append(keys, rs.Data[i])
		}

		if len(keys) > 0 {
			var aitems []Article
			userMap := map[uint64]UserMini{}
			categoryMap := map[uint64]CategoryMini{}

			rs := db.Hmget("article", keys)
			if rs.State == "ok" {
				for i := 0; i < (len(rs.Data) - 1); i += 2 {
					item := Article{}
					json.Unmarshal(rs.Data[i+1], &item)
					if !item.Hidden {
						aitems = append(aitems, item)
						userMap[item.UID] = UserMini{}
						categoryMap[item.CID] = CategoryMini{}
					}
				}
			}

			userKeys := make([][]byte, 0, len(userMap))
			for k := range userMap {
				userKeys = append(userKeys, youdb.I2b(k))
			}
			rs = db.Hmget("user", userKeys)
			if rs.State == "ok" {
				for i := 0; i < (len(rs.Data) - 1); i += 2 {
					item := UserMini{}
					json.Unmarshal(rs.Data[i+1], &item)
					userMap[item.ID] = item
				}
			}

			categoryKeys := make([][]byte, 0, len(categoryMap))
			for k := range categoryMap {
				categoryKeys = append(categoryKeys, youdb.I2b(k))
			}
			rs = db.Hmget("category", categoryKeys)
			if rs.State == "ok" {
				for i := 0; i < (len(rs.Data) - 1); i += 2 {
					item := CategoryMini{}
					json.Unmarshal(rs.Data[i+1], &item)
					categoryMap[item.ID] = item
				}
			}

			for _, article := range aitems {
				user := userMap[article.UID]
				category := categoryMap[article.CID]
				item := ArticleFeedListItem{
					ID:          article.ID,
					UID:         article.UID,
					Name:        user.Name,
					Cname:       html.EscapeString(category.Name),
					Title:       html.EscapeString(article.Title),
					AddTimeFmt:  util.TimeFmt(article.AddTime, time.RFC3339, tz),
					EditTimeFmt: util.TimeFmt(article.EditTime, time.RFC3339, tz),
				}

				contentRune := []rune(article.Content)
				if len(contentRune) > 150 {
					contentRune := []rune(article.Content)
					item.Des = string(contentRune[:150])
				} else {
					item.Des = article.Content
				}
				item.Des = html.EscapeString(item.Des)

				items = append(items, item)
			}

			keys = keys[:0]
			if len(items) >= limit {
				break
			}
		}
	}

	return items
}
