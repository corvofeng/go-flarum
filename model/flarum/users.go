package flarum

import (
	"fmt"
	"goyoubbs/model"
	"time"
)

// BaseUser 基础的用户类
type BaseUser struct {
	Type        string `json:"type"`
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
	u.Type = "user"
}

// GetDefaultAttributes 获取属性
func (u BaseUser) GetDefaultAttributes(obj interface{}) {
	uObj := obj.(model.User)
	fmt.Println(uObj)
}
