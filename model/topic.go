package model

import (
	"errors"
	"fmt"
	"strconv"

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

	CommentCount uint64

	ClientIP string `json:"clientip"`
	IsSticky bool

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

	/**
	 * When in search page, every article item have the highlight content,
	 * we need tell users that we do not return a random list.
	 */
	HighlightContent template.HTML `json:"highlight_content"`
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
	articleBaseList, err := sqlGetTopicByList(gormDB, []uint64{aid})
	var obj Topic
	if len(articleBaseList) == 0 {
		return obj, errors.New("no result")
	}
	obj = articleBaseList[0]

	return obj, err
}

// GetWeight 获取当前帖子的权重
/**
 * (Log10(QView) * 2 + 4 * comments)/ QAge
 *
 * db (*sql.DB): TODO
 * redisDB (redis.Client): TODO
 */
func (article *Topic) GetWeight(redisDB *redis.Client) float64 {
	// var editTime time.Time
	// var now = time.Now()
	// if article.EditTime == 0 {
	// 	editTime = time.Unix(now.Unix()-2*24*3600, 0)
	// } else {
	// 	editTime = time.Unix(int64(article.EditTime), 0)
	// 	sTime := time.Unix(1577836800, 0) // 2020-01-01 08:00:00
	// 	if editTime.Before(sTime) {
	// 		editTime = time.Unix(now.Unix()-2*24*3600, 0)
	// 	}
	// }
	// if article.ClickCnt == 0 { // 避免出现0的情况
	// 	article.ClickCnt = 1
	// }
	// qAge := now.Sub(editTime).Hours()
	// weight := (math.Log10(float64(article.ClickCnt))*2 + 4*float64(article.GetCommentsSize(db))) / (qAge * 1.0)
	return 0
}

// CreateFlarumTopic 创建flarum的帖子
// 帖子中, category和tag是不同的数据
// category是帖子比较大的分类, 每个帖子只能有一个
// tag只是这个帖子具有的某种特征, 每个帖子可以有多个tag
func (topic *Topic) CreateFlarumTopic(gormDB *gorm.DB) (bool, error) {
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

	tx.Save(&topic)

	result = tx.Commit()
	if result.Error != nil {
		logger.Error("Create reply with error", result.Error)
		return false, result.Error
	}

	return true, nil
}

// sqlGetTopicByList 获取帖子信息, NOTE: 请尽量调用该函数, 而不是自己去写sql语句
func sqlGetTopicByList(gormDB *gorm.DB, articleList []uint64) (topics []Topic, err error) {
	err = gormDB.Preload("Tags").Find(&topics, articleList).Error
	return
}

// tagID 为0 表示全部主题
func SQLGetTopicByTag(gormDB *gorm.DB, redisDB *redis.Client, tagID, start uint64, limit uint64) (topics []Topic, err error) {
	logger := util.GetLogger()
	var tag Tag

	ormFilter := gormDB.Preload("Tags").Limit(int(limit)).Offset(int(start))
	if tagID != 0 {
		tag, err = SQLGetTagByID(gormDB, tagID)
		if err != nil {
			logger.Error("Can't get tag by id", tagID)
		}
		err = ormFilter.Model(&tag).Association("Topics").Find(&topics)
	} else {
		err = ormFilter.Find(&topics).Error
	}

	if err != nil {
		logger.Errorf("Can't get all topics for %d `%s`", tagID, err)
	}

	return topics, err
}

func SQLGetTopicByUser(gormDB *gorm.DB, userID, start uint64, limit uint64) (topics []Topic, err error) {
	var user User
	ormFilter := gormDB.Preload("Tags").Limit(int(limit)).Offset(int(start))
	if userID != 0 {
		user, err = SQLUserGetByID(gormDB, userID)
		if err != nil {
			return
		}
		err = ormFilter.Where("user_id = ?", user.ID).Find(&topics).Error
	} else {
		err = ormFilter.Find(&topics).Error
	}

	return
}

// only for rank
func sqlGetAllArticleWithCID(cid uint64, active bool) ([]ArticleMini, error) {
	var articles []ArticleMini
	// var rows *sql.Rows
	// var err error
	// logger := util.GetLogger()

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
	// for rows.Next() {
	// 	obj := ArticleMini{}
	// 	err = rows.Scan(&obj.ID)
	// 	if err != nil {
	// 		logger.Error("Scan failed", err)
	// 		return articles, err
	// 	}

	// 	articles = append(articles, obj)
	// }

	return articles, nil
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

func (topic *Topic) toKeyForComments() string {
	return fmt.Sprintf("comments-article-%d", topic.ID)
}

// CacheCommentList 缓存当前话题对应的评论ID, 该函数可以用于进行增加或是减少
// 注意这里是有顺序的, 顺序为发帖时间
func (topic *Topic) CacheCommentList(redisDB *redis.Client, comments []Comment, done chan bool) error {
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
