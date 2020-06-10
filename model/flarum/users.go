package flarum

import (
	"time"
)

// BaseUser 基础的用户类
type BaseUser struct {
	BaseResources

	Username    string `json:"username"`
	Displayname string `json:"displayName"`
	AvatarURL   string `json:"avatarUrl"`

	JoinTime        string `json:"joinTime"`
	DiscussionCount int    `json:"discussionCount"`
	CommentCount    int    `json:"commentCount"`
	CanEdit         bool   `json:"canEdit"`
	CanDelete       bool   `json:"canDelete"`
	LastSeenAt      string `json:"lastSeenAt"`
	CanSuspend      bool   `json:"canSuspend"`
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

// UserRelations 用户所具有的关系
type UserRelations struct {
	Groups RelationArray `json:"groups"`
}

// DoInit 初始化用户类
func (u *BaseUser) DoInit(id uint64) {
	u.setID(id)
	u.setType("users")
}

// GetType 获取类型
func (u *BaseUser) GetType() string {
	return u.Type
}

// GetID 获取ID信息
func (u *BaseUser) GetID() uint64 {
	return u.id
}

// // GetAttributes 获取属性
// func (u *BaseUser) GetAttributes() map[string]interface{} {
// 	// uObj := obj.(model.User)
// 	// fmt.Println(uObj)
// 	return nil
// }
