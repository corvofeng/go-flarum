package model

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"html/template"

	"github.com/corvofeng/go-flarum/model/flarum"
	"github.com/corvofeng/go-flarum/util"

	"github.com/go-redis/redis/v7"
	"gorm.io/gorm"
)

// Topic 基础的文档类, 在数据库表中的字段
type Topic struct {
	gorm.Model
	ID     uint64 `gorm:"primaryKey"`
	UserID uint64 `gorm:"column:user_id;index"`

	Title   string `json:"title"`
	Content string `json:"content"`

	FirstPostID uint64

	LastPostID     uint64
	LastPostUserID uint64
	LastPostAt     time.Time

	ClickCnt     uint64 `json:"clickcnt"` // 不保证精确性
	CommentCount uint64

	AddTime  uint64 `json:"addtime"`
	EditTime uint64 `json:"edittime"`

	ClientIP string `json:"clientip"`

	Tags []Tag `gorm:"many2many:topic_tags;"`
}

// TopicTags 帖子的标签
// 使用gorm 的many2many, 不需要单独初始化了
type TopicTag struct {
	gorm.Model
	TopicID uint64 `gorm:"primaryKey"`
	TagID   uint64 `gorm:"primaryKey"`
}

// ArticleMini 缩略版的Article信息
type ArticleMini struct {
	Topic
	Ruid   uint64 `json:"ruid"`
	Hidden bool   `json:"hidden"`
}

// ArticleListItem data strucy only used in page.
type ArticleListItem struct {
	Topic
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
func SQLArticleGetByID(gormDB *gorm.DB, redisDB *redis.Client, aid uint64) (Topic, error) {
	articleBaseList := sqlGetTopicByList(gormDB, redisDB, []uint64{aid})
	var obj Topic
	if len(articleBaseList) == 0 {
		return obj, errors.New("No result")
	}
	obj = articleBaseList[0]

	return obj, nil
}

// GetCommentsSize 获取评论
/*
 * db (*sql.DB): TODO
 * redisDB (redis.Client): TODO
 */
// func (topic *Topic) GetCommentsSize(db *sql.DB) uint64 {
// 	return topic.CommentCount
// }

// GetWeight 获取当前帖子的权重
/**
 * (Log10(QView) * 2 + 4 * comments)/ QAge
 *
 * db (*sql.DB): TODO
 * redisDB (redis.Client): TODO
 */
func (article *Topic) GetWeight(redisDB *redis.Client) float64 {
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
	// qAge := now.Sub(editTime).Hours()
	// weight := (math.Log10(float64(article.ClickCnt))*2 + 4*float64(article.GetCommentsSize(db))) / (qAge * 1.0)
	return 0
}

// CreateFlarumTopic 创建flarum的帖子
// 帖子中, category和tag是不同的数据
// category是帖子比较大的分类, 每个帖子只能有一个
// tag只是这个帖子具有的某种特征, 每个帖子可以有多个tag
func (topic *Topic) CreateFlarumTopic(gormDB *gorm.DB, tags flarum.RelationArray) (bool, error) {
	logger := util.GetLogger()
	tx := gormDB.Begin()
	defer clearGormTransaction(tx)

	result := tx.Create(&topic)

	if result.Error != nil {
		logger.Error("Can't creat topic with error", result.Error)
		return false, result.Error
	}

	comment := Comment{
		Reply: Reply{
			AID:      topic.ID,
			UID:      topic.UserID,
			Content:  topic.Content,
			Number:   1,
			ClientIP: topic.ClientIP,
		},
	}
	result = tx.Create(&comment.Reply)
	if result.Error != nil {
		logger.Error("Can't create first commet for topic with error", result.Error)
		return false, result.Error
	}
	topic.LastPostID = comment.ID
	topic.FirstPostID = comment.ID

	for _, tid := range tags.Data {
		fmt.Println(tid)
	}

	tx.Save(&topic)

	result = tx.Commit()
	if result.Error != nil {
		logger.Error("Create reply with error", result.Error)
		return false, result.Error
	}

	return true, nil
}

// updateFlarumTag 更新评论信息
func (article *Topic) updateFlarumTag(tx *sql.Tx, tags flarum.RelationArray) (bool, error) {
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

// ToArticleListItem 转换为可以做列表的内容
// func (ab *Topic) ToArticleListItem(gormDB *gorm.DB, redisDB *redis.Client, tz int) ArticleListItem {
// 	item := ArticleListItem{
// 		Topic: *ab,
// 	}
// 	item.EditTimeFmt = item.UpdatedAt.UTC().String()
// 	item.AddTimeFmt = item.CreatedAt.UTC().String()
// 	if item.LastPostID != 0 {
// 		lastComment, err := SQLCommentByID(gormDB, redisDB, item.LastPostID, tz)
// 		if err != nil {
// 			util.GetLogger().Errorf("Can't get last comment(%d)for article(%d)", item.LastPostID, item.ID)
// 		} else {
// 			lc := lastComment.toCommentListItem(   redisDB, tz)
// 			item.LastComment = &lc
// 		}
// 	}

// 	return item
// }

// SQLArticleGetByList 通过id列表获取对应的帖子
func SQLArticleGetByList(gormDB *gorm.DB, redisDB *redis.Client, articleList []uint64, tz int) ArticlePageInfo {
	var items []ArticleListItem
	var hasPrev, hasNext bool
	var firstKey, firstScore, lastKey, lastScore uint64
	m := make(map[uint64]ArticleListItem)

	// articleBaseList := sqlGetTopicByList(gormDB, db, redisDB, articleList)
	// for _, articleBase := range articleBaseList {
	// 	m[articleBase.ID] = articleBase.ToArticleListItem(gormDB, db, redisDB, tz)
	// }

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

// sqlGetTopicByList 获取帖子信息, NOTE: 请尽量调用该函数, 而不是自己去写sql语句
func sqlGetTopicByList(gormDB *gorm.DB, redisDB *redis.Client, articleList []uint64) (items []Topic) {
	logger := util.GetLogger()
	result := gormDB.Find(&items, articleList)
	if result.Error != nil {
		logger.Errorf("Can't get article list by ", articleList)
	}
	for _, article := range items {
		article.ClickCnt = GetArticleCntFromRedisDB(redisDB, article.ID)
	}
	return
}

// SQLArticleGetByCID 根据页码获取某个分类的列表
func SQLTopicGetByTag(gormDB *gorm.DB, redisDB *redis.Client, tagID, start, limit uint64, tz int) []Topic {
	return SQLCIDArticleList(gormDB, redisDB, tagID, start, limit, tz)
}

// SQLTopicGetByUID 根据创建用户获取帖子列表
func SQLTopicGetByUID(gormDB *gorm.DB, redisDB *redis.Client, uid, page, limit uint64, tz int) ArticlePageInfo {
	var rows *sql.Rows
	var err error
	var pageInfo ArticlePageInfo
	var articleList []uint64
	logger := util.GetLogger()
	defer rowsClose(rows)

	// rows, err = db.Query(
	// 	"SELECT id FROM topic WHERE user_id = ? and active = 1 ORDER BY created_at DESC limit ? offset ?",
	// 	uid, limit, (page-1)*limit,
	// )
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
	pageInfo = SQLArticleGetByList(gormDB, redisDB, articleList, tz)

	pageInfo.PageNum = page
	pageInfo.PageNext = page + 1
	pageInfo.PagePrev = page - 1
	if len(articleList) == int(limit) {
		pageInfo.HasNext = true
	}
	return pageInfo
}

// SQLCIDArticleList 返回某个节点的主题
// tagID 为0 表示全部主题
func SQLCIDArticleList(gormDB *gorm.DB, redisDB *redis.Client, tagID, start uint64, limit uint64, tz int) []Topic {
	logger := util.GetLogger()
	var topics []Topic
	// var err error
	if tagID != 0 {
		tag, err := SQLGetTagByID(gormDB, tagID)
		if err != nil {
			logger.Error("Can't get tag by id", tagID)
		}
		err = gormDB.Limit(int(limit)).Offset(int(start)).Model(&tag).Association("Topics").Find(&topics)
		if err != nil {
			logger.Error("Can't get topics by tag id", tagID)
		}
	} else {
		rlt := gormDB.Limit(int(limit)).Offset(int(start)).Find(&topics)
		if rlt.Error != nil {
			logger.Info("Can't get all topics", rlt.Error)
		}
	}
	return topics
}

// only for rank
func sqlGetAllArticleWithCID(cid uint64, active bool) ([]ArticleMini, error) {
	var articles []ArticleMini
	var rows *sql.Rows
	var err error
	logger := util.GetLogger()

	// activeData := 0

	// if active {
	// 	activeData = 1
	// } else {
	// 	activeData = 0
	// }

	// if cid == 0 {
	// 	// cid为0, 查询所有节点
	// 	rows, err = db.Query(
	// 		"SELECT t_list.topic_id FROM (SELECT topic_id FROM `topic_tag`) as t_list LEFT JOIN topic ON t_list.topic_id = topic.id WHERE active = ?",
	// 		activeData)
	// } else {
	// 	rows, err = db.Query(
	// 		"SELECT t_list.topic_id FROM (SELECT topic_id FROM `topic_tag` where tag_id = ?) as t_list LEFT JOIN topic ON t_list.topic_id = topic.id WHERE active = ?",
	// 		cid, activeData)
	// }
	// defer rowsClose(rows)

	// if err != nil {
	// 	logger.Error("Can't get topic_tag info", err)
	// 	return articles, err
	// }
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
/*
func (article *Topic) IncrArticleCntFromRedisDB(redisDB *redis.Client) uint64 {
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
*/

// GetArticleCntFromRedisDB 获取当前帖子的点击次数
/*
 * sqlDB (*sql.DB): TODO
 * redisDB (*redis.Client): TODO
 */
func (article *Topic) GetArticleCntFromRedisDB(redisDB *redis.Client) uint64 {
	article.ClickCnt = GetArticleCntFromRedisDB(redisDB, article.ID)
	return article.ClickCnt
}

// GetArticleCntFromRedisDB 从不同的数据库中获取点击数
func GetArticleCntFromRedisDB(redisDB *redis.Client, aid uint64) uint64 {
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
func SQLArticleList(gormDB *gorm.DB, redisDB *redis.Client, start uint64, btnAct string, limit uint64, tz int) []Topic {
	return SQLCIDArticleList(
		gormDB, redisDB, 0, start, limit, tz,
	)
}

func (topic *Topic) toKeyForComments() string {
	return fmt.Sprintf("comments-article-%d", topic.ID)
}

// CacheCommentList 缓存当前话题对应的评论ID, 该函数可以用于进行增加或是减少
// 注意这里是有顺序的, 顺序为发帖时间
func (topic *Topic) CacheCommentList(redisDB *redis.Client, comments []CommentListItem, done chan bool) error {
	logger := util.GetLogger()
	logger.Debugf("Cache comment list for: %d, and %d comments", topic.ID, len(comments))
	for _, c := range comments {
		_, err := rankRedisDB.ZAddNX(topic.toKeyForComments(), &redis.Z{
			Score:  float64(c.CreatedAt.Unix()),
			Member: c.ID},
		).Result()
		util.CheckError(err, "更新redis中的话题的评论信息")
	}
	done <- true
	return nil
}

// GetCommentIDList 获取帖子已经排序好的评论列表
func (topic *Topic) GetCommentIDList(redisDB *redis.Client) (comments []uint64) {
	rdsData, _ := rankRedisDB.ZRange(topic.toKeyForComments(), 0, -1).Result()
	for _, _cid := range rdsData {
		cid, _ := strconv.ParseUint(_cid, 10, 64)
		comments = append(comments, cid)
	}
	return
}

func (article *Topic) CleanCache() {
	logger := util.GetLogger()
	rankRedisDB.Del(article.toKeyForComments())
	logger.Info("Delete comment list cache for: ", article.ID)
}
