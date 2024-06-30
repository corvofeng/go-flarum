package model

import (
	"database/sql"
	"html/template"
	"strconv"

	"github.com/corvofeng/go-flarum/util"

	"github.com/go-redis/redis/v7"
	"gorm.io/gorm"
)

type (
	// Reply 会在数据库中保存的信息
	Reply struct {
		gorm.Model
		ID       uint64 `gorm:"primaryKey"`
		AID      uint64 `gorm:"column:topic_id"`
		UID      uint64 `gorm:"column:user_id"`
		Number   uint64 `json:"number"`
		Content  string `json:"content"`
		ClientIP string `json:"clientip"`
		AddTime  uint64 `json:"addtime"`
	}

	ReplyLikes struct {
		gorm.Model
		UserID  uint64 `gorm:"column:user_id;index"`
		ReplyID uint64 `gorm:"column:reply_id;index"`
	}

	// Comment 评论信息
	Comment struct {
		Reply
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
func PreProcessUserMention(gormDB *gorm.DB, redisDB *redis.Client, tz int, userComment string) string {

	mentionDict := make(map[string]string)
	for _, mentionStr := range mentionRegexp.FindAllStringSubmatch(userComment, -1) {
		cid, err := strconv.ParseUint(mentionStr[2], 10, 64)
		if err != nil {
			util.GetLogger().Warning("Can't process mention", mentionStr[0])
			continue
		}
		comment, err := SQLCommentByID(gormDB, redisDB, cid, tz)
		if err != nil {
			util.GetLogger().Warningf("Can't comment %d with error %v", cid, err)
			continue
		}
		user, err := SQLUserGetByName(gormDB, mentionStr[1])
		replData := makeMention(mentionStr, comment, user)
		mentionDict[mentionStr[0]] = replData
	}

	newPost := replaceAllMentions(userComment, mentionDict)
	return newPost
}

func sqlGetRepliesBaseByList(gormDB *gorm.DB, redisDB *redis.Client, repliesList []uint64) (items []Reply) {
	logger := util.GetLogger()
	result := gormDB.Find(&items, repliesList)
	if result.Error != nil {
		logger.Errorf("Can't get replies list by ", repliesList)
	}
	return
}

func (cb *Reply) toComment(gormDB *gorm.DB, redisDB *redis.Client, tz int) Comment {
	c := Comment{
		Reply: *cb,
		Likes: cb.getUserLikes(gormDB, redisDB),
	}
	c.AddTimeFmt = cb.CreatedAt.String()

	// 预防XSS漏洞
	c.ContentFmt = template.HTML(ContentFmt(cb.Content))

	c.UserName = GetUserNameByID(gormDB, redisDB, cb.UID)
	c.Avatar = GetAvatarByID(gormDB, redisDB, cb.UID)
	return c
}

func (cb *Reply) getUserLikes(gormDB *gorm.DB, redisDB *redis.Client) (likes []uint64) {
	rows, _ := gormDB.Model(&ReplyLikes{}).Where("reply_id = ?", cb.ID).Rows()
	defer rows.Close()
	for rows.Next() {
		var r ReplyLikes
		gormDB.ScanRows(rows, &r)
		likes = append(likes, r.UserID)
	}
	return
}

func (cb *Reply) getUserComments(gormDB *gorm.DB, redisDB *redis.Client) (comments []uint64) {
	rows, _ := gormDB.Model(&Reply{}).Where("user_id = ?", cb.UID).Rows()
	defer rows.Close()
	for rows.Next() {
		var r Reply
		gormDB.ScanRows(rows, &r)
		comments = append(comments, r.ID)
	}
	return
}

func (comment *Comment) toCommentListItem(redisDB *redis.Client, tz int) CommentListItem {
	item := CommentListItem{
		Comment: *comment,
	}
	return item
}

func sqlCommentListByTopicID(gormDB *gorm.DB, redisDB *redis.Client, topicID uint64, limit uint64, tz int) (comments []Comment, err error) {
	var rows *sql.Rows
	defer rowsClose(rows)

	var baseComments []Reply
	gormDB.Order("number asc").Where("topic_id = ?", topicID).Limit(int(limit)).Find(&baseComments)
	for _, bc := range baseComments {
		comments = append(comments, bc.toComment(gormDB, redisDB, tz))
	}
	return
}

func sqlCommentListByUserID(gormDB *gorm.DB, redisDB *redis.Client, userID uint64, limit uint64, tz int) (comments []Comment, err error) {
	var baseComments []Reply
	gormDB.Where("user_id = ?", userID).Limit(int(limit)).Find(&baseComments)

	for _, bc := range baseComments {
		comments = append(comments, bc.toComment(gormDB, redisDB, tz))
	}
	return
}

// SQLCommentByID 获取一条评论
func SQLCommentByID(gormDB *gorm.DB, redisDB *redis.Client, cid uint64, tz int) (Comment, error) {
	logger := util.GetLogger()
	var c Reply
	result := gormDB.First(&c, cid)

	if result.Error != nil {
		logger.Error("Can't find commet with error", result.Error)
		return Comment{}, result.Error
	}
	return c.toComment(gormDB, redisDB, tz), nil
}

// SQLCommentListByCID 获取某条评论
func SQLCommentListByCID(gormDB *gorm.DB, redisDB *redis.Client, commentID uint64, limit uint64, tz int) CommentPageInfo {
	var items []CommentListItem
	var hasPrev, hasNext bool
	var firstKey, lastKey uint64
	var err error
	logger := util.GetLogger()

	comment, err := SQLCommentByID(gormDB, redisDB, commentID, tz)
	if err != nil {
		logger.Errorf("Query comments failed for cid(%d)", commentID)
	}
	items = append(items, comment.toCommentListItem(redisDB, tz))

	return CommentPageInfo{
		Items:    items,
		HasPrev:  hasPrev,
		HasNext:  hasNext,
		FirstKey: firstKey,
		LastKey:  lastKey,
	}
}

// SQLCommentListByList 获取某条评论
func SQLCommentListByList(gormDB *gorm.DB, redisDB *redis.Client, commentList []uint64, tz int) CommentPageInfo {
	var items []CommentListItem
	var hasPrev, hasNext bool
	var firstKey, lastKey uint64

	baseComments := sqlGetRepliesBaseByList(gormDB, redisDB, commentList)
	for _, bc := range baseComments {
		c := bc.toComment(gormDB, redisDB, tz)
		items = append(items, c.toCommentListItem(redisDB, tz))
	}

	return CommentPageInfo{
		Items:    items,
		HasPrev:  hasPrev,
		HasNext:  hasNext,
		FirstKey: firstKey,
		LastKey:  lastKey,
	}
}

// SQLCommentListByTopic 获取帖子的所有评论
func SQLCommentListByTopic(gormDB *gorm.DB, redisDB *redis.Client, topicID uint64, limit uint64, tz int) CommentPageInfo {
	var items []CommentListItem
	var hasPrev, hasNext bool
	var firstKey, lastKey uint64
	var err error
	logger := util.GetLogger()

	comments, err := sqlCommentListByTopicID(gormDB, redisDB, topicID, limit, tz)
	if err != nil {
		logger.Errorf("Query comments failed for %d", topicID)
	}
	for _, c := range comments {
		items = append(items, c.toCommentListItem(redisDB, tz))
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
func SQLCommentListByUser(gormDB *gorm.DB, redisDB *redis.Client, userID uint64, limit uint64, tz int) CommentPageInfo {
	var items []CommentListItem
	var hasPrev, hasNext bool
	var firstKey, lastKey uint64
	var err error
	logger := util.GetLogger()

	comments, err := sqlCommentListByUserID(gormDB, redisDB, userID, limit, tz)
	if err != nil {
		logger.Errorf("Query comments failed for user %d", userID)
	}
	for _, c := range comments {
		items = append(items, c.toCommentListItem(redisDB, tz))
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
func (comment *Comment) CreateFlarumComment(gormDB *gorm.DB) (bool, error) {
	logger := util.GetLogger()

	tx := gormDB.Begin()
	defer clearGormTransaction(tx)

	var topic Topic
	result := gormDB.First(&topic, &comment.AID)
	if result.Error != nil {
		logger.Error("Can't find topic with error", result.Error)
		return false, result.Error
	}

	var lastComment Reply
	result = gormDB.First(&lastComment, topic.LastPostID)
	if result.Error != nil {
		logger.Error("Can't find last commet with error", result.Error)
		return false, result.Error
	}
	comment.Number = lastComment.Number + 1

	result = tx.Create(&comment.Reply)
	if result.Error != nil {
		return false, result.Error
	}
	topic.CommentCount += 1
	topic.LastPostID = comment.ID
	topic.LastPostUserID = comment.UID
	topic.LastPostAt = comment.CreatedAt

	tx.Save(&topic)
	logger.Debugf("Update comment number %d success", comment.ID)

	result = tx.Commit()
	if result.Error != nil {
		logger.Error("Create reply with error", result.Error)
		return false, result.Error
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

func (comment *Comment) sqlUpdateComment(tx *sql.Tx, newContent string) (bool, error) {
	comment.Content = newContent
	_, err := tx.Exec("UPDATE `reply`"+
		" set content = ?, updated_at = ?"+
		" where id = ?",
		comment.Content,
		util.TimeNow(),
		comment.ID,
	)
	if util.CheckError(err, "更新评论内容") {
		return false, err
	}
	return true, nil
}

func (comment *Comment) sqlCreateHistory(tx *sql.Tx, newContent string, uID uint64) (bool, error) {
	_, err := tx.Exec(
		("INSERT INTO `history` " +
			" (`user_id`, `reply_id`, `topic_id`, `content`, created_at)" +
			" VALUES " +
			" (?, ?, ?, ?, ?)"),
		uID,
		comment.ID,
		comment.AID,
		comment.Content,
		util.TimeNow(),
	)
	if err != nil {
		util.GetLogger().Errorf("Can't create history because of %s", err)
		return false, err
	}
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
func (comment *Comment) DoLike(gormDB *gorm.DB, redisDB *redis.Client, user *User, isLiked bool) {
	if isLiked {
		rl := ReplyLikes{UserID: user.ID, ReplyID: comment.ID}
		gormDB.Create(&rl)
	} else {
		gormDB.Unscoped().Where("user_id = ? and reply_id = ?", user.ID, comment.ID).Delete(&ReplyLikes{})
	}
}
