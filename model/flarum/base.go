package flarum

import (
	"reflect"
	"strconv"
)

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

	// EGroup 信息
	EGroup EResourceType = "group"
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
	BaseResources
}

// BaseResources 基础的资源结构
type BaseResources struct {
	// Issue-5: flarum needs id as string
	id uint64
	ID string `json:"id"`

	Type string `json:"type"`
}

// IRelation 具有的一些函数
type IRelation interface {
	// field, data
	// BindRelation(string, interface{})
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
	BaseResources
	Attributes    IDataBase `json:"attributes"`
	Relationships IRelation `json:"relationships"`
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

	/**
	 * Data API返回是的主要数据, 有一点很坑:
	 *    在disscussion帖子信息时 此变量是数组类型, 是要展示的主题合集
	 *    但是在请求post评论信息是 此变量是字典类型, 当前的评论所对应的一个主题
	 *
	 *  ALERT: 这里必须使用interface{}类型, 并且赋值时只能使用SetData函数
	 */
	Data interface{} `json:"data"`

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

	var defaultRelation IRelation

	switch resourceType {
	case EBaseUser:
		data = &BaseUser{}
		defaultRelation = &UserRelations{}
		break
	case ECurrentUser:
		data = &CurrentUser{}
		defaultRelation = &UserRelations{}
		break
	case EDiscussion:
		data = &Discussion{}
		defaultRelation = &DiscussionRelations{}
		break
	case EForum:
		data = &Forum{}
		break
	case ETAG:
		data = &Tag{}
		break
	case EPost:
		data = &Post{}
		defaultRelation = &PostRelations{}
		break
	case EGroup:
		data = &Group{}
		break
	}
	data.DoInit()
	obj = Resource{
		Attributes:    data,
		Relationships: defaultRelation,
	}
	obj.BaseResources.SetType(data.GetType())
	return obj
}

// NewAPIDoc 新建一个APIDoc对象
func NewAPIDoc() APIDoc {
	apiDoc := APIDoc{}
	apiDoc.Links = make(map[string]string)
	return apiDoc
}

// SetData 设置为字典类型的数据
func (apiDoc *APIDoc) SetData(data interface{}) {
	apiDoc.Data = data
}

// BindRelations 绑定关系
func (r *Resource) BindRelations(field string, data IRelation) {
	reflect.ValueOf(r.Relationships).Elem().FieldByName(field).Set(reflect.ValueOf(data))
}

// SetID 绑定ID
func (r *BaseResources) SetID(id uint64) {
	r.id = id
	r.ID = strconv.FormatUint(id, 10)
}

// SetType 绑定类型
func (r *BaseResources) SetType(t string) {
	r.Type = t
}

// GetID 获取ID
func (r *BaseResources) GetID() uint64 {
	return r.id
}

// GetType 绑定类型
func (r *BaseResources) GetType() string {
	return r.Type
}

// InitBaseResources 初始化一个基础资源
func InitBaseResources(id uint64, t string) BaseRelation {
	br := BaseRelation{}
	br.SetID(id)
	br.SetType(t)
	return br
}
