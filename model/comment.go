package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"strconv"

	"goyoubbs/util"

	"github.com/ego008/youdb"
	"github.com/go-redis/redis/v7"
)

// Comment 评论信息
type Comment struct {
	ID       uint64 `json:"id"`
	Aid      uint64 `json:"aid"`
	UID      uint64 `json:"uid"`
	Number   uint64 `json:"number"`
	Content  string `json:"content"`
	ClientIp string `json:"clientip"`
	AddTime  uint64 `json:"addtime"`
}

// CommentListItem 页面中的评论
type CommentListItem struct {
	ID         uint64 `json:"id"`
	Aid        uint64 `json:"aid"`
	UID        uint64 `json:"uid"`
	Number     uint64 `json:"number"`
	Name       string `json:"name"`
	UserName   string `json:"username"`
	Avatar     string `json:"avatar"`
	Content    string `json:"content"`
	ContentFmt template.HTML
	AddTime    uint64 `json:"addtime"`
	AddTimeFmt string `json:"addtimefmt"`
}

type CommentPageInfo struct {
	Items    []CommentListItem `json:"items"`
	HasPrev  bool              `json:"hasprev"`
	HasNext  bool              `json:"hasnext"`
	FirstKey uint64            `json:"firstkey"`
	LastKey  uint64            `json:"lastkey"`
}

// SQLSaveComment 在数据库中存储评论
func (comment *Comment) SQLSaveComment(db *sql.DB) {
	rows, err := db.Exec(
		("INSERT INTO `reply`" +
			"(`user_id`, `topic_id`, `client_ip`, `content`, `created_at`, `updated_at`)" +
			"VALUES " +
			"(?, ?, ?, ?, ?, ?)"),
		comment.UID,
		comment.Aid,
		comment.ClientIp,
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

// SQLCommentListByPage 获取帖子的所有评论
func SQLCommentListByPage(db *sql.DB, redisDB *redis.Client, topicID uint64, tz int) CommentPageInfo {
	var items []CommentListItem
	var hasPrev, hasNext bool
	var firstKey, lastKey uint64
	var rows *sql.Rows
	var err error
	rows, err = db.Query(
		("SELECT id, user_id, topic_id, content, created_at " +
			" FROM  reply WHERE topic_id = ?"),
		topicID,
	)
	defer func() {
		if rows != nil {
			rows.Close() //可以关闭掉未scan连接一直占用
		}
	}()

	for rows.Next() {
		item := CommentListItem{}
		err = rows.Scan(&item.ID, &item.UID, &item.Aid, &item.Content, &item.AddTime)
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
		err = rows.Scan(&item.ID, &item.UID, &item.Aid, &item.Content, &item.AddTime)
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
				ID:         citem.ID,
				Aid:        citem.Aid,
				UID:        citem.UID,
				Name:       user.Name,
				Avatar:     user.Avatar,
				AddTime:    citem.AddTime,
				AddTimeFmt: util.TimeFmt(citem.AddTime, "2006-01-02 15:04", tz),
				ContentFmt: template.HTML(util.ContentFmt(citem.Content)),
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
