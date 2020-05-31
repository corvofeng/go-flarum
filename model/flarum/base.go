package flarum

// EResourceType flarum中的资源类型
// E是enum的意思
type EResourceType string

const (
	// EBaseUser 基础用户
	EBaseUser EResourceType = "baseuser"

	// ECurrentUser 当前用户
	ECurrentUser EResourceType = "current"

	// EBaseDiscussion 帖子基础资源
	EBaseDiscussion EResourceType = "base_discussion"

	// EDiscussion 帖子资源
	EDiscussion EResourceType = "discussion"

	// EForum 论坛信息
	EForum EResourceType = "forum"
)

// IDataBase flarum数据
type IDataBase interface {
	// GetAttributes()
	DoInit()
	GetDefaultAttributes(obj interface{})
}

// BaseRelation flarum中的基础资源关系
type BaseRelation struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// RelationDict 字典形式的关系
type RelationDict struct {
	Data BaseRelation `json:"data"`
}

// RelationArray 数组形式的关系
type RelationArray struct {
	Data []BaseRelation `json:"data"`
}

// Resource flarum资源
type Resource struct {
	ID            int           `json:"id"`
	Type          string        `json:"type"`
	Attributes    IDataBase     `json:"attributes"`
	Relationships []interface{} `json:"relationships"`
}

// Session flarum session数据
type Session struct {
	UserID    int    `json:"userId"`
	CsrfToken string `json:"csrfToken"`
}

// APIDoc flarum api document
type APIDoc struct {
	Links    []string    `json:"links"`
	Data     []IDataBase `json:"data"`
	Included []IDataBase `json:"included"`
}

// CoreData flarum需要返回的数据
type CoreData struct {
	Resources   []IDataBase         `json:"resource"`
	Sessions    Session             `json:"session"`
	Locale      string              `json:"locale"`
	Locales     []map[string]string `json:"locales"`
	APIDocument APIDoc              `json:"apiDocument"`
}

// NewResource 根据类型初始化一个资源
func NewResource(resourceType EResourceType) IDataBase {

	var obj IDataBase
	switch resourceType {
	case EBaseUser:
		obj = &BaseUser{}
		break
	case ECurrentUser:
		obj = &CurrentUser{}
		break

	case EForum:
		obj = &Forum{}
		break
	}
	obj.DoInit()
	return obj
}
