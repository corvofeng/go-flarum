package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"goyoubbs/util"
	"html/template"

	"github.com/go-redis/redis/v7"

	"github.com/ego008/youdb"
)

// ArticleBase 基础的文档类
type ArticleBase struct {
	ID  uint64 `json:"id"`
	UID uint64 `json:"uid"`
	CID uint64 `json:"cid"`
}

// Article store in database
type Article struct {
	ArticleBase
	RUID              uint64 `json:"ruid"`
	Title             string `json:"title"`
	Content           string `json:"content"`
	ClientIP          string `json:"clientip"`
	Tags              string `json:"tags"`
	AddTime           uint64 `json:"addtime"`
	EditTime          uint64 `json:"edittime"`
	Comments          uint64 `json:"comments"`
	ClickCnt          uint64 `json:"clickcnt"`
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
	Ruid     uint64 `json:"ruid"`
	Title    string `json:"title"`
	EditTime uint64 `json:"edittime"`
	Comments uint64 `json:"comments"`
	Hidden   bool   `json:"hidden"`
}

// ArticleListItem data strucy only used in page.
type ArticleListItem struct {
	ArticleBase
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	Cname       string `json:"cname"`
	Ruid        uint64 `json:"ruid"`
	Rname       string `json:"rname"`
	Title       string `json:"title"`
	EditTime    uint64 `json:"edittime"`
	EditTimeFmt string `json:"edittimefmt"`
	ClickCnt    uint64 `json:"clickcnt"`
	Comments    uint64 `json:"comments"`

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
	logger := util.GetLogger()
	obj := Article{}
	rows, err := db.Query(
		"SELECT id, node_id, user_id, title, content, created_at, updated_at, father_topic_id, client_ip FROM topic WHERE id = ? and active !=0",
		aid,
	)
	defer func() {
		if rows != nil {
			rows.Close() //可以关闭掉未scan连接一直占用
		}
	}()
	util.CheckError(err, fmt.Sprintf("Query %d failed", aid))

	if rows.Next() {
		err = rows.Scan(
			&obj.ID,
			&obj.CID,
			&obj.UID,
			&obj.Title,
			&obj.Content,
			&obj.AddTime,
			&obj.EditTime,
			&obj.FatherTopicID,
			&obj.ClientIP,
		)

		if util.CheckError(err, fmt.Sprintf("scan %d", aid)) {
			return obj, errors.New("No result")
		}
		obj.GetArticleCntFromRedisDB(db, redisDB)
	} else {
		logger.Debug("failed query ", aid)
		return obj, errors.New("No result")
	}

	return obj, nil
}

// FlarumArticleGetByID 获取主题
func FlarumArticleGetByID(db *sql.DB, aid uint64) map[string]interface{} {
	dictData := executeQuery(
		db,
		"SELECT id, node_id, user_id, title, content, created_at, updated_at, father_topic_id, client_ip FROM topic WHERE id = ? and active !=0",
		aid,
	)

	if len(dictData) == 0 {
		return nil
	}

	return dictData[0]
}

// LoadDictData 从dict cursor的结果中获取数据
func (article *Article) LoadDictData(map[string]interface{}) {

}

// GetCommentsSize 获取评论
/*
 * db (*sql.DB): TODO
 * cntDB (*youdb.DB): TODO
 * redisDB (redis.Client): TODO
 */
func (article *Article) GetCommentsSize(db *sql.DB) int {
	var count int

	row := db.QueryRow(
		"SELECT count(*) FROM reply where topic_id = ?",
		article.ID,
	)

	err := row.Scan(&count)
	if util.CheckError(err, "查询评论数量") {
		count = 0
	}
	// if count != 0 {
	// 	util.GetLogger().Debugf("The %s has comments %d", article.Title, count)
	// }

	return count
}

// GetWeight 获取当前帖子的权重
/**
 * (Log10(QView) * 2 + 4 * comments)/ QAge
 *
 * db (*sql.DB): TODO
 * redisDB (redis.Client): TODO
 */
func (article *Article) GetWeight(db *sql.DB, cntDB *youdb.DB, redisDB *redis.Client) float64 {
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
	qAge := now.Sub(editTime).Hours()
	weight := (math.Log10(float64(article.ClickCnt))*2 + 4*float64(article.GetCommentsSize(db))) / (qAge * 1.0)

	return weight
}

// SQLCreateTopic 创建主题
func (article *Article) SQLCreateTopic(db *sql.DB) bool {
	row, err := db.Exec(
		("INSERT INTO `topic` " +
			" (`node_id`, `user_id`, `title`, `content`, created_at, updated_at, client_ip, father_topic_id, active)" +
			" VALUES " +
			" (?, ?, ?, ?, ?, ?, ?, ?, ?)"),
		article.CID,
		article.UID,
		article.Title,
		article.Content,
		article.AddTime,
		article.EditTime,
		article.ClientIP,
		article.FatherTopicID,
		article.Active,
	)
	if util.CheckError(err, "创建主题") {
		return false
	}
	aid, err := row.LastInsertId()
	article.ID = uint64(aid)

	return true
}

// SQLArticleUpdate 更新当前帖子
func (article *Article) SQLArticleUpdate(db *sql.DB, cntDB *youdb.DB, redisDB *redis.Client) bool {
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

// SQLArticleGetByList 通过id列表获取对应的帖子
func SQLArticleGetByList(db *sql.DB, redisDB *redis.Client, articleList []uint64, tz int) ArticlePageInfo {
	var items []ArticleListItem
	var hasPrev, hasNext bool
	var firstKey, firstScore, lastKey, lastScore uint64
	var rows *sql.Rows
	var err error
	logger := util.GetLogger()
	articleListStr := ""

	var articlePageInfo ArticlePageInfo

	if len(articleList) == 0 {
		logger.Warning("SQLArticleGetByList: Can't process the article list empty")
		return articlePageInfo
	}

	for _, v := range articleList {
		if len(articleListStr) > 0 {
			articleListStr += ", "
		}
		articleListStr += strconv.FormatInt(int64(v), 10)
	}
	sql := "select id, title, node_id, user_id, updated_at from topic where id in (" + articleListStr + ")"

	rows, err = db.Query(sql)
	defer func() {
		if rows != nil {
			rows.Close() //可以关闭掉未scan连接一直占用
		}
	}()
	if err != nil {
		logger.Errorf("Query failed,err:%v", err)
		return ArticlePageInfo{}
	}
	m := make(map[uint64]ArticleListItem)
	for rows.Next() {
		item := ArticleListItem{}
		err = rows.Scan(&item.ID, &item.Title, &item.CID, &item.UID, &item.EditTime) //不scan会导致连接不释放
		item.Avatar = GetAvatarByID(db, redisDB, item.UID)

		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			continue
		}
		item.ClickCnt = GetArticleCntFromRedisDB(db, redisDB, item.ID)
		item.EditTimeFmt = util.TimeFmt(item.EditTime, "2006-01-02 15:04", tz)
		item.Cname = GetCategoryNameByCID(db, redisDB, item.CID)
		m[item.ID] = item
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

// SQLCIDArticleListByPage 根据页码获取某个分类的列表
func SQLCIDArticleListByPage(db *sql.DB, redisDB *redis.Client, nodeID, page, limit uint64, tz int) ArticlePageInfo {
	articleList := GetTopicListByPageNum(nodeID, page, limit)
	var pageInfo ArticlePageInfo
	if len(articleList) == 0 {
		// TODO: remove it
		articleIteratorStart := GetCIDArticleMax(nodeID)
		fmt.Println("Current iterator is in ", articleIteratorStart)
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
	util.CheckError(err, "更新帖子权重")
}

// SQLCIDArticleList 返回某个节点的主题
// nodeID 为0 表示全部主题
func SQLCIDArticleList(db *sql.DB, redisDB *redis.Client, nodeID, start uint64, btnAct string, limit uint64, tz int) ArticlePageInfo {
	var items []ArticleListItem
	var hasPrev, hasNext bool
	var firstKey, firstScore, lastKey, lastScore uint64
	var rows *sql.Rows
	var err error
	logger := util.GetLogger()
	valueList := "id, title, user_id, node_id, updated_at"
	selectList := " active != 0 "
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

	defer func() {
		if rows != nil {
			rows.Close() //可以关闭掉未scan连接一直占用
		}
	}()
	if err != nil {
		logger.Errorf("Query failed,err:%v", err)
		return ArticlePageInfo{}
	}
	for rows.Next() {
		item := ArticleListItem{}
		err = rows.Scan(&item.ID, &item.Title, &item.UID, &item.CID, &item.EditTime) //不scan会导致连接不释放
		item.Avatar = GetAvatarByID(db, redisDB, item.UID)
		item.EditTimeFmt = util.TimeFmt(item.EditTime, "2006-01-02 15:04", tz)
		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			continue
		}
		item.ClickCnt = GetArticleCntFromRedisDB(db, redisDB, item.ID)
		item.Cname = GetCategoryNameByCID(db, redisDB, item.CID)
		items = append(items, item)
	}
	if len(items) > 0 {
		firstKey = items[0].ID
		lastKey = items[len(items)-1].ID
		hasNext = true
		hasPrev = true

		// 前一页, 后一页的判断其实比较复杂的, 这里只针对最容易的情况进行了判断,
		// 因为帖子较多时, 这算是一种近似

		// 如果最开始的帖子ID为1, 那肯定是没有了前一页了
		if items[0].ID == 1 || start < uint64(limit) {
			hasPrev = false
		}

		// 查询出的数量比要求的数量要少, 说明没有下一页
		if uint64(len(items)) < limit {
			hasNext = false
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

func sqlGetAllArticleWithCID(db *sql.DB, cid uint64, active bool) ([]ArticleMini, error) {
	var articles []ArticleMini
	var rows *sql.Rows
	var err error
	if active {
		rows, err = db.Query("SELECT id FROM topic WHERE node_id = ? AND active = ?", cid, 1)
	} else {
		rows, err = db.Query("SELECT id FROM topic WHERE node_id = ? AND active = ?", cid, 0)
	}
	defer func() {
		if rows != nil {
			rows.Close() //可以关闭掉未scan连接一直占用
		}
	}()
	if err != nil {
		fmt.Printf("Query failed,err:%v", err)
		return articles, err
	}
	for rows.Next() {
		obj := ArticleMini{}
		err = rows.Scan(&obj.ID)
		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
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
			// 从cntDB中查找
			// clickCnt, _ = cntDB.Hincr("article_views", youdb.I2b(aid), 1)
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
	// if data == 0 {
	// rep := cntDB.Hget("article_views", youdb.I2b(aid))
	// if rep != nil {
	// data = rep.Uint64()
	// redisDB.HSet("article_views", fmt.Sprintf("%d", aid), data)
	// }
	// }
	return data
}

// SQLArticleList 返回所有节点的主题
func SQLArticleList(db *sql.DB, redisDB *redis.Client, start uint64, btnAct string, limit uint64, tz int) ArticlePageInfo {
	return SQLCIDArticleList(
		db, redisDB, 0, start, btnAct, limit, tz,
	)
}

func ArticleList(db *youdb.DB, cmd, tb, key, score string, limit, tz int) ArticlePageInfo {
	var items []ArticleListItem
	var keys [][]byte
	var hasPrev, hasNext bool
	var firstKey, firstScore, lastKey, lastScore uint64

	keyStart := youdb.DS2b(key)
	scoreStart := youdb.DS2b(score)
	if cmd == "zrscan" {
		rs := db.Zrscan(tb, keyStart, scoreStart, limit)
		if rs.State == "ok" {
			for i := 0; i < (len(rs.Data) - 1); i += 2 {
				keys = append(keys, rs.Data[i])
			}
		}
	} else if cmd == "zscan" {
		rs := db.Zscan(tb, keyStart, scoreStart, limit)
		if rs.State == "ok" {
			for i := len(rs.Data) - 2; i >= 0; i -= 2 {
				keys = append(keys, rs.Data[i])
			}
		}
	}

	if len(keys) > 0 {
		var aitems []ArticleMini
		userMap := map[uint64]UserMini{}
		categoryMap := map[uint64]CategoryMini{}

		rs := db.Hmget("article", keys)
		if rs.State == "ok" {
			for i := 0; i < (len(rs.Data) - 1); i += 2 {
				item := ArticleMini{}
				json.Unmarshal(rs.Data[i+1], &item)
				if !item.Hidden {
					aitems = append(aitems, item)
					userMap[item.UID] = UserMini{}
					if item.Ruid > 0 {
						userMap[item.Ruid] = UserMini{}
					}
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
			item := ArticleListItem{
				ArticleBase: ArticleBase{
					ID:  article.ID,
					UID: article.UID,
					CID: article.CID,
				},
				Name:        user.Name,
				Avatar:      user.Avatar,
				Cname:       category.Name,
				Ruid:        article.Ruid,
				Title:       article.Title,
				EditTime:    article.EditTime,
				EditTimeFmt: util.TimeFmt(article.EditTime, "2006-01-02 15:04", tz),
				Comments:    article.Comments,
			}
			if article.Ruid > 0 {
				item.Rname = userMap[article.Ruid].Name
			}
			items = append(items, item)
			if firstKey == 0 {
				firstKey = item.ID
				firstScore = item.EditTime
			}
			lastKey = item.ID
			lastScore = item.EditTime
		}

		// not fix hidden article
		rs = db.Zscan(tb, youdb.I2b(firstKey), youdb.I2b(firstScore), 1)
		if rs.State == "ok" {
			hasPrev = true
		}
		rs = db.Zrscan(tb, youdb.I2b(lastKey), youdb.I2b(lastScore), 1)
		if rs.State == "ok" {
			hasNext = true
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

func ArticleGetRelative(aid uint64, tags string) ArticleRelative {
	if len(tags) == 0 {
		return ArticleRelative{}
	}
	getMax := 10
	// scanMax := 100

	var aitems []ArticleLi
	var titems []string

	// tagsLow := strings.ToLower(tags)

	// ctagMap := map[string]struct{}{}

	aidCount := map[uint64]int{}

	// for _, tag := range strings.Split(tagsLow, ",") {
	// 	ctagMap[tag] = struct{}{}
	// 	rs := db.Hrscan("tag:"+tag, []byte(""), scanMax)
	// 	if rs.State == "ok" {
	// 		for i := 0; i < len(rs.Data)-1; i += 2 {
	// 			aid2 := youdb.B2i(rs.Data[i])
	// 			if aid2 != aid {
	// 				if _, ok := aidCount[aid2]; ok {
	// 					aidCount[aid2] += 1
	// 				} else {
	// 					aidCount[aid2] = 1
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	if len(aidCount) > 0 {

		type Kv struct {
			Key   uint64
			Value int
		}

		var ss []Kv
		for k, v := range aidCount {
			ss = append(ss, Kv{k, v})
		}

		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value > ss[j].Value
		})

		var akeys [][]byte
		j := 0
		for _, kv := range ss {
			akeys = append(akeys, youdb.I2b(kv.Key))
			j++
			if j == getMax {
				break
			}
		}

		// rs := db.Hmget("article", akeys)
		// if rs.State == "ok" {
		// 	tmpMap := map[string]struct{}{}
		// 	for i := 0; i < len(rs.Data)-1; i += 2 {
		// 		item := ArticleLi{}
		// 		json.Unmarshal(rs.Data[i+1], &item)
		// 		aitems = append(aitems, item)
		// 		for _, tag := range strings.Split(strings.ToLower(item.Tags), ",") {
		// 			if _, ok := ctagMap[tag]; !ok {
		// 				tmpMap[tag] = struct{}{}
		// 			}
		// 		}
		// 	}

		// 	for k := range tmpMap {
		// 		titems = append(titems, k)
		// 	}
		// }
	}

	return ArticleRelative{
		Articles: aitems,
		Tags:     titems,
	}
}

func UserArticleList(db *youdb.DB, cmd, tb, key string, limit, tz int) ArticlePageInfo {
	var items []ArticleListItem
	var keys [][]byte

	var hasPrev, hasNext bool
	var firstKey, lastKey uint64

	keyStart := youdb.DS2b(key)
	if cmd == "hrscan" {
		rs := db.Hrscan(tb, keyStart, limit)
		if rs.State == "ok" {
			for i := 0; i < (len(rs.Data) - 1); i += 2 {
				keys = append(keys, rs.Data[i])
			}
		}
	} else if cmd == "hscan" {
		rs := db.Hscan(tb, keyStart, limit)
		if rs.State == "ok" {
			for i := len(rs.Data) - 2; i >= 0; i -= 2 {
				keys = append(keys, rs.Data[i])
			}
		}
	}

	if len(keys) > 0 {
		var aitems []ArticleMini
		userMap := map[uint64]UserMini{}
		categoryMap := map[uint64]CategoryMini{}

		rs := db.Hmget("article", keys)
		if rs.State == "ok" {
			for i := 0; i < (len(rs.Data) - 1); i += 2 {
				item := ArticleMini{}
				json.Unmarshal(rs.Data[i+1], &item)
				aitems = append(aitems, item)
				userMap[item.UID] = UserMini{}
				if item.Ruid > 0 {
					userMap[item.Ruid] = UserMini{}
				}
				categoryMap[item.CID] = CategoryMini{}
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
			item := ArticleListItem{
				ArticleBase: ArticleBase{
					ID:  article.ID,
					UID: article.UID,
					CID: article.CID,
				},
				Name:        user.Name,
				Avatar:      user.Avatar,
				Cname:       category.Name,
				Ruid:        article.Ruid,
				Title:       article.Title,
				EditTime:    article.EditTime,
				EditTimeFmt: util.TimeFmt(article.EditTime, "2006-01-02 15:04", tz),
				Comments:    article.Comments,
			}
			if article.Ruid > 0 {
				item.Rname = userMap[article.Ruid].Name
			}
			items = append(items, item)
			if firstKey == 0 {
				firstKey = item.ID
			}
			lastKey = item.ID
		}

		rs = db.Hscan(tb, youdb.I2b(firstKey), 1)
		if rs.State == "ok" {
			hasPrev = true
		}
		rs = db.Hrscan(tb, youdb.I2b(lastKey), 1)
		if rs.State == "ok" {
			hasNext = true
		}
	}

	return ArticlePageInfo{
		Items:    items,
		HasPrev:  hasPrev,
		HasNext:  hasNext,
		FirstKey: firstKey,
		LastKey:  lastKey,
	}
}

func ArticleNotificationList(db *youdb.DB, ids string, tz int) ArticlePageInfo {
	var items []ArticleListItem
	var keys [][]byte

	for _, v := range strings.Split(ids, ",") {
		keys = append(keys, youdb.DS2b(v))
	}

	if len(keys) > 0 {
		var aitems []ArticleMini
		userMap := map[uint64]UserMini{}
		categoryMap := map[uint64]CategoryMini{}

		rs := db.Hmget("article", keys)
		if rs.State == "ok" {
			for i := 0; i < (len(rs.Data) - 1); i += 2 {
				item := ArticleMini{}
				json.Unmarshal(rs.Data[i+1], &item)
				aitems = append(aitems, item)
				userMap[item.UID] = UserMini{}
				if item.Ruid > 0 {
					userMap[item.Ruid] = UserMini{}
				}
				categoryMap[item.CID] = CategoryMini{}
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
			item := ArticleListItem{
				ArticleBase: ArticleBase{
					ID:  article.ID,
					UID: article.UID,
					CID: article.CID,
				},
				Name:        user.Name,
				Avatar:      user.Avatar,
				Cname:       category.Name,
				Ruid:        article.Ruid,
				Title:       article.Title,
				EditTime:    article.EditTime,
				EditTimeFmt: util.TimeFmt(article.EditTime, "2006-01-02 15:04", tz),
				Comments:    article.Comments,
			}
			if article.Ruid > 0 {
				item.Rname = userMap[article.Ruid].Name
			}
			items = append(items, item)
		}
	}

	return ArticlePageInfo{Items: items}
}

func ArticleSearchList(db *youdb.DB, where, kw string, limit, tz int) ArticlePageInfo {
	var items []ArticleListItem

	var aitems []Article
	userMap := map[uint64]UserMini{}
	categoryMap := map[uint64]CategoryMini{}

	startKey := []byte("")
	for {
		rs := db.Hrscan("article", startKey, limit)
		if rs.State == "ok" {
			for i := 0; i < (len(rs.Data) - 1); i += 2 {
				startKey = rs.Data[i]
				aitem := Article{}
				json.Unmarshal(rs.Data[i+1], &aitem)
				if !aitem.Hidden {
					var getIt bool
					if where == "title" {
						if strings.Index(strings.ToLower(aitem.Title), kw) >= 0 {
							getIt = true
						}
					} else {
						if strings.Index(strings.ToLower(aitem.Content), kw) >= 0 {
							getIt = true
						}
					}
					if getIt {
						aitems = append(aitems, aitem)
						userMap[aitem.UID] = UserMini{}
						if aitem.RUID > 0 {
							userMap[aitem.RUID] = UserMini{}
						}
						categoryMap[aitem.CID] = CategoryMini{}
						if len(aitems) == limit {
							break
						}
					}
				}
			}
			if len(aitems) == limit {
				break
			}
		} else {
			break
		}
	}

	if len(aitems) > 0 {
		userKeys := make([][]byte, 0, len(userMap))
		for k := range userMap {
			userKeys = append(userKeys, youdb.I2b(k))
		}
		rs := db.Hmget("user", userKeys)
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
			item := ArticleListItem{
				ArticleBase: ArticleBase{
					ID:  article.ID,
					UID: article.UID,
					CID: article.CID,
				},
				Name:        user.Name,
				Avatar:      user.Avatar,
				Cname:       category.Name,
				Ruid:        article.RUID,
				Title:       article.Title,
				EditTime:    article.EditTime,
				EditTimeFmt: util.TimeFmt(article.EditTime, "2006-01-02 15:04", tz),
				Comments:    article.Comments,
			}
			if article.RUID > 0 {
				item.Rname = userMap[article.RUID].Name
			}
			items = append(items, item)
		}
	}

	return ArticlePageInfo{Items: items}
}

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
