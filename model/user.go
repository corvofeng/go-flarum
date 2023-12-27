package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"zoe/model/flarum"
	"zoe/util"

	"github.com/go-redis/redis/v7"
	"github.com/google/go-github/github"

	"gorm.io/gorm"
)

// User store in database
type User struct {
	gorm.Model
	ID       uint64 `gorm:"primaryKey"`
	Name     string `json:"name" gorm:"index:idx_name,unique"`
	Nickname string `json:"nickname"`
	Gender   string `json:"gender"`
	Flag     int    `json:"flag"`
	Avatar   string `json:"avatar"`
	Password string `json:"password"`
	Email    string `json:"email"`
	URL      string `json:"url"`
	Articles uint64 `json:"articles"`
	Replies  uint64 `json:"replies"`
	RegTime  uint64 `json:"regtime"`
	About    string `json:"about"`
	Hidden   bool   `json:"hidden"`
	Session  string `json:"session"`
	Token    string `json:"token"`

	Description string
	WebSite     string
	Reputation  uint64 `json:"reputation"` // 声望值

	// Preferences *flarum.Preferences `gorm:"foreignKey:PreferencesRefer"`
	Preferences []byte
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

// SQLUserGet 获取用户
// 当你不确定用户传来的是用户名还是用户id时, 可以调用该函数获取用户
func SQLUserGet(gormDB *gorm.DB, _userID string) (User, error) {
	var err error
	var user User
	var userID uint64
	logger := util.GetLogger()

	for {
		// 如果通过用户名可以获取到用户, 那么马上退出并返回
		// if user, err = SQLUserGetByName(gormDB, _userID); err == nil {
		// 	break
		// }

		if userID, err = strconv.ParseUint(_userID, 10, 64); err != nil {
			logger.Error("Can't get user id for ", _userID)
			break
		}
		if user, err = SQLUserGetByID(gormDB, userID); err != nil {
			logger.Error("Can't get user by err: ", err)
			break
		}
		break //lint:ignore SA4004 ignore this!
	}
	return user, err
}

// func sqlGetUserByList(db *sql.DB, redisDB *redis.Client, userIDList []uint64) (users []User) {
// 	var err error
// 	var rows *sql.Rows
// 	var userListStr []string
// 	logger := util.GetLogger()
// 	defer rowsClose(rows)
// 	if len(userIDList) == 0 {
// 		logger.Warning("sqlGetUserByList: Can't process the empty user list")
// 		return
// 	}
// 	for _, v := range userIDList {
// 		userListStr = append(userListStr, strconv.FormatInt(int64(v), 10))
// 	}
// 	qFieldList := []string{
// 		"id", "name", "nickname", "password", "reputation",
// 		"email", "avatar", "website",
// 		"description", "token", "created_at",
// 	}
// 	sql := fmt.Sprintf("select %s from user where id in (%s)",
// 		strings.Join(qFieldList, ","),
// 		strings.Join(userListStr, ","))
// 	rows, err = db.Query(sql)
// 	if err != nil {
// 		logger.Errorf("Query failed,err: %v", err)
// 		return
// 	}
// 	for rows.Next() {
// 		obj := User{}
// 		err = rows.Scan(
// 			&obj.ID,
// 			&obj.Name,
// 			&obj.Nickname,
// 			&obj.Password,
// 			&obj.Reputation,
// 			&obj.Email,
// 			&obj.Avatar,
// 			&obj.URL,
// 			&obj.About,
// 			&obj.Token,
// 			&obj.RegTime,
// 		)
// 		if err != nil {
// 			logger.Errorf("Scan failed,err:%v", err)
// 			continue
// 		}
// 		users = append(users, obj)
// 	}

// 	return users
// }

// func (user *User) toUserListItem(db *sql.DB, redisDB *redis.Client, tz int) UserListItem {
// 	item := UserListItem{
// 		User: *user,
// 	}
// 	item.RegTimeFmt = util.TimeFmt(item.RegTime, util.TIME_FMT, tz)
// 	return item
// }

// SQLUserGetByID 获取数据库用户
func SQLUserGetByID(gormDB *gorm.DB, uid uint64) (User, error) {
	user := User{}
	result := gormDB.First(&user, uid)
	return user, result.Error
}

// SQLUserGetByName 获取数据库中用户
func SQLUserGetByName(gormDB *gorm.DB, name string) (User, error) {
	user := User{}
	result := gormDB.Where("name = ?", name).First(&user)
	return user, result.Error
}

// SQLUserGetByEmail 获取数据库中用户
func SQLUserGetByEmail(gormDB *gorm.DB, email string) (User, error) {
	user := User{}
	result := gormDB.Where("email = ?", email).First(&user)
	return user, result.Error
}

// SQLUserUpdate 更新用户信息
// func (user *User) SQLUserUpdate(db *sql.DB) bool {
// 	_, err := db.Exec(
// 		"UPDATE `user` "+
// 			"set email=?,"+
// 			"description=?,"+
// 			"website=?"+
// 			" where id=?",
// 		user.Email,
// 		user.About,
// 		user.URL,
// 		user.ID,
// 	)
// 	if util.CheckError(err, "更新用户信息") {
// 		return false
// 	}
// 	return true
// }

func SQLUserRegister(gormDB *gorm.DB, name, email, password string) (User, error) {
	user := User{
		Name:       name,
		Email:      email,
		Password:   password,
		Reputation: 20,
		Avatar:     "/static/avatar/3.jpg",
	}
	result := gormDB.Create(&user)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

// SQLGithubSync github用户同步信息
func (user *User) SQLGithubSync(gormDB *gorm.DB, gu *github.User) {
	logger := util.GetLogger()
	if user.Email != gu.GetEmail() {
		logger.Errorf("Wanna modify %s, but give %s", user.Email, gu.GetEmail())
		return
	}

	if user.About == "" {
		gormDB.Model(user).Update("description", gu.GetBio())
	}
	if user.URL == "" {
		gormDB.Model(user).Update("website", gu.GetBlog())
	}
	if user.Avatar == "" {
		gormDB.Model(user).Update("avatar", gu.GetAvatarURL())
	}
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
func SQLGithubRegister(gormDB *gorm.DB, gu *github.User) (User, error) {
	user := User{
		Name:        gu.GetLogin(),
		Email:       gu.GetEmail(),
		Reputation:  20,
		Password:    "NoPassWordForGithub",
		Avatar:      gu.GetAvatarURL(),
		Description: gu.GetBio(),
		WebSite:     gu.GetBlog(),
	}
	result := gormDB.Create(&user)
	logger := util.GetLogger()

	if result.Error != nil {
		logger.Error("Can't creat user for", gu.GetLogin())
		return User{}, result.Error
	}

	logger.Infof("Create user %d-%s success", user.ID, gu.GetLogin())
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
func (user *User) CanEdit(aobjBase *Topic) bool {
	if user == nil {
		return false
	}
	if aobjBase == nil {
		return false
	}
	return user.IsAdmin() || user.ID == aobjBase.UserID
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
}

// GetAvatarByID 获取用户头像
func GetAvatarByID(gormDB *gorm.DB, redisDB *redis.Client, uid uint64) string {
	var avatar string
	logger := util.GetLogger()

	rep, err := redisDB.HGet("avatar", fmt.Sprintf("%d", uid)).Result()
	if err != redis.Nil {
		return rep
	}

	user, err := SQLUserGetByID(gormDB, uid)
	if util.CheckError(err, "查询用户") {
		return avatar
	}
	avatar = user.Avatar

	redisDB.HSet("avatar", user.toKey(), avatar)
	logger.Debugf("avatar not found for %d %s but we refresh!", user.ID, user.Name)
	return avatar
}

// GetUserNameByID 获取用户名称
func GetUserNameByID(gormDB *gorm.DB, redisDB *redis.Client, uid uint64) string {
	var username string
	logger := util.GetLogger()

	rep, err := redisDB.HGet("username", fmt.Sprintf("%d", uid)).Result()
	if err != redis.Nil {
		return rep
	}

	user, err := SQLUserGetByID(gormDB, uid)
	if util.CheckError(err, "查询用户") {
		return username
	}
	username = user.Name

	redisDB.HSet("username", user.toKey(), username)
	logger.Debugf("username not found for %d %s but we refresh!", user.ID, user.Name)
	return username
}

// SetPreference 更新用户配置信息
// 数据库中使用了blob的数据类型, 查看数据时, 需要进行转换:
//
//	SELECT CONVERT(`preferences` USING utf8) FROM `user`;
func (user *User) SetPreference(gormDB *gorm.DB, redisDB *redis.Client, preference flarum.Preferences) {
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

	result := gormDB.Model(user).Update("preferences", data)
	if result.Error != nil {
		logger.Error("Update user preferences error", result.Error, user.Preferences)
	}
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
