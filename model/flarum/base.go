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
	EForum EResourceType = "forums"

	// ETAG 标签信息
	ETAG EResourceType = "tag"

	// EPost 评论信息
	EPost EResourceType = "post"
)

// IDataBase flarum数据
type IDataBase interface {
	// GetAttributes()
	DoInit()
	GetType() string
	// GetID() uint64
	GetDefaultAttributes(obj interface{})
}

// BaseRelation flarum中的基础资源关系
type BaseRelation struct {
	Type string `json:"type"`
	ID   uint64 `json:"id"`
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
	ID            uint64      `json:"id"`
	Type          string      `json:"type"`
	Attributes    *IDataBase  `json:"attributes"`
	Relationships interface{} `json:"relationships"`
}

// Session flarum session数据
type Session struct {
	UserID    int    `json:"userId"`
	CsrfToken string `json:"csrfToken"`
}

// APIDoc flarum api将会返回的结果
type APIDoc struct {
	/**
	 * Links 当前可点的链接:
	 * 		first: 首页
	 * 		next: 下一页
	 * 		prev: 前一页
	 */
	Links map[string]string `json:"links"`

	Data     []Resource `json:"data"`
	Included []Resource `json:"included"`
}

// CoreData flarum需要返回的数据
type CoreData struct {
	Resources   []Resource        `json:"resources"`
	Sessions    Session           `json:"session"`
	Locale      string            `json:"locale"`
	Locales     map[string]string `json:"locales"`
	APIDocument APIDoc            `json:"apiDocument"`
}

// NewResource 根据类型初始化一个资源
func NewResource(resourceType EResourceType) Resource {
	var obj Resource
	var data IDataBase
	switch resourceType {
	case EBaseUser:
		data = &BaseUser{}
		break
	case ECurrentUser:
		data = &CurrentUser{}
		break
	case EDiscussion:
		data = &Discussion{}
		break
	case EForum:
		data = &Forum{}
		break
	case ETAG:
		data = &Tag{}
		break
	case EPost:
		data = &Post{}
		break
	}
	data.DoInit()
	obj = Resource{
		Type:       data.GetType(),
		Attributes: &data,
	}
	return obj
}

// NewAPIDoc 新建一个APIDoc对象
func NewAPIDoc() APIDoc {
	apiDoc := APIDoc{}
	apiDoc.Links = make(map[string]string)
	return apiDoc
}
