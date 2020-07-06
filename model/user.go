package model

import (
	"database/sql"
	"errors"
	"fmt"

	"goyoubbs/util"

	"github.com/ego008/youdb"
	"github.com/go-redis/redis/v7"
)

// User store in database
type User struct {
	ID            uint64 `json:"id"`
	Name          string `json:"name"`
	Gender        string `json:"gender"`
	Flag          int    `json:"flag"`
	Avatar        string `json:"avatar"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	URL           string `json:"url"`
	Articles      uint64 `json:"articles"`
	Replies       uint64 `json:"replies"`
	RegTime       uint64 `json:"regtime"`
	LastPostTime  uint64 `json:"lastposttime"`
	LastReplyTime uint64 `json:"lastreplytime"`
	LastLoginTime uint64 `json:"lastlogintime"`
	About         string `json:"about"`
	Notice        string `json:"notice"`
	NoticeNum     int    `json:"noticenum"`
	Hidden        bool   `json:"hidden"`
	Session       string `json:"session"`
	Token         string `json:"token"`
	Reputation    uint64 `json:"reputation"`
}

type UserMini struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type UserPageInfo struct {
	Items    []User `json:"items"`
	HasPrev  bool   `json:"hasprev"`
	HasNext  bool   `json:"hasnext"`
	FirstKey uint64 `json:"firstkey"`
	LastKey  uint64 `json:"lastkey"`
}

// SQLUserListByFlag 从数据库中查找用户列表
/*
 * db (*youdb.DB): TODO
 * cmd (TODO): TODO
 * tb (TODO): TODO
 * key (string): TODO
 * limit (int): TODO
 */
func SQLUserListByFlag(sqlDB *sql.DB, db *youdb.DB, cmd, tb, key string, limit int) UserPageInfo {
	var items []User
	// var keys [][]byte
	var hasPrev, hasNext bool
	var firstKey, lastKey uint64

	return UserPageInfo{
		Items:    items,
		HasPrev:  hasPrev,
		HasNext:  hasNext,
		FirstKey: firstKey,
		LastKey:  lastKey,
	}
}

// SQLUserGetByID 获取数据库用户
func SQLUserGetByID(db *sql.DB, uid uint64) (User, error) {
	obj := User{}
	logger := util.GetLogger()

	rows, err := db.Query(
		"SELECT id, name, password, reputation, email, avatar, website, description, token, created_at FROM user WHERE id = ?",
		uid,
	)
	defer rowsClose(rows)
	if err != nil {
		logger.Errorf("Query failed,err:%v", err)
		return obj, err
	}
	for rows.Next() {
		err = rows.Scan(
			&obj.ID,
			&obj.Name,
			&obj.Password,
			&obj.Reputation,
			&obj.Email,
			&obj.Avatar,
			&obj.URL,
			&obj.About,
			&obj.Token,
			&obj.RegTime,
		)

		if err != nil {
			logger.Errorf("Scan failed,err:%v", err)
			return obj, errors.New("No result")
		}
	}

	return obj, nil
}

// SQLUserGetByName 获取数据库中用户
func SQLUserGetByName(db *sql.DB, name string) (User, error) {
	obj := User{}
	logger := util.GetLogger()

	rows, err := db.Query(
		"SELECT id, name, password, reputation, email, avatar, website, token, created_at FROM user WHERE name =  ?",
		name)
	defer rowsClose(rows)
	if err != nil {
		logger.Errorf("Query failed,err:%v", err)
		return obj, err
	}
	if rows.Next() {
		err = rows.Scan(
			&obj.ID,
			&obj.Name,
			&obj.Password,
			&obj.Reputation,
			&obj.Email,
			&obj.Avatar,
			&obj.URL,
			&obj.Token,
			&obj.RegTime,
		)
		if err != nil {
			logger.Errorf("Scan failed,err:%v", err)
			return obj, errors.New("No result")
		}
	} else {
		return obj, errors.New("No result")
	}
	return obj, nil
}

// StrID 返回string类型的ID值
func (user *User) StrID() string {
	return fmt.Sprintf("%d", user.ID)
}

// SQLUserUpdate 更新用户信息
func (user *User) SQLUserUpdate(db *sql.DB) bool {
	_, err := db.Exec(
		"UPDATE `user` "+
			"set email=?,"+
			"description=?,"+
			"website=?"+
			" where id=?",
		user.Email,
		user.About,
		user.URL,
		user.ID,
	)
	if util.CheckError(err, "更新用户信息") {
		return false
	}
	return true
}

// SQLRegister 用户注册
func (user *User) SQLRegister(db *sql.DB) bool {
	row, err := db.Exec(
		("INSERT INTO `user` " +
			" (`name`, `email`, `urlname`, `password`, `reputation`, `avatar`)" +
			" VALUES " +
			" (?, ?,?, ?, ?, ?)"),
		user.Name,
		user.Name,
		user.Name,
		user.Password,
		20, // 初始声望值20
		"/static/avatar/3.jpg",
	)
	if util.CheckError(err, "用户注册") {
		return false
	}
	uid, err := row.LastInsertId()
	user.ID = uint64(uid)

	return true
}

// IsForbid 检查当前用户是否被禁用
func (user *User) IsForbid() bool {
	if user == nil {
		return true
	}
	// flag为0 并且 声望值较小
	if user.Flag == 0 && user.Reputation < 10 {
		return true
	}
	return false
}

// CanReply 检查当前用户是否可以回复帖子
func (user *User) CanReply() bool {
	return !user.IsForbid()
}

// CanCreateTopic 检查当前用户是否可以创建帖子
func (user *User) CanCreateTopic() bool {
	return !user.IsForbid()
}

// IsAdmin 检查当前用户是否为管理员
func (user *User) IsAdmin() bool {
	if user == nil {
		return false
	}
	if user.Reputation > 99 {
		return true
	}
	return false
}

// CanEdit 检查当前用户是否可以编辑帖子
func (user *User) CanEdit(aobjBase *ArticleBase) bool {
	if user == nil {
		return false
	}
	if aobjBase == nil {
		return false
	}
	return user.IsAdmin() || user.ID == aobjBase.UID
}

// SaveAvatar 更新用户头像
func (user *User) SaveAvatar(sqlDB *sql.DB, cntDB *youdb.DB, redisDB *redis.Client, avatar string) {
	logger := util.GetLogger()

	if user == nil {
		return
	}

	_, err := sqlDB.Exec("UPDATE user SET avatar = ? WHERE id = ?", avatar, user.ID)
	if err != nil {
		logger.Error("Set ", user, " avatar ", avatar, " failed!!")
		return
	}

	redisDB.HSet("avatar", fmt.Sprintf("%d", user.ID), avatar)
	logger.Notice("Refresh user avatar", user)
	return
}

// GetAvatarByID 获取用户头像
func GetAvatarByID(sqlDB *sql.DB, redisDB *redis.Client, uid uint64) string {
	var avatar string
	logger := util.GetLogger()

	rep, err := redisDB.HGet("avatar", fmt.Sprintf("%d", uid)).Result()
	if err != redis.Nil {
		return rep
	}

	user, err := SQLUserGetByID(sqlDB, uid)
	if util.CheckError(err, "查询用户") {
		return avatar
	}
	avatar = user.Avatar

	redisDB.HSet("avatar", fmt.Sprintf("%d", uid), avatar)
	logger.Debugf("avatar not found for %d %s but we refresh!", user.ID, user.Name)
	return avatar
}

// GetUserNameByID 获取用户名称
func GetUserNameByID(db *sql.DB, redisDB *redis.Client, uid uint64) string {
	var username string
	logger := util.GetLogger()

	rep, err := redisDB.HGet("username", fmt.Sprintf("%d", uid)).Result()
	if err != redis.Nil {
		return rep
	}

	user, err := SQLUserGetByID(db, uid)
	if util.CheckError(err, "查询用户") {
		return username
	}
	username = user.Name

	redisDB.HSet("username", fmt.Sprintf("%d", uid), username)
	logger.Debugf("username not found for %d %s but we refresh!", user.ID, user.Name)
	return username
}

// RefreshCSRF 刷新CSRF token
func (user *User) RefreshCSRF(redisDB *redis.Client) string {
	t := util.GetNewToken()
	redisDB.HSet("csrf", fmt.Sprintf("%d", user.ID), t)
	return t
}

// VerifyCSRFToken 确认用户CSRF token
func (user *User) VerifyCSRFToken(redisDB *redis.Client, token string) bool {
	rep, err := redisDB.HGet("csrf", fmt.Sprintf("%d", user.ID)).Result()
	if err != nil {
		util.GetLogger().Warningf("Can't get csrf token for user(%s): %s", user.ID, err)
		return false
	}
	return util.VerifyToken(token, rep)
}

// CachedToRedis 缓存当前用户的信息至Redis
func (user *User) CachedToRedis(redisDB *redis.Client) error {
	return rSet(redisDB, "user", user.StrID(), user)
}

// CleareRedisCache 缓存当前用户的信息至Redis
func (user *User) CleareRedisCache(redisDB *redis.Client) error {
	return rDel(redisDB, "user", user.StrID())
}

// RedisGetUserByID 从Redis中获取缓存的用户
func RedisGetUserByID(redisDB *redis.Client, uid string) (User, error) {
	user := User{}
	err := rGet(redisDB, "user", uid, &user)
	return user, err
}
