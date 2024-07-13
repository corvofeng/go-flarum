package model

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/corvofeng/go-flarum/model/flarum"
	"github.com/corvofeng/go-flarum/util"

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

var DefaultPreference = flarum.Preferences{
	Locale: "en",
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
		if user, err = SQLUserGetByName(gormDB, _userID); err == nil {
			break
		}

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

func SQLUserRegister(gormDB *gorm.DB, name, email, password string) (User, error) {
	user := User{
		Name:        name,
		Email:       email,
		Password:    password,
		Reputation:  20,
		Avatar:      fmt.Sprintf("https://robohash.org/%s", name),
		Preferences: []byte(`{}`),
	}

	result := gormDB.Create(&user)
	if result.Error != nil {
		if result.Error.Error() == fmt.Sprintf("Error 1062: Duplicate entry '%s' for key 'idx_name'", name) {
			return User{}, fmt.Errorf("用户名已经存在")
		}
		return User{}, result.Error
	}
	user.SetPreference(gormDB, DefaultPreference)
	return user, nil
}

// SQLGithubSync github用户同步信息
func (user *User) SQLGithubSync(gormDB *gorm.DB, gu *github.User) {
	logger := util.GetLogger()
	if user.Email != gu.GetEmail() {
		logger.Errorf("Wanna modify %s, but give %s", user.Email, gu.GetEmail())
		return
	}
	logger.Debugf("Sync user %s(%d) with github %+v", user.Name, user.ID, gu)
	if user.About == "" {
		gormDB.Model(user).Updates(User{Description: gu.GetBlog()})
	}
	if user.URL == "" {
		gormDB.Model(user).Updates(User{WebSite: gu.GetBlog()})
	}
	if user.Avatar == "" {
		gormDB.Model(user).Updates(User{Avatar: gu.GetAvatarURL()})
	}
	if user.Preferences == nil {
		user.Preferences = []byte(`{}`)
	}
	user.SetPreference(gormDB, DefaultPreference)
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
// func (user *User) SaveAvatar(redisDB *redis.Client, avatar string) {
// 	logger := util.GetLogger()

// 	if user == nil {
// 		return
// 	}
// 	_, err := sqlDB.Exec("UPDATE user SET avatar = ? WHERE id = ?", avatar, user.ID)
// 	if err != nil {
// 		logger.Error("Set ", user, " avatar ", avatar, " failed!!")
// 		return
// 	}
// 	redisDB.HSet("avatar", user.toKey(), avatar)
// 	logger.Notice("Refresh user avatar", user)
// }

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
func (user *User) SetPreference(gormDB *gorm.DB, preference flarum.Preferences) {
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
