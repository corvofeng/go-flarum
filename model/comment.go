package model

import (
	"database/sql"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"

	"goyoubbs/util"

	"github.com/go-redis/redis/v7"
)

type (
	// CommentBase 会在数据库中保存的信息
	CommentBase struct {
		ID       uint64 `json:"id"`
		AID      uint64 `json:"aid"`
		UID      uint64 `json:"uid"`
		Number   uint64 `json:"number"`
		Content  string `json:"content"`
		ClientIP string `json:"clientip"`
		AddTime  uint64 `json:"addtime"`
	}

	// Comment 评论信息
	Comment struct {
		CommentBase
		UserName   string `json:"username"`
		Avatar     string `json:"avatar"`
		ContentFmt template.HTML
		AddTimeFmt string   `json:"addtimefmt"`
		Likes      []uint64 // 点赞的用户
	}

	// CommentListItem 页面中的评论
	CommentListItem struct {
		Comment

		Name string `json:"name"`
	}

	// CommentPageInfo 页面中显示的内容
	CommentPageInfo struct {
		Items    []CommentListItem `json:"items"`
		HasPrev  bool              `json:"hasprev"`
		HasNext  bool              `json:"hasnext"`
		FirstKey uint64            `json:"firstkey"`
		LastKey  uint64            `json:"lastkey"`
	}
)

// PreProcessUserMention 预处理用户的引用
// #14
func PreProcessUserMention(sqlDB *sql.DB, redisDB *redis.Client, tz int, userComment string) string {

	mentionDict := make(map[string]string)
	for _, mentionStr := range mentionRegexp.FindAllStringSubmatch(userComment, -1) {
		cid, err := strconv.ParseUint(mentionStr[2], 10, 64)
		if err != nil {
			util.GetLogger().Warning("Can't process mention", mentionStr[0])
			continue
		}
		comment, err := SQLGetCommentByID(sqlDB, redisDB, cid, tz)
		if err != nil {
			util.GetLogger().Warningf("Can't comment %d with error %v", cid, err)
			continue
		}
		user, err := SQLUserGetByName(sqlDB, mentionStr[1])
		replData := makeMention(mentionStr, comment, user)
		mentionDict[mentionStr[0]] = replData
	}

	newPost := replaceAllMentions(userComment, mentionDict)
	return newPost
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
		Likes:       cb.getUserLikes(db, redisDB),
	}
	c.AddTimeFmt = util.TimeFmt(cb.AddTime, time.RFC3339, tz)

	// 预防XSS漏洞
	c.ContentFmt = template.HTML(ContentFmt(cb.Content))

	c.UserName = GetUserNameByID(db, redisDB, cb.UID)
	c.Avatar = GetAvatarByID(db, redisDB, cb.UID)
	return c
}

func (cb *CommentBase) getUserLikes(db *sql.DB, redisDB *redis.Client) (likes []uint64) {
	var rows *sql.Rows
	defer rowsClose(rows)
	logger := util.GetLogger()

	sql := "SELECT reply_id, user_id FROM `reply_likes` where reply_id=?"
	rows, err := db.Query(sql, cb.ID)
	if err != nil {
		logger.Error("Can't get likes", err.Error())
		return
	}
	for rows.Next() {
		var cid uint64
		var uid uint64
		err = rows.Scan(&cid, &uid)
		if err != nil {
			logger.Errorf("Scan failed,err:%v", err)
		}

		likes = append(likes, uid)
	}
	return
}

func (comment *Comment) toCommentListItem(db *sql.DB, redisDB *redis.Client, tz int) CommentListItem {
	item := CommentListItem{
		Comment: *comment,
	}
	return item
}

func sqlCommentListByTopicID(db *sql.DB, redisDB *redis.Client, topicID uint64, limit uint64, tz int) (comments []Comment, err error) {
	var rows *sql.Rows
	defer rowsClose(rows)
	logger := util.GetLogger()

	rows, err = db.Query("SELECT id FROM reply where topic_id = ? LIMIT ?", topicID, limit)
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

func sqlCommentListByUserID(db *sql.DB, redisDB *redis.Client, userID uint64, limit uint64, tz int) (comments []Comment, err error) {
	var rows *sql.Rows
	defer rowsClose(rows)
	logger := util.GetLogger()

	rows, err = db.Query("SELECT id FROM `reply` where user_id = ? order by created_at desc limit ?", userID, limit)
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
func SQLGetCommentByID(db *sql.DB, redisDB *redis.Client, cid uint64, tz int) (Comment, error) {
	logger := util.GetLogger()
	comments := sqlGetCommentsBaseByList(db, redisDB, []uint64{cid})
	if len(comments) == 0 {
		logger.Debugf("Error get comment(%d)", cid)
		return Comment{}, fmt.Errorf("Can't find comment")
	}
	return comments[0].toComment(db, redisDB, tz), nil
}

// SQLCommentListByID 获取某条评论
func SQLCommentListByID(db *sql.DB, redisDB *redis.Client, commentID uint64, limit uint64, tz int) CommentPageInfo {
	var items []CommentListItem
	var hasPrev, hasNext bool
	var firstKey, lastKey uint64
	var err error
	logger := util.GetLogger()

	comment, err := SQLGetCommentByID(db, redisDB, commentID, tz)
	if err != nil {
		logger.Errorf("Query comments failed for cid(%d)", commentID)
	}
	items = append(items, comment.toCommentListItem(db, redisDB, tz))

	return CommentPageInfo{
		Items:    items,
		HasPrev:  hasPrev,
		HasNext:  hasNext,
		FirstKey: firstKey,
		LastKey:  lastKey,
	}
}

// SQLCommentListByList 获取某条评论
func SQLCommentListByList(db *sql.DB, redisDB *redis.Client, commentList []uint64, tz int) CommentPageInfo {
	var items []CommentListItem
	var hasPrev, hasNext bool
	var firstKey, lastKey uint64
	baseComments := sqlGetCommentsBaseByList(db, redisDB, commentList)
	for _, bc := range baseComments {
		c := bc.toComment(db, redisDB, tz)
		items = append(items, c.toCommentListItem(db, redisDB, tz))
	}

	return CommentPageInfo{
		Items:    items,
		HasPrev:  hasPrev,
		HasNext:  hasNext,
		FirstKey: firstKey,
		LastKey:  lastKey,
	}
}

// SQLCommentListByPage 获取帖子的所有评论
func SQLCommentListByPage(db *sql.DB, redisDB *redis.Client, topicID uint64, limit uint64, tz int) CommentPageInfo {
	var items []CommentListItem
	var hasPrev, hasNext bool
	var firstKey, lastKey uint64
	var err error
	logger := util.GetLogger()

	comments, err := sqlCommentListByTopicID(db, redisDB, topicID, limit, tz)
	if err != nil {
		logger.Errorf("Query comments failed for %d", topicID)
	}
	for _, c := range comments {
		items = append(items, c.toCommentListItem(db, redisDB, tz))
	}

	return CommentPageInfo{
		Items:    items,
		HasPrev:  hasPrev,
		HasNext:  hasNext,
		FirstKey: firstKey,
		LastKey:  lastKey,
	}
}

// SQLCommentListByUser 获取某个用户的帖子信息
func SQLCommentListByUser(db *sql.DB, redisDB *redis.Client, userID uint64, limit uint64, tz int) CommentPageInfo {
	var items []CommentListItem
	var hasPrev, hasNext bool
	var firstKey, lastKey uint64
	var err error
	logger := util.GetLogger()

	comments, err := sqlCommentListByUserID(db, redisDB, userID, limit, tz)
	if err != nil {
		logger.Errorf("Query comments failed for user %d", userID)
	}
	for _, c := range comments {
		items = append(items, c.toCommentListItem(db, redisDB, tz))
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
// TODO: deprecated
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
			ContentFmt(item.Content))

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

// CreateFlarumComment 创建flarum的评论
func (comment *Comment) CreateFlarumComment(db *sql.DB) (bool, error) {
	tx, err := db.Begin()
	logger := util.GetLogger()
	defer clearTransaction(tx)
	if err != nil {
		return false, err
	}
	if ok, err := comment.sqlCreateComment(tx); !ok {
		return false, err
	}

	logger.Debugf("Create comment %d success", comment.ID)

	if ok, err := comment.sqlUpdateNumber(tx); !ok {
		return false, err
	}

	logger.Debugf("Update comment number %d success", comment.ID)
	if err := tx.Commit(); err != nil {
		logger.Error("Create reply with error", err)
		return false, err
	}

	return true, nil
}

func (comment *Comment) sqlCreateComment(tx *sql.Tx) (bool, error) {
	row, err := tx.Exec(
		("INSERT INTO `reply` " +
			" (`user_id`, `topic_id`, `content`, created_at, updated_at, client_ip)" +
			" VALUES " +
			" (?, ?, ?, ?, ?, ?)"),
		comment.UID,
		comment.AID,
		comment.Content,
		comment.AddTime,
		comment.AddTime,
		comment.ClientIP,
	)
	if err != nil {
		return false, err
	}
	cid, err := row.LastInsertId()
	comment.ID = uint64(cid)

	return true, nil
}

func (comment *Comment) sqlUpdateNumber(tx *sql.Tx) (bool, error) {
	// 锁表
	logger := util.GetLogger()
	var lastReplyID uint64
	var lastReplyNumber uint64
	var replyCount uint64

	row, err := tx.Query(
		("SELECT reply.id, reply.number, t.reply_count" +
			" FROM " +
			" (SELECT last_post_id, reply_count FROM `topic` WHERE id = ? FOR UPDATE ) AS t" +
			" LEFT JOIN reply ON t.last_post_id = reply.id"),
		comment.AID,
	)
	if err != nil {
		return false, err
	}

	if row.Next() {
		row.Scan(&lastReplyID, &lastReplyNumber, &replyCount)
	} else {
		logger.Warningf("Can't get last post for topic %d", comment.AID)
		lastReplyID = 0
		lastReplyNumber = 0
		replyCount = 0
	}
	rowsClose(row) // 查询之后, 立刻关闭, 否则后面的语句无法执行

	logger.Debugf("Get last reply (%d,%d) for article: %d cnt: %d", lastReplyID, lastReplyNumber, comment.AID, replyCount)

	comment.Number = lastReplyNumber + 1
	replyCount = replyCount + 1
	_, err = tx.Exec(
		("UPDATE `topic` SET" +
			" last_post_id=?," +
			" reply_count=?" +
			" where id=?"),
		comment.ID,
		replyCount,
		comment.AID,
	)
	if err != nil {
		return false, err
	}

	_, err = tx.Exec(
		("UPDATE `reply` SET" +
			" number=?" +
			" where id=?"),
		comment.Number,
		comment.ID,
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

// DoLike 用户的点赞
func (comment *Comment) DoLike(db *sql.DB, redisDB *redis.Client, user *User, isLiked bool) {
	sql := ""
	if isLiked {
		sql = "INSERT INTO `reply_likes` (`reply_id`, `user_id`) VALUES (?, ?)"
	} else {

		sql = "DELETE FROM `reply_likes` WHERE `reply_likes`.`reply_id` = ? AND `reply_likes`.`user_id` = ?"
	}
	_, err := db.Exec(sql, comment.ID, user.ID)
	logger := util.GetLogger()
	if err != nil {
		logger.Warning("Can't do sql", sql, err.Error())
	}
}
