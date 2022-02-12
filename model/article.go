package model

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"html/template"
	"zoe/model/flarum"
	"zoe/util"

	"github.com/go-redis/redis/v7"
	"gorm.io/gorm"
)

// ArticleBase 基础的文档类, 在数据库表中的字段
type ArticleBase struct {
	gorm.Model
	ID  uint64 `json:"id"`
	UID uint64 `gorm:"column:user_id;index"`
	CID uint64 `gorm:"column:category_id;index"`

	Title   string `json:"title"`
	Content string `json:"content"`

	FirstPostID uint64
	LastPostID  uint64

	ClickCnt   uint64 `json:"clickcnt"` // 不保证精确性
	ReplyCount uint64

	AddTime  uint64 `json:"addtime"`
	EditTime uint64 `json:"edittime"`

	ClientIP string `json:"clientip"`
}

// 设置comment内容
func (ArticleBase) TableName() string {
	return "topic"
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
	AddTimeFmt  string `json:"addtimefmt"`
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
func SQLArticleGetByID(gormDB *gorm.DB, db *sql.DB, redisDB *redis.Client, aid uint64) (Article, error) {
	articleBaseList := sqlGetArticleBaseByList(gormDB, db, redisDB, []uint64{aid})
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
 * redisDB (redis.Client): TODO
 */
func (article *Article) GetCommentsSize(db *sql.DB) uint64 {
	return article.ReplyCount
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
			" (?, ?, ?, ?, ?, ?, ?, ?, ?)"),
		article.CID,
		article.UID,
		article.Title,
		article.Content,
		article.AddTime,
		article.EditTime,
		article.ClientIP,
		article.ReplyCount,
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
func (article *Article) SQLArticleUpdate(gormDB *gorm.DB, db *sql.DB, redisDB *redis.Client) bool {
	// 更新记录必须要被保存, 配合数据库中的father_topic_id来实现
	// 每次更新主题, 会将前帖子复制为一个新帖子(active=0不被看见),
	// 当前帖子的id没有变化, 但是father_topic_id变为这个新的帖子.
	// 通过father_topic_id组成了链表的关系

	// 以当前帖子为模板创建一个新的帖子
	// 对象中只有简单的数据结构, 浅拷贝即可, 需要将其设为不可见
	oldArticle, err := SQLArticleGetByID(gormDB, db, redisDB, article.ID)
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

// ToArticleListItem 转换为可以做列表的内容
func (ab *ArticleBase) ToArticleListItem(gormDB *gorm.DB, sqlDB *sql.DB, redisDB *redis.Client, tz int) ArticleListItem {
	item := ArticleListItem{
		ArticleBase: *ab,
	}
	item.EditTimeFmt = util.TimeFmt(item.EditTime, util.TIME_FMT, tz)
	item.AddTimeFmt = util.TimeFmt(item.AddTime, util.TIME_FMT, tz)
	item.Cname = GetCategoryNameByCID(sqlDB, redisDB, item.CID)
	if item.LastPostID != 0 {
		lastComment, err := SQLCommentByID(gormDB, sqlDB, redisDB, item.LastPostID, tz)
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
func SQLArticleGetByList(gormDB *gorm.DB, db *sql.DB, redisDB *redis.Client, articleList []uint64, tz int) ArticlePageInfo {
	var items []ArticleListItem
	var hasPrev, hasNext bool
	var firstKey, firstScore, lastKey, lastScore uint64
	articleBaseList := sqlGetArticleBaseByList(gormDB, db, redisDB, articleList)
	m := make(map[uint64]ArticleListItem)

	for _, articleBase := range articleBaseList {
		m[articleBase.ID] = articleBase.ToArticleListItem(gormDB, db, redisDB, tz)
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
func sqlGetArticleBaseByList(gormDB *gorm.DB, db *sql.DB, redisDB *redis.Client, articleList []uint64) (items []ArticleBase) {
	logger := util.GetLogger()
	result := gormDB.Find(&items, articleList)
	if result.Error != nil {
		logger.Errorf("Can't get article list by ", articleList)

	}
	for _, article := range items {
		article.ClickCnt = GetArticleCntFromRedisDB(db, redisDB, article.ID)
	}
	return
}

// SQLArticleGetByCID 根据页码获取某个分类的列表
func SQLArticleGetByCID(gormDB *gorm.DB, db *sql.DB, redisDB *redis.Client, nodeID, page, limit uint64, tz int) ArticlePageInfo {
	var pageInfo ArticlePageInfo
	articleList := GetTopicListByPageNum(nodeID, page, limit)
	logger := util.GetLogger()
	logger.Debug("Get article list", page, limit, articleList)
	if len(articleList) == 0 {
		// TODO: remove it
		articleIteratorStart := GetCIDArticleMax(nodeID)
		pageInfo = SQLCIDArticleList(gormDB, db, redisDB, nodeID, articleIteratorStart, "next", limit, tz)
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
		pageInfo = SQLArticleGetByList(gormDB, db, redisDB, articleList, tz)
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
func SQLArticleGetByUID(gormDB *gorm.DB, db *sql.DB, redisDB *redis.Client, uid, page, limit uint64, tz int) ArticlePageInfo {
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
			logger.Errorf("Scan failed,err:%v", err)
			continue
		}
		articleList = append(articleList, aid)
	}

	logger.Debug("Get article list", page, limit, articleList)
	pageInfo = SQLArticleGetByList(gormDB, db, redisDB, articleList, tz)

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
func SQLCIDArticleList(gormDB *gorm.DB, db *sql.DB, redisDB *redis.Client, nodeID, start uint64, btnAct string, limit uint64, tz int) ArticlePageInfo {
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

	pageInfo := SQLArticleGetByList(gormDB, db, redisDB, articleList, tz)

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
		logger := util.GetLogger()
		logger.Errorf("Get %d with error :%v", aid, err)
		data = 0
	}
	return data
}

// SQLArticleList 返回所有节点的主题
func SQLArticleList(gormDB *gorm.DB, db *sql.DB, redisDB *redis.Client, start uint64, btnAct string, limit uint64, tz int) ArticlePageInfo {
	return SQLCIDArticleList(
		gormDB, db, redisDB, 0, start, btnAct, limit, tz,
	)
}

func (article *Article) toKeyForComments() string {
	return fmt.Sprintf("comments-article-%d", article.ID)
}

// CacheCommentList 缓存当前话题对应的评论ID, 该函数可以用于进行增加或是减少
// 注意这里是有顺序的, 顺序为发帖时间
func (article *Article) CacheCommentList(redisDB *redis.Client, comments []CommentListItem, done chan bool) error {
	logger := util.GetLogger()
	logger.Debugf("Cache comment list for: %d, and %d comments", article.ID, len(comments))
	for _, c := range comments {
		_, err := rankRedisDB.ZAddNX(article.toKeyForComments(), &redis.Z{
			Score:  float64(c.CreatedAt.Unix()),
			Member: c.ID},
		).Result()
		util.CheckError(err, "更新redis中的话题的评论信息")
	}
	done <- true
	return nil
}

// GetCommentIDList 获取帖子已经排序好的评论列表
func (article *Article) GetCommentIDList(redisDB *redis.Client) (comments []uint64) {
	rdsData, _ := rankRedisDB.ZRange(article.toKeyForComments(), 0, -1).Result()
	for _, _cid := range rdsData {
		cid, _ := strconv.ParseUint(_cid, 10, 64)
		comments = append(comments, cid)
	}
	return
}

func (article *Article) CleanCache() {
	logger := util.GetLogger()
	rankRedisDB.Del(article.toKeyForComments())
	logger.Info("Delete comment list cache for: ", article.ID)
}
