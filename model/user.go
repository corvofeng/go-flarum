package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"../util"

	"github.com/ego008/youdb"
)

type User struct {
	Id            uint64 `json:"id"`
	Name          string `json:"name"`
	Gender        string `json:"gender"`
	Flag          int    `json:"flag"`
	Avatar        string `json:"avatar"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	Url           string `json:"url"`
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
	Id     uint64 `json:"id"`
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

func UserGetById(db *youdb.DB, uid uint64) (User, error) {
	obj := User{}
	rs := db.Hget("user", youdb.I2b(uid))
	if rs.State == "ok" {
		json.Unmarshal(rs.Data[0], &obj)
		return obj, nil
	}
	return obj, errors.New(rs.State)
}

func UserUpdate(db *youdb.DB, obj User) error {
	jb, _ := json.Marshal(obj)
	return db.Hset("user", youdb.I2b(obj.Id), jb)
}

func UserGetByName(db *youdb.DB, name string) (User, error) {
	obj := User{}
	rs := db.Hget("user_name2uid", []byte(name))
	if rs.State == "ok" {
		rs2 := db.Hget("user", rs.Data[0])
		if rs2.State == "ok" {
			json.Unmarshal(rs2.Data[0], &obj)
			return obj, nil
		}
		return obj, errors.New(rs2.State)
	}
	return obj, errors.New(rs.State)
}

func UserGetIdByName(db *youdb.DB, name string) string {
	rs := db.Hget("user_name2uid", []byte(name))
	if rs.State == "ok" {
		return youdb.B2ds(rs.Data[0])
	}
	return ""
}

func UserListByFlag(db *youdb.DB, cmd, tb, key string, limit int) UserPageInfo {
	var items []User
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
		rs := db.Hmget("user", keys)
		if rs.State == "ok" {
			for i := 0; i < (len(rs.Data) - 1); i += 2 {
				item := User{}
				json.Unmarshal(rs.Data[i+1], &item)
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
	}

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

	rows, err := db.Query(
		"SELECT id, name, password, reputation, email, avatar, website, token, created_at FROM user WHERE id =  ?",
		uid)

	defer func() {
		if rows != nil {
			rows.Close() //可以关闭掉未scan连接一直占用
		}
	}()
	if err != nil {
		fmt.Printf("Query failed,err:%v", err)
		return obj, err
	}
	for rows.Next() {
		err = rows.Scan(
			&obj.Id,
			&obj.Name,
			&obj.Password,
			&obj.Reputation,
			&obj.Email,
			&obj.Avatar,
			&obj.Url,
			&obj.Token,
			&obj.RegTime,
		)

		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			return obj, errors.New("No result")
		}
	}

	return obj, nil
}

// SQLUserGetByName 获取数据库中用户
func SQLUserGetByName(db *sql.DB, name string) (User, error) {
	obj := User{}

	rows, err := db.Query(
		"SELECT id, name, password, reputation, email, avatar, website, token, created_at FROM user WHERE name =  ?",
		name)

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
			&obj.Name,
			&obj.Password,
			&obj.Reputation,
			&obj.Email,
			&obj.Avatar,
			&obj.Url,
			&obj.Token,
			&obj.RegTime,
		)

		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			return obj, errors.New("No result")
		}
	}

	return obj, nil
}

// SQLRegister 用户注册
func (user *User) SQLRegister(db *sql.DB) bool {
	row, err := db.Exec(
		("INSERT INTO `user` " +
			" (`name`, `email`, `urlname`, `password`, `reputation`)" +
			" VALUES " +
			" (?, ?,?, ?, ?)"),
		user.Name,
		user.Name,
		user.Name,
		user.Password,
		20, // 初始声望值20
	)
	if util.CheckError(err, "用户注册") {
		return false
	}
	uid, err := row.LastInsertId()
	user.Id = uint64(uid)

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

// CanEdit 检查当前用户是否可以创建帖子
func (user *User) CanEdit() bool {
	return user.IsAdmin()
}

// SaveAvatar 更新用户头像
func (user *User) SaveAvatar(db *sql.DB, cntDB *youdb.DB, avatar string) {
	logger := util.GetLogger()

	if user == nil {
		return
	}
	_, err := db.Exec("UPDATE user SET avatar = ? WHERE id = ?", avatar, user.Id)
	if err != nil {
		logger.Error("Set ", user, " avatar ", avatar, " failed!!")
		return
	}

	cntDB.Hdel("avatar", youdb.I2b(user.Id))
	logger.Notice("Refresh user avatar", user)
	return
}

// GetAvatarByID 获取用户头像
func GetAvatarByID(db *sql.DB, cntDB *youdb.DB, uid uint64) string {
	var avatar string
	logger := util.GetLogger()

	rs := cntDB.Hget("avatar", youdb.I2b(uid))
	if rs.State == "ok" {
		return rs.String()
	}

	user, err := SQLUserGetByID(db, uid)
	if util.CheckError(err, "查询用户") {
		return avatar
	}
	avatar = user.Avatar

	cntDB.Hset("avatar", youdb.I2b(uid), []byte(avatar))
	logger.Debug("key not found for ", user.Id, user.Name, "but we refresh!")
	return avatar
}
