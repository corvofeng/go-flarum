package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"sort"
	"strconv"
	"strings"
	"time"

	"goyoubbs/util"
	"github.com/ego008/youdb"
)

type Article struct {
	Id           uint64 `json:"id"`
	Uid          uint64 `json:"uid"`
	Cid          uint64 `json:"cid"`
	RUid         uint64 `json:"ruid"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	ClientIp     string `json:"clientip"`
	Tags         string `json:"tags"`
	AddTime      uint64 `json:"addtime"`
	EditTime     uint64 `json:"edittime"`
	Comments     uint64 `json:"comments"`
	CloseComment bool   `json:"closecomment"`
	Hidden       bool   `json:"hidden"` // Depreacte, do not use it.

	// 帖子被管理员修改后, 已经保存的旧的帖子ID
	FatherTopicID uint64 `json:"fathertopicid"`

	// 记录当前帖子是否可以被用户看到, 与上面的hidden类似
	Active uint64 `json:"active"`
}

type ArticleMini struct {
	Id       uint64 `json:"id"`
	Uid      uint64 `json:"uid"`
	Cid      uint64 `json:"cid"`
	Ruid     uint64 `json:"ruid"`
	Title    string `json:"title"`
	EditTime uint64 `json:"edittime"`
	Comments uint64 `json:"comments"`
	Hidden   bool   `json:"hidden"`
}

type ArticleListItem struct {
	Id          uint64 `json:"id"`
	Uid         uint64 `json:"uid"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	Cid         uint64 `json:"cid"`
	Cname       string `json:"cname"`
	Ruid        uint64 `json:"ruid"`
	Rname       string `json:"rname"`
	Title       string `json:"title"`
	EditTime    uint64 `json:"edittime"`
	EditTimeFmt string `json:"edittimefmt"`
	ClickCnt    uint64 `json:"clickcnt"`
	Comments    uint64 `json:"comments"`
}

type ArticlePageInfo struct {
	Items      []ArticleListItem `json:"items"`
	HasPrev    bool              `json:"hasprev"`
	HasNext    bool              `json:"hasnext"`
	FirstKey   uint64            `json:"firstkey"`
	FirstScore uint64            `json:"firstscore"`
	LastKey    uint64            `json:"lastkey"`
	PageNum    uint64            `json:"pagenum"`
	LastScore  uint64            `json:"lastscore"`
}

type ArticleLi struct {
	Id    uint64 `json:"id"`
	Title string `json:"title"`
	Tags  string `json:"tags"`
}

type ArticleRelative struct {
	Articles []ArticleLi
	Tags     []string
}

type ArticleFeedListItem struct {
	Id          uint64
	Uid         uint64
	Name        string
	Cname       string
	Title       string
	AddTimeFmt  string
	EditTimeFmt string
	Des         string
}

// 文章添加、编辑后传给后台任务的信息
type ArticleTag struct {
	Id      uint64
	OldTags string
	NewTags string
}

// SQLArticleGetByID 通过 article id获取内容
func SQLArticleGetByID(db *sql.DB, aid string) (Article, error) {
	obj := Article{}
	rows, err := db.Query(
		"SELECT id, node_id, user_id, title, content, created_at, updated_at, client_ip FROM topic WHERE id = ? and active !=0",
		aid,
	)
	defer func() {
		if rows != nil {
			rows.Close() //可以关闭掉未scan连接一直占用
		}
	}()

	if err != nil {
		fmt.Printf("Query failed,err:%v", err)
	}

	for rows.Next() {
		err = rows.Scan(
			&obj.Id,
			&obj.Cid,
			&obj.Uid,
			&obj.Title,
			&obj.Content,
			&obj.AddTime,
			&obj.EditTime,
			&obj.ClientIp,
		)

		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			return obj, errors.New("No result")
		}
	}

	return obj, nil
}

// SQLCreateTopic 创建主题
func (article *Article) SQLCreateTopic(db *sql.DB) bool {
	row, err := db.Exec(
		("INSERT INTO `topic` " +
			" (`node_id`, `user_id`, `title`, `content`, created_at, updated_at, client_ip, father_topic_id, active)" +
			" VALUES " +
			" (?, ?, ?, ?, ?, ?, ?, ?, ?)"),
		article.Cid,
		article.Uid,
		article.Title,
		article.Content,
		article.AddTime,
		article.EditTime,
		article.ClientIp,
		article.FatherTopicID,
		article.Active,
	)
	if util.CheckError(err, "创建主题") {
		return false
	}
	aid, err := row.LastInsertId()
	article.Id = uint64(aid)

	return true
}

// SQLArticleUpdate 更新当前帖子
func (article *Article) SQLArticleUpdate(db *sql.DB) bool {
	// 更新记录必须要被保存, 配合数据库中的father_topic_id来实现
	// 每次更新主题, 会将前帖子复制为一个新帖子(active=0不被看见),
	// 当前帖子的id没有变化, 但是father_topic_id变为这个新的帖子.
	// 通过father_topic_id组成了链表的关系

	// 以当前帖子为模板创建一个新的帖子
	// 对象中只有简单的数据结构, 浅拷贝即可, 需要将其设为不可见
	oldArticle, err := SQLArticleGetByID(db, strconv.FormatUint(article.Id, 10))
	oldArticle.Active = 0
	if util.CheckError(err, "修改时拷贝") {
		return false
	}
	oldArticle.SQLCreateTopic(db)

	//" (`node_id`, `user_id`, `title`, `content`, created_at, updated_at, client_ip, father_topic_id, active)" +
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
		article.Cid,
		article.Uid,
		article.EditTime,
		article.ClientIp,
		oldArticle.Id,
		article.Id,
	)
	if util.CheckError(err, "更新帖子") {
		return false
	}

	return true
}

// SQLArticleGetByList 通过id列表获取对应的帖子
func SQLArticleGetByList(db *sql.DB, cacheDB *youdb.DB, articleList []int) ArticlePageInfo {
	var items []ArticleListItem
	var hasPrev, hasNext bool
	var firstKey, firstScore, lastKey, lastScore uint64
	var rows *sql.Rows
	var err error
	logger := util.GetLogger()
	articleListStr := ""

	for _, v := range articleList {
		if len(articleListStr) > 0 {
			articleListStr += ", "
		}
		articleListStr += strconv.Itoa(v)
	}
	sql := "select id, title, user_id from topic where id in (" + articleListStr + ")"

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
	for rows.Next() {
		item := ArticleListItem{}
		err = rows.Scan(&item.Id, &item.Title, &item.Uid) //不scan会导致连接不释放
		item.Avatar = GetAvatarByID(db, cacheDB, item.Uid)

		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			continue
		}
		rep := cacheDB.Hget("article_views", youdb.I2b(item.Id))
		item.ClickCnt = rep.Uint64()
		items = append(items, item)
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

func ArticleGetById(db *youdb.DB, aid string) (Article, error) {
	obj := Article{}
	rs := db.Hget("article", youdb.DS2b(aid))
	if rs.State == "ok" {
		json.Unmarshal(rs.Data[0], &obj)
		return obj, nil
	}
	return obj, errors.New(rs.State)
}

// SQLCidArticleList 返回某个节点的主题
// nodeID 为0 表示全部主题
func SQLCidArticleList(db *sql.DB, cntDB *youdb.DB, nodeID, start uint64, btnAct string, limit, tz int) ArticlePageInfo {
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
		err = rows.Scan(&item.Id, &item.Title, &item.Uid, &item.Cid, &item.EditTime) //不scan会导致连接不释放
		item.Avatar = GetAvatarByID(db, cntDB, item.Uid)
		item.EditTimeFmt = util.TimeFmt(item.EditTime, "2006-01-02 15:04", tz)
		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			continue
		}
		rep := cntDB.Hget("article_views", youdb.I2b(item.Id))
		item.ClickCnt = rep.Uint64()
		items = append(items, item)
	}
	if len(items) > 0 {
		firstKey = items[0].Id
		lastKey = items[len(items)-1].Id
		hasNext = true
		hasPrev = true

		// 前一页, 后一页的判断其实比较复杂的, 这里只针对最容易的情况进行了判断,
		// 因为帖子较多时, 这算是一种近似

		// 如果最开始的帖子ID为1, 那肯定是没有了前一页了
		if items[0].Id == 1 || start < uint64(limit) {
			hasPrev = false
		}

		// 查询出的数量比要求的数量要少, 说明没有下一页
		if len(items) < limit {
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

// SQLArticleList 返回所有节点的主题
func SQLArticleList(db *sql.DB, cntDB *youdb.DB, start uint64, btnAct string, limit, tz int) ArticlePageInfo {
	return SQLCidArticleList(
		db, cntDB, 0, start, btnAct, limit, tz,
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
					userMap[item.Uid] = UserMini{}
					if item.Ruid > 0 {
						userMap[item.Ruid] = UserMini{}
					}
					categoryMap[item.Cid] = CategoryMini{}
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
				userMap[item.Id] = item
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
				categoryMap[item.Id] = item
			}
		}

		for _, article := range aitems {
			user := userMap[article.Uid]
			category := categoryMap[article.Cid]
			item := ArticleListItem{
				Id:          article.Id,
				Uid:         article.Uid,
				Name:        user.Name,
				Avatar:      user.Avatar,
				Cid:         article.Cid,
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
				firstKey = item.Id
				firstScore = item.EditTime
			}
			lastKey = item.Id
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

func ArticleGetRelative(db *youdb.DB, aid uint64, tags string) ArticleRelative {
	if len(tags) == 0 {
		return ArticleRelative{}
	}
	getMax := 10
	scanMax := 100

	var aitems []ArticleLi
	var titems []string

	tagsLow := strings.ToLower(tags)

	ctagMap := map[string]struct{}{}

	aidCount := map[uint64]int{}

	for _, tag := range strings.Split(tagsLow, ",") {
		ctagMap[tag] = struct{}{}
		rs := db.Hrscan("tag:"+tag, []byte(""), scanMax)
		if rs.State == "ok" {
			for i := 0; i < len(rs.Data)-1; i += 2 {
				aid2 := youdb.B2i(rs.Data[i])
				if aid2 != aid {
					if _, ok := aidCount[aid2]; ok {
						aidCount[aid2] += 1
					} else {
						aidCount[aid2] = 1
					}
				}
			}
		}
	}

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

		rs := db.Hmget("article", akeys)
		if rs.State == "ok" {
			tmpMap := map[string]struct{}{}
			for i := 0; i < len(rs.Data)-1; i += 2 {
				item := ArticleLi{}
				json.Unmarshal(rs.Data[i+1], &item)
				aitems = append(aitems, item)
				for _, tag := range strings.Split(strings.ToLower(item.Tags), ",") {
					if _, ok := ctagMap[tag]; !ok {
						tmpMap[tag] = struct{}{}
					}
				}
			}

			for k := range tmpMap {
				titems = append(titems, k)
			}
		}
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
				userMap[item.Uid] = UserMini{}
				if item.Ruid > 0 {
					userMap[item.Ruid] = UserMini{}
				}
				categoryMap[item.Cid] = CategoryMini{}
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
				userMap[item.Id] = item
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
				categoryMap[item.Id] = item
			}
		}

		for _, article := range aitems {
			user := userMap[article.Uid]
			category := categoryMap[article.Cid]
			item := ArticleListItem{
				Id:          article.Id,
				Uid:         article.Uid,
				Name:        user.Name,
				Avatar:      user.Avatar,
				Cid:         article.Cid,
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
				firstKey = item.Id
			}
			lastKey = item.Id
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
				userMap[item.Uid] = UserMini{}
				if item.Ruid > 0 {
					userMap[item.Ruid] = UserMini{}
				}
				categoryMap[item.Cid] = CategoryMini{}
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
				userMap[item.Id] = item
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
				categoryMap[item.Id] = item
			}
		}

		for _, article := range aitems {
			user := userMap[article.Uid]
			category := categoryMap[article.Cid]
			item := ArticleListItem{
				Id:          article.Id,
				Uid:         article.Uid,
				Name:        user.Name,
				Avatar:      user.Avatar,
				Cid:         article.Cid,
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
						userMap[aitem.Uid] = UserMini{}
						if aitem.RUid > 0 {
							userMap[aitem.RUid] = UserMini{}
						}
						categoryMap[aitem.Cid] = CategoryMini{}
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
				userMap[item.Id] = item
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
				categoryMap[item.Id] = item
			}
		}

		for _, article := range aitems {
			user := userMap[article.Uid]
			category := categoryMap[article.Cid]
			item := ArticleListItem{
				Id:          article.Id,
				Uid:         article.Uid,
				Name:        user.Name,
				Avatar:      user.Avatar,
				Cid:         article.Cid,
				Cname:       category.Name,
				Ruid:        article.RUid,
				Title:       article.Title,
				EditTime:    article.EditTime,
				EditTimeFmt: util.TimeFmt(article.EditTime, "2006-01-02 15:04", tz),
				Comments:    article.Comments,
			}
			if article.RUid > 0 {
				item.Rname = userMap[article.RUid].Name
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
						userMap[item.Uid] = UserMini{}
						categoryMap[item.Cid] = CategoryMini{}
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
					userMap[item.Id] = item
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
					categoryMap[item.Id] = item
				}
			}

			for _, article := range aitems {
				user := userMap[article.Uid]
				category := categoryMap[article.Cid]
				item := ArticleFeedListItem{
					Id:          article.Id,
					Uid:         article.Uid,
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
