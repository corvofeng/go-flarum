package flarum

/**
 * 与topic行为一致
 *	refer to:
 *		view/flarum/src/Api/Serializer/DiscussionSerializer.php
 *
 * Flarum 中为什么称这个变量为Discussion, 这是根据数据库的内容定义来的:
 *   数据库中:
 *      Discussion 为一个议题
 * 		Post 为议题下放的评论
 * 	用户创建时, 可以
 */

// BaseDiscussion 基础类
type BaseDiscussion struct {
	BaseResources

	Title string `json:"title"`
	Slug  string `json:"slug"`
}

// Discussion 帖子或是讨论
// view/flarum/migrations/2015_02_24_000000_create_discussions_table.php
type Discussion struct {
	BaseDiscussion

	CommentCount     uint64 `json:"commentCount"`
	ParticipantCount int    `json:"participantCount"`
	LastPostNumber   int    `json:"lastPostNumber"`

	// 第一个评论的信息, 通常由作者创建
	CreatedAt   string `json:"createdAt"`
	StartPostID string
	StartUserID string

	// 最后一次评论的信息
	LastPostedAt string `json:"lastPostedAt"`
	LastPostedID string
	LastUserID   string

	CanReply bool `json:"canReply"`
	// CanRename bool `json:"canRename"`
	// CanDelete bool `json:"canDelete"`
	// CanHide   bool `json:"canHide"`
	// CanLock   bool `json:"canLock"`

	// IsHidden   bool `json:"isHidden"`
	// IsApproved bool `json:"isApproved"`
	// IsLocked   bool `json:"isLocked"`
	// IsSticky   bool `json:"isSticky"`

	// HiddenAt   string `json:"hiddenAt"`
	// LastReadAt string `json:"lastReadAt"`
	// Subscription bool `json:"subscription"`
	// LastReadPostNumber int    `json:"lastReadPostNumber"`
}

// DiscussionRelations 帖子具有的关系
type DiscussionRelations struct {
	User RelationDict `json:"user"` // 创建帖子的用户
	// LastPostedUser RelationDict `json:"lastPostedUser"`
	// FirstPost      RelationDict `json:"firstPost"`

	Tags  RelationArray `json:"tags"`
	Posts RelationArray `json:"posts"`
	// LatestViews    RelationArray `json:"latestViews"`
	// RecipientUsers RelationArray `json:"recipientUsers"`

	// OldRecipientUsers  RelationArray `json:"oldRecipientUsers"`
	// RecipientGroups    RelationArray `json:"recipientGroups"`
	// OldRecipientGroups RelationArray `json:"oldRecipientGroups"`
}

// DoInit 初始化一篇帖子
func (d *BaseDiscussion) DoInit() {
	d.SetType("discussions")
}

// GetType 获取类型
func (d *BaseDiscussion) GetType() string {
	return d.Type
}

// GetID 获取ID信息
func (d *BaseDiscussion) GetID() uint64 {
	return d.id
}

// // GetAttributes 获取属性
// func (d *BaseDiscussion) GetAttributes() map[string]interface{} {
// 	// uObj := obj.(model.User)
// 	// fmt.Println(uObj)
// 	return nil
// }
