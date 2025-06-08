package model

import (
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
		Avatar     string `json:"avatar"`
		ContentFmt template.HTML
		Likes      []uint64 // 点赞的用户
	}
)

// PreProcessUserMention 预处理用户的引用
// #14
func PreProcessUserMention(gormDB *gorm.DB, redisDB *redis.Client, tz int, userComment string) string {
	logger := util.GetLogger()

	mentionDict := make(map[string]string)
	for _, mentionStr := range mentionRegexp.FindAllStringSubmatch(userComment, -1) {
		cid, err := strconv.ParseUint(mentionStr[2], 10, 64)
		if err != nil {
			logger.Warning("Can't process mention", mentionStr[0])
			continue
		}
		comment, err := SQLCommentByID(gormDB, redisDB, cid, tz)
		if err != nil {
			logger.Warningf("Can't comment %d with error %v", cid, err)
			continue
		}
		user, err := SQLUserGetByName(gormDB, mentionStr[1])
		if err != nil {
			logger.Warningf("can't get user `%d` with error %v", mentionStr[1], err)
			continue
		}
		logger.Debugf("Mention %s, comment %d, user %d", mentionStr[0], comment.ID, user.ID)
		replData := makeMention(mentionStr, comment, user)
		mentionDict[mentionStr[0]] = replData
	}

	newPost := replaceAllMentions(userComment, mentionDict)
	return newPost
}

func sqlGetRepliesBaseByList(gormDB *gorm.DB, redisDB *redis.Client, repliesList []uint64, tz int) (comments []Comment, err error) {
	var baseComments []Reply
	err = gormDB.Find(&baseComments, repliesList).Error
	for _, bc := range baseComments {
		comments = append(comments, bc.toComment(gormDB, redisDB, tz))
	}
	return
}

func (cb *Reply) toComment(gormDB *gorm.DB, redisDB *redis.Client, tz int) Comment {
	c := Comment{
		Reply: *cb,
		Likes: cb.getUserLikes(gormDB, redisDB),
	}

	// 预防XSS漏洞
	c.ContentFmt = template.HTML(ContentFmt(cb.Content))
	c.Avatar = GetAvatarByID(gormDB, redisDB, cb.UID)
	return c
}

func (cb *Reply) getUserLikes(gormDB *gorm.DB, redisDB *redis.Client) (likes []uint64) {
	gormDB.Model(&ReplyLikes{}).Where("reply_id = ?", cb.ID).Pluck("user_id", &likes)
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

func sqlCommentListByTopicID(gormDB *gorm.DB, redisDB *redis.Client, topicID uint64, limit uint64, tz int) (comments []Comment, err error) {
	var replys []Reply
	err = gormDB.Order("number asc").Where("topic_id = ?", topicID).Limit(int(limit)).Find(&replys).Error
	for _, bc := range replys {
		comments = append(comments, bc.toComment(gormDB, redisDB, tz))
	}
	return
}

func sqlCommentListByUserID(gormDB *gorm.DB, redisDB *redis.Client, userID uint64, limit uint64, tz int) (comments []Comment, err error) {
	var baseComments []Reply
	err = gormDB.Where("user_id = ?", userID).Limit(int(limit)).Find(&baseComments).Error
	for _, bc := range baseComments {
		comments = append(comments, bc.toComment(gormDB, redisDB, tz))
	}
	return
}

// SQLCommentByID 获取一条评论
func SQLCommentByID(gormDB *gorm.DB, redisDB *redis.Client, cid uint64, tz int) (Comment, error) {
	var c Reply
	err := gormDB.First(&c, cid).Error
	return c.toComment(gormDB, redisDB, tz), err
}

// SQLCommentListByCID 获取某条评论
func SQLCommentListByCID(gormDB *gorm.DB, redisDB *redis.Client, commentID uint64, limit uint64, tz int) ([]Comment, error) {
	comment, err := SQLCommentByID(gormDB, redisDB, commentID, tz)
	return []Comment{comment}, err
}

// SQLCommentListByList 获取某条评论
func SQLCommentListByList(gormDB *gorm.DB, redisDB *redis.Client, commentList []uint64, tz int) ([]Comment, error) {
	return sqlGetRepliesBaseByList(gormDB, redisDB, commentList, tz)
}

// SQLCommentListByTopic 获取帖子的所有评论
func SQLCommentListByTopic(gormDB *gorm.DB, redisDB *redis.Client, topicID uint64, limit uint64, tz int) ([]Comment, error) {
	return sqlCommentListByTopicID(gormDB, redisDB, topicID, limit, tz)
}

// SQLCommentListByUser 获取某个用户的帖子信息
func SQLCommentListByUser(gormDB *gorm.DB, redisDB *redis.Client, userID uint64, limit uint64, tz int) ([]Comment, error) {
	return sqlCommentListByUserID(gormDB, redisDB, userID, limit, tz)
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

	tx.Save(&topic)
	logger.Debugf("Update comment number %d success", comment.ID)

	result = tx.Commit()
	if result.Error != nil {
		logger.Error("Create reply with error", result.Error)
		return false, result.Error
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
