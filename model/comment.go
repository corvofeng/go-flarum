package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"

	"goyoubbs/util"

	"github.com/ego008/youdb"
	"github.com/go-redis/redis/v7"
)

// CommentBase 会在数据库中保存的信息
type CommentBase struct {
	ID       uint64 `json:"id"`
	AID      uint64 `json:"aid"`
	UID      uint64 `json:"uid"`
	Number   uint64 `json:"number"`
	Content  string `json:"content"`
	ClientIP string `json:"clientip"`
	AddTime  uint64 `json:"addtime"`
}

// Comment 评论信息
type Comment struct {
	CommentBase
	UserName   string `json:"username"`
	Avatar     string `json:"avatar"`
	ContentFmt template.HTML
	AddTimeFmt string `json:"addtimefmt"`
}

// CommentListItem 页面中的评论
type CommentListItem struct {
	Comment

	Name string `json:"name"`
}

type CommentPageInfo struct {
	Items    []CommentListItem `json:"items"`
	HasPrev  bool              `json:"hasprev"`
	HasNext  bool              `json:"hasnext"`
	FirstKey uint64            `json:"firstkey"`
	LastKey  uint64            `json:"lastkey"`
}

// sqlSaveComment 在数据库中存储评论
func (comment *Comment) sqlSaveComment(tx *sql.Tx) (bool, error) {
	rows, err := tx.Exec(
		("INSERT INTO `reply`" +
			"(`user_id`, `topic_id`, `number`, `client_ip`, `content`, `created_at`, `updated_at`)" +
			"VALUES " +
			"(?, ?, ?, ?, ?, ?, ?)"),
		comment.UID,
		comment.AID,
		comment.Number,
		comment.ClientIP,
		comment.Content,
		comment.AddTime,
		comment.AddTime,
	)
	if util.CheckError(err, "回复失败") {
		return false, err
	}
	cid, err := rows.LastInsertId()
	comment.ID = uint64(cid)
	return true, nil
}

// SQLSaveComment 在数据库中存储评论
func (comment *Comment) SQLSaveComment(db *sql.DB) {
	rows, err := db.Exec(
		("INSERT INTO `reply`" +
			"(`user_id`, `topic_id`, `client_ip`, `content`, `created_at`, `updated_at`)" +
			"VALUES " +
			"(?, ?, ?, ?, ?, ?)"),
		comment.UID,
		comment.AID,
		comment.ClientIP,
		comment.Content,
		comment.AddTime,
		comment.AddTime,
	)
	if util.CheckError(err, "回复失败") {
		return
	}
	cid, err := rows.LastInsertId()
	comment.ID = uint64(cid)
}

func sqlGetCommentsBaseByList(db *sql.DB, redisDB *redis.Client, commentsList []uint64) (items []CommentBase) {
	var err error
	var rows *sql.Rows
	var commentListStr []string
	logger := util.GetLogger()
	defer rowsClose(rows)

	if len(commentsList) == 0 {
		logger.Warning("SQLArticleGetByList: Can't process the article list empty")
		return
	}

	for _, v := range commentsList {
		commentListStr = append(commentListStr, strconv.FormatInt(int64(v), 10))
	}
	qFieldList := []string{
		"id", "user_id", "topic_id", "content",
		"number", "created_at",
	}
	sql := fmt.Sprintf("select %s from reply where id in (%s)",
		strings.Join(qFieldList, ","),
		strings.Join(commentListStr, ","))

	rows, err = db.Query(sql)
	if err != nil {
		logger.Errorf("Query failed,err:%v", err)
		return
	}

	for rows.Next() {
		item := CommentBase{}
		err = rows.Scan(
			&item.ID, &item.UID, &item.AID, &item.Content,
			&item.Number, &item.AddTime,
		)
		if err != nil {
			logger.Errorf("Scan failed,err:%v", err)
			continue
		}
		items = append(items, item)
	}
	return
}
func (cb *CommentBase) toComment(db *sql.DB, redisDB *redis.Client, tz int) Comment {
	c := Comment{
		CommentBase: *cb,
	}
	c.AddTimeFmt = util.TimeFmt(cb.AddTime, time.RFC3339, tz)

	// 预防XSS漏洞
	c.ContentFmt = template.HTML(util.ContentFmt(cb.Content))

	c.UserName = GetUserNameByID(db, redisDB, cb.UID)
	c.Avatar = GetAvatarByID(db, redisDB, cb.UID)
	return c
}

func sqlCommentListByTopicID(db *sql.DB, redisDB *redis.Client, topicID uint64, tz int) (comments []Comment, err error) {
	var rows *sql.Rows
	defer rowsClose(rows)
	logger := util.GetLogger()

	rows, err = db.Query("SELECT id FROM reply where topic_id = ?", topicID)
	if err != nil {
		logger.Errorf("Query failed,err:%v", err)
		return
	}

	var commentList []uint64
	for rows.Next() {
		var item uint64
		err = rows.Scan(&item)
		if err != nil {
			logger.Errorf("Scan failed,err:%v", err)
			continue
		}
		commentList = append(commentList, item)
	}
	baseComments := sqlGetCommentsBaseByList(db, redisDB, commentList)
	for _, bc := range baseComments {
		comments = append(comments, bc.toComment(db, redisDB, tz))
	}
	return
}

// SQLGetCommentByID 获取一条评论
func SQLGetCommentByID(db *sql.DB, redisDB *redis.Client, cid uint64, tz int) Comment {
	logger := util.GetLogger()
	comments := sqlGetCommentsBaseByList(db, redisDB, []uint64{cid})
	if len(comments) == 0 {
		logger.Warningf("Error get comment(%d)", cid)
		return Comment{}
	}
	return comments[0].toComment(db, redisDB, tz)
}

// SQLCommentListByPage 获取帖子的所有评论
func SQLCommentListByPage(db *sql.DB, redisDB *redis.Client, topicID uint64, tz int) CommentPageInfo {
	var items []CommentListItem
	var hasPrev, hasNext bool
	var firstKey, lastKey uint64
	var rows *sql.Rows
	var err error
	defer rowsClose(rows)
	logger := util.GetLogger()

	comments, err := sqlCommentListByTopicID(db, redisDB, topicID, tz)
	if err != nil {
		logger.Errorf("Query comments failed for %d", topicID)
	}
	for _, c := range comments {
		item := CommentListItem{
			Comment: c,
		}
		items = append(items, item)
	}

	return CommentPageInfo{
		Items:    items,
		HasPrev:  hasPrev,
		HasNext:  hasNext,
		FirstKey: firstKey,
		LastKey:  lastKey,
	}
}

// SQLCommentList 获取在数据库中存储的评论
func SQLCommentList(db *sql.DB, redisDB *redis.Client, topicID, start uint64, btnAct string, limit, tz int) CommentPageInfo {
	var items []CommentListItem
	var hasPrev, hasNext bool
	var firstKey, lastKey uint64
	var rows *sql.Rows
	var err error
	logger := util.GetLogger()
	if btnAct == "" || btnAct == "next" {
		rows, err = db.Query(
			("SELECT id, user_id, topic_id, content, created_at " +
				" FROM  reply WHERE topic_id = ? And id > ?" +
				" ORDER BY id limit ?"),
			topicID,
			start,
			limit,
		)
	} else if btnAct == "prev" {
		rows, err = db.Query(
			("SELECT id, user_id, topic_id, content, created_at " +
				" FROM  reply WHERE topic_id = ? And id <= ?" +
				" ORDER BY id limit ?"),
			topicID,
			start,
			limit,
		)
	} else {
		logger.Error("Get wrond btn", btnAct)
	}
	defer func() {
		if rows != nil {
			rows.Close() //可以关闭掉未scan连接一直占用
		}
	}()
	for rows.Next() {
		item := CommentListItem{}
		err = rows.Scan(&item.ID, &item.UID, &item.AID, &item.Content, &item.AddTime)
		item.Avatar = GetAvatarByID(db, redisDB, item.UID)
		item.UserName = GetUserNameByID(db, redisDB, item.UID)

		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			continue
		}

		item.AddTimeFmt = util.TimeFmt(item.AddTime, "2006-01-02 15:04", tz)

		// 预防XSS漏洞
		item.ContentFmt = template.HTML(
			util.ContentFmt(item.Content))

		items = append(items, item)
	}

	if len(items) > 0 {
		firstKey = items[0].ID
		lastKey = items[len(items)-1].ID
		hasNext = true
		hasPrev = true
	}

	if start < uint64(limit) {
		hasPrev = false
	}
	if len(items) < limit {
		hasNext = false
	}

	return CommentPageInfo{
		Items:    items,
		HasPrev:  hasPrev,
		HasNext:  hasNext,
		FirstKey: firstKey,
		LastKey:  lastKey,
	}
}

func CommentGetByKey(db *youdb.DB, aid uint64, cid uint64) (Comment, error) {
	obj := Comment{}
	rs := db.Hget("article_comment:"+strconv.Itoa(int(aid)), youdb.I2b(cid))
	if rs.State == "ok" {
		json.Unmarshal(rs.Data[0], &obj)
		return obj, nil
	}
	return obj, errors.New(rs.State)
}

func CommentSetByKey(db *youdb.DB, aid uint64, cid uint64, obj Comment) error {
	jb, _ := json.Marshal(obj)
	return db.Hset("article_comment:"+strconv.Itoa(int(aid)), youdb.I2b(cid), jb)
}

func CommentDelByKey(db *youdb.DB, aid uint64, cid uint64) error {
	return db.Hdel("article_comment:"+strconv.Itoa(int(aid)), youdb.I2b(cid))
}

func CommentList(db *youdb.DB, cmd, tb, key string, limit, tz int) CommentPageInfo {
	var items []CommentListItem
	var citems []Comment
	userMap := map[uint64]UserMini{}
	var userKeys [][]byte
	var hasPrev, hasNext bool
	var firstKey, lastKey uint64

	keyStart := youdb.DS2b(key)
	if cmd == "hrscan" {
		rs := db.Hrscan(tb, keyStart, limit)
		if rs.State == "ok" {
			for i := len(rs.Data) - 2; i >= 0; i -= 2 {
				item := Comment{}
				json.Unmarshal(rs.Data[i+1], &item)
				citems = append(citems, item)
				userMap[item.UID] = UserMini{}
				userKeys = append(userKeys, youdb.I2b(item.UID))
			}
		}
	} else if cmd == "hscan" {
		rs := db.Hscan(tb, keyStart, limit)
		if rs.State == "ok" {
			for i := 0; i < (len(rs.Data) - 1); i += 2 {
				item := Comment{}
				json.Unmarshal(rs.Data[i+1], &item)
				citems = append(citems, item)
				userMap[item.UID] = UserMini{}
				userKeys = append(userKeys, youdb.I2b(item.UID))
			}
		}
	}

	if len(citems) > 0 {
		rs := db.Hmget("user", userKeys)
		if rs.State == "ok" {
			for i := 0; i < (len(rs.Data) - 1); i += 2 {
				item := UserMini{}
				json.Unmarshal(rs.Data[i+1], &item)
				userMap[item.ID] = item
			}
		}

		for _, citem := range citems {
			user := userMap[citem.UID]
			item := CommentListItem{
				Comment: Comment{
					CommentBase: CommentBase{
						ID:      citem.ID,
						AID:     citem.AID,
						UID:     citem.UID,
						AddTime: citem.AddTime,
					},
					Avatar:     user.Avatar,
					AddTimeFmt: util.TimeFmt(citem.AddTime, "2006-01-02 15:04", tz),
					ContentFmt: template.HTML(util.ContentFmt(citem.Content)),
				},
				Name: user.Name,
			}
			items = append(items, item)
			if firstKey == 0 {
				firstKey = item.ID
			}
			lastKey = item.ID
		}

		rs = db.Hrscan(tb, youdb.I2b(firstKey), 1)
		if rs.State == "ok" {
			hasPrev = true
		}
		rs = db.Hscan(tb, youdb.I2b(lastKey), 1)
		if rs.State == "ok" {
			hasNext = true
		}
	}

	return CommentPageInfo{
		Items:    items,
		HasPrev:  hasPrev,
		HasNext:  hasNext,
		FirstKey: firstKey,
		LastKey:  lastKey,
	}
}
