package flarum

import (
	"encoding/json"
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
	DoInit(uint64)
	GetType() string
	// GetID() uint64
	GetAttributes() (map[string]interface{}, error)
}

// BaseRelation flarum中的基础资源关系
type BaseRelation struct {
	BaseResources
}

// Struct2Map 将结构体转换为json
/**
 * From https://stackoverflow.com/a/42849112 也许这样的方式并不快, 但一定是bug最少的
 * 如果成为了瓶颈再考虑优化吧
 */
func Struct2Map(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj) // Convert to a json string
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &newMap) // Convert to a map
	return
}

// -------------  BaseResources ---------------

// BaseResources 基础的资源结构
/**
 * 请不要直接修改结构体中的变量
 */
type BaseResources struct {
	// Issue-5: flarum needs id as string
	id uint64

	ID string `json:"id"`

	Type string `json:"type"`
}

// DoInit 空函数, 占位使用
func (r *BaseResources) DoInit() {}

// setID 绑定ID
func (r *BaseResources) setID(id uint64) {
	r.id = id
	r.ID = strconv.FormatUint(id, 10)
}

// SetType 绑定类型
func (r *BaseResources) setType(t string) {
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

// GetAttributes 获取结构体的属性值
/**
 * 基类将会默认拥有这一函数, 但是理论上讲该函数不该被调用, 就和base resources不该被使用一样
 */
func (r *BaseResources) GetAttributes() (map[string]interface{}, error) {
	panic("Please write your own get attributes")
}

// -------------  BaseResources ---------------

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

// GetAttributes 获取结构体的属性值, 基类将会继承这一函数
func (r *Resource) GetAttributes() (map[string]interface{}, error) {
	return Struct2Map(r)
}

// Session flarum session数据
type Session struct {
	UserID    uint64 `json:"userId"`
	CsrfToken string `json:"csrfToken"`
}

// APIDoc flarum api将会返回的结果
type APIDoc struct {
	/**
	 * 虽然感觉没有在用, 但是需要保留
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

// CoreData flarum页面需要返回的数据
type CoreData struct {
	Resources   []Resource        `json:"resources"`
	Sessions    Session           `json:"session"`
	Locale      string            `json:"locale"`
	Locales     map[string]string `json:"locales"`
	APIDocument APIDoc            `json:"apiDocument"`
}

// NewResource 根据类型初始化一个资源
func NewResource(resourceType EResourceType, id uint64) Resource {
	var obj Resource
	var data IDataBase
	var defaultRelation IRelation

	switch resourceType {
	// golang no need break
	case EBaseUser:
		data = &BaseUser{}
		defaultRelation = &UserRelations{}
	case ECurrentUser:
		data = &CurrentUser{}
		defaultRelation = &UserRelations{}
	case EDiscussion:
		data = &Discussion{}
		defaultRelation = &DiscussionRelations{}
	case EForum:
		data = &Forum{}
		defaultRelation = &ForumRelations{}
	case ETAG:
		data = &Tag{}
		defaultRelation = &TagRelations{}
	case EPost:
		data = &Post{}
		defaultRelation = &PostRelations{}
	case EGroup:
		data = &Group{}
	}
	data.DoInit(id)
	obj = Resource{
		Attributes:    data,
		Relationships: defaultRelation,
	}
	obj.setID(id)
	obj.setType(data.GetType())
	return obj
}

// newAPIDoc 新建一个APIDoc对象
func newAPIDoc() APIDoc {
	apiDoc := APIDoc{}
	apiDoc.Links = make(map[string]string)
	apiDoc.Data = []Resource{}
	apiDoc.Included = []Resource{}
	return apiDoc
}

// NewCoreData 新建一个CoreData对象
// 使用方法:
// 	coreData := flarum.NewCoreData()
// 	apiDoc := &coreData.APIDocument // 注意, 获取到的是指针
func NewCoreData() CoreData {
	coreData := CoreData{}
	coreData.APIDocument = newAPIDoc()
	return coreData
}

// NewAdminCoreData 新建一个CoreData对象
// 使用方法:
// 	coreData := flarum.NewAdminCoreData()
// 	apiDoc := &coreData.APIDocument // 注意, 获取到的是指针
func NewAdminCoreData() AdminCoreData {
	adminCoreData := AdminCoreData{}
	adminCoreData.APIDocument = newAPIDoc()
	return adminCoreData
}

// SetData 设置为字典类型的数据
func (apiDoc *APIDoc) SetData(data interface{}) {
	apiDoc.Data = data
}

// AppendResources 添加资源
func (apiDoc *APIDoc) AppendResources(res Resource) {
	apiDoc.Included = append(apiDoc.Included, res)
}

// AppendResources 添加资源
func (coreData *CoreData) AppendResources(res Resource) {
	coreData.APIDocument.AppendResources(res)
	coreData.Resources = append(coreData.Resources, res)
}

// AddCurrentUser 增加当前用户的信息
func (coreData *CoreData) AddCurrentUser(user Resource) {
	coreData.AppendResources(user)
}

// AddSessionData 添加用户的session信息, 仅用于csrf
func (coreData *CoreData) AddSessionData(user Resource, csrf string) {
	coreData.Sessions = Session{
		UserID:    user.GetID(),
		CsrfToken: csrf,
	}
}

// BindRelations 绑定关系
func (r *Resource) BindRelations(field string, data IRelation) {
	reflect.ValueOf(r.Relationships).Elem().FieldByName(field).Set(reflect.ValueOf(data))
}

// InitBaseResources 初始化一个基础资源
func InitBaseResources(id uint64, t string) BaseRelation {
	br := BaseRelation{}
	br.setID(id)
	br.setType(t)
	return br
}
