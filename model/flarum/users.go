package flarum

import (
	"time"
)

// BaseUser 基础的用户类
type BaseUser struct {
	Type        string `json:"type"`
	ID          uint64 `json:"id"`
	Username    string `json:"username"`
	Displayname string `json:"displayName"`
	AvatarURL   string `json:"avatarUrl"`
}

// CurrentUser 当前用户信息
type CurrentUser struct {
	BaseUser
	IsEmailConfirmed bool   `json:"isEmailConfirmed"`
	Email            string `json:"email"`

	MarkedAllAsReadAt       time.Time `json:"markedAllAsReadAt"`
	UnreadNotificationCount int       `json:"unreadNotificationCount"`
	NewNotificationCount    int       `json:"newNotificationCount"`
	Preferences             []string  `json:"preferences"`
}

// DoInit 初始化用户类
func (u *BaseUser) DoInit() {
	u.Type = "users"
}

// GetType 获取类型
func (u *BaseUser) GetType() string {
	return u.Type
}

// GetID 获取ID信息
func (u *BaseUser) GetID() uint64 {
	return u.ID
}

// GetDefaultAttributes 获取属性
func (u *BaseUser) GetDefaultAttributes(obj interface{}) {
	// uObj := obj.(model.User)
	// fmt.Println(uObj)
}
