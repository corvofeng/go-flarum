package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"goyoubbs/model/flarum"
	"goyoubbs/util"

	"github.com/go-redis/redis/v7"
	"github.com/google/go-github/github"
)

// User store in database
type User struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	Nickname   string `json:"nickname"`
	Gender     string `json:"gender"`
	Flag       int    `json:"flag"`
	Avatar     string `json:"avatar"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	URL        string `json:"url"`
	Articles   uint64 `json:"articles"`
	Replies    uint64 `json:"replies"`
	RegTime    uint64 `json:"regtime"`
	About      string `json:"about"`
	Hidden     bool   `json:"hidden"`
	Session    string `json:"session"`
	Token      string `json:"token"`
	Reputation uint64 `json:"reputation"` // 声望值

	Preferences *flarum.Preferences
}

// UserListItem 用户信息
type UserListItem struct {
	User
	RegTimeFmt string `json:"regtime"`

	Notice    string `json:"notice"`
	NoticeNum uint64 `json:"noticenum"`

	LastPostTime     uint64
	LastReplyTime    uint64
	LastLoginTime    uint64
	LastPostTimeFmt  string `json:"lastposttime"`
	LastReplyTimeFmt string `json:"lastreplytime"`
	LastLoginTimeFmt string `json:"lastlogintime"`
}

// UserMini 简单用户
type UserMini struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// UserPageInfo 用户页面
type UserPageInfo struct {
	Items    []User `json:"items"`
	HasPrev  bool   `json:"hasprev"`
	HasNext  bool   `json:"hasnext"`
	FirstKey uint64 `json:"firstkey"`
	LastKey  uint64 `json:"lastkey"`
}

// toKey 作为Redis中的key值存储
func (user *User) toKey() string {
	return fmt.Sprintf("%d", user.ID)
}

// StrID 返回string类型的ID值
func (user *User) StrID() string {
	return fmt.Sprintf("%d", user.ID)
}

// IsValid 当前用户是否有效
func (user *User) IsValid() bool {
	return user.ID != 0
}

// SQLUserListByFlag 从数据库中查找用户列表
/*
 * db (*youdb.DB): TODO
 * cmd (TODO): TODO
 * tb (TODO): TODO
 * key (string): TODO
 * limit (int): TODO
 */
func SQLUserListByFlag(sqlDB *sql.DB, cmd, tb, key string, limit int) UserPageInfo {
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

// SQLUserGet 获取用户
// 当你不确定用户传来的是用户名还是用户id时, 可以调用该函数获取用户
func SQLUserGet(sqlDB *sql.DB, _userID string) (User, error) {
	var err error
	var user User
	var userID uint64
	logger := util.GetLogger()

	for true {
		// 如果通过用户名可以获取到用户, 那么马上退出并返回
		if user, err = SQLUserGetByName(sqlDB, _userID); err == nil {
			break
		}

		if userID, err = strconv.ParseUint(_userID, 10, 64); err != nil {
			logger.Error("Can't get user id for ", _userID)
			break
		}
		if user, err = SQLUserGetByID(sqlDB, userID); err != nil {
			logger.Error("Can't get user by err: ", err)
			break
		}
		break
	}
	return user, err
}

func sqlGetUserByList(db *sql.DB, redisDB *redis.Client, userIDList []uint64) (users []User) {
	var err error
	var rows *sql.Rows
	var userListStr []string
	logger := util.GetLogger()
	defer rowsClose(rows)

	if len(userIDList) == 0 {
		logger.Warning("sqlGetUserByList: Can't process the empty user list")
		return
	}

	for _, v := range userIDList {
		userListStr = append(userListStr, strconv.FormatInt(int64(v), 10))
	}
	qFieldList := []string{
		"id", "name", "nickname", "password", "reputation",
		"email", "avatar", "website",
		"description", "token", "created_at",
	}
	sql := fmt.Sprintf("select %s from user where id in (%s)",
		strings.Join(qFieldList, ","),
		strings.Join(userListStr, ","))

	rows, err = db.Query(sql)
	if err != nil {
		logger.Errorf("Query failed,err: %v", err)
		return
	}
	for rows.Next() {
		obj := User{}
		err = rows.Scan(
			&obj.ID,
			&obj.Name,
			&obj.Nickname,
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
			continue
		}
		users = append(users, obj)
	}

	return users
}

func (user *User) toUserListItem(db *sql.DB, redisDB *redis.Client, tz int) UserListItem {
	item := UserListItem{
		User: *user,
	}
	item.RegTimeFmt = util.TimeFmt(item.RegTime, util.TIME_FMT, tz)
	return item
}

// SQLUserGetByID 获取数据库用户
func SQLUserGetByID(sqlDB *sql.DB, uid uint64) (User, error) {
	obj := User{}
	users := sqlGetUserByList(sqlDB, nil, []uint64{uid})
	if len(users) == 0 {
		return obj, fmt.Errorf("Can't find user %d", uid)
	}
	return users[0], nil
}

// SQLUserGetByName 获取数据库中用户
func SQLUserGetByName(sqlDB *sql.DB, name string) (User, error) {
	var uid uint64
	obj := User{}
	logger := util.GetLogger()

	rows, err := sqlDB.Query("SELECT id FROM user WHERE name =  ?", name)
	defer rowsClose(rows)
	if err != nil {
		logger.Errorf("Query failed,err:%v", err)
		return obj, err
	}

	if rows.Next() {
		err = rows.Scan(&uid)
		if err != nil {
			logger.Errorf("Scan failed,err:%v", err)
			return obj, errors.New("No result")
		}
	} else {
		return obj, errors.New("No result")
	}

	return SQLUserGetByID(sqlDB, uid)
}

// SQLUserGetByEmail 获取数据库中用户
func SQLUserGetByEmail(sqlDB *sql.DB, email string) (User, error) {
	var uid uint64
	obj := User{}
	logger := util.GetLogger()

	rows, err := sqlDB.Query("SELECT id FROM user WHERE email =  ?", email)
	defer rowsClose(rows)
	if err != nil {
		logger.Errorf("Query failed,err:%v", err)
		return obj, err
	}

	if rows.Next() {
		err = rows.Scan(&uid)
		if err != nil {
			logger.Errorf("Scan failed,err:%v", err)
			return obj, errors.New("No result")
		}
	} else {
		return obj, errors.New("No result")
	}

	return SQLUserGetByID(sqlDB, uid)
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

// SQLGithubSync github用户同步信息
func (user *User) SQLGithubSync(sqlDB *sql.DB, gu *github.User) {
	logger := util.GetLogger()
	if user.Email != gu.GetEmail() {
		logger.Errorf("Wanna modify %s, but give %s", user.Email, gu.GetEmail())
		return
	}
	if user.About == "" {
		user.UpdateField(sqlDB, "description", gu.GetBio())
	}
	if user.URL == "" {
		user.UpdateField(sqlDB, "website", gu.GetBlog())
	}
	if user.Avatar == "" {
		user.UpdateField(sqlDB, "avatar", gu.GetAvatarURL())
	}
	*user, _ = SQLUserGetByID(sqlDB, user.ID)
}

// UpdateField 用户更新数据
func (user *User) UpdateField(sqlDB *sql.DB, field string, value string) {
	_, err := sqlDB.Exec(
		fmt.Sprintf("UPDATE `user` set %s=? where id=?", field),
		value,
		user.ID,
	)
	util.CheckError(err, fmt.Sprintf("更新用户信息%d (%s:%s)", user.ID, field, value))
}

// SQLGithubRegister github用户注册
func SQLGithubRegister(sqlDB *sql.DB, gu *github.User) (User, error) {
	logger := util.GetLogger()
	user := User{}
	row, err := sqlDB.Exec(
		("INSERT INTO `user` " +
			" (`name`, `email`, `is_email_confirmed`, `urlname`,`nickname`, `password`, `reputation`, `avatar`, `description`, `website`, `created_at`)" +
			" VALUES " +
			" (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"),
		gu.GetLogin(),
		gu.GetEmail(),
		"1",
		gu.GetLogin(),
		gu.GetName(),
		"NoPassWordForGithub",
		20, // 初始声望值20
		gu.GetAvatarURL(),
		gu.GetBio(),
		gu.GetBlog(),
		uint64(time.Now().UTC().Unix()),
	)
	if util.CheckError(err, "用户注册") {
		return user, err
	}
	uid, err := row.LastInsertId()
	if err != nil {
		logger.Error("Get insert id err", err)
	}
	logger.Infof("Create user %d-%s success", uid, gu.GetLogin())
	user, err = SQLUserGetByID(sqlDB, uint64(uid))
	if err != nil {
		logger.Error("Get user id err", err)
	}

	return user, nil
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
func (user *User) SaveAvatar(sqlDB *sql.DB, redisDB *redis.Client, avatar string) {
	logger := util.GetLogger()

	if user == nil {
		return
	}

	_, err := sqlDB.Exec("UPDATE user SET avatar = ? WHERE id = ?", avatar, user.ID)
	if err != nil {
		logger.Error("Set ", user, " avatar ", avatar, " failed!!")
		return
	}

	redisDB.HSet("avatar", user.toKey(), avatar)
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

	redisDB.HSet("avatar", user.toKey(), avatar)
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

	redisDB.HSet("username", user.toKey(), username)
	logger.Debugf("username not found for %d %s but we refresh!", user.ID, user.Name)
	return username
}

// GetPreference 获取用户定义配置
func (user *User) GetPreference(sqlDB *sql.DB, redisDB *redis.Client) {
	logger := util.GetLogger()
	rows, err := sqlDB.Query(
		"SELECT `preferences` FROM `user` WHERE id=?",
		user.ID,
	)
	defer rows.Close()
	if err != nil {
		logger.Error("Get preferences", err.Error())
		return
	}

	if rows.Next() {
		var data []byte
		rows.Scan(&data)
		err = json.Unmarshal(data, &user.Preferences)
		if err != nil {
			logger.Error("Load preferences", err.Error(), data)
			user.Preferences = &flarum.Preferences{}
		}
	}
	return
}

// SetPreference 更新用户配置信息
// 数据库中使用了blob的数据类型, 查看数据时, 需要进行转换:
//  SELECT CONVERT(`preferences` USING utf8) FROM `user`;
func (user *User) SetPreference(sqlDB *sql.DB, redisDB *redis.Client, preference flarum.Preferences) {
	logger := util.GetLogger()
	if user.Preferences == nil {
		logger.Warning("Can't process user with no preferences", user.ID, user.Name)
		return
	}

	data, err := json.Marshal(preference)
	if err != nil {
		logger.Error("Convert preferences to json error", err.Error(), user.Preferences)
		return
	}

	_, err = sqlDB.Exec(
		"UPDATE `user` "+
			"set preferences=?"+
			" where id=?",
		string(data),
		user.ID,
	)

	if err != nil {
		logger.Error("Update user preferences error", err.Error(), user.Preferences)
	}
	user.GetPreference(sqlDB, redisDB)
}

// RefreshCSRF 刷新CSRF token
func (user *User) RefreshCSRF(redisDB *redis.Client) string {
	t := util.GetNewToken()
	redisDB.HSet("csrf", user.toKey(), t)
	return t
}

// VerifyCSRFToken 确认用户CSRF token
func (user *User) VerifyCSRFToken(redisDB *redis.Client, token string) bool {
	rep, err := redisDB.HGet("csrf", user.toKey()).Result()
	if err != nil {
		util.GetLogger().Warningf("Can't get csrf token for user(%s): %s", user.ID, err)
		return false
	}
	return util.VerifyToken(token, rep)
}

// RefreshCache 刷新当前用户的信息
func (user *User) RefreshCache(redisDB *redis.Client) {
	user.CleareRedisCache(redisDB)
	user.CachedToRedis(redisDB)
	redisDB.HDel("avatar", user.toKey())
	redisDB.HDel("username", user.toKey())
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
