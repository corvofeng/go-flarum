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

import "time"

// BaseDiscussion 基础类
type BaseDiscussion struct {
	Type string `json:"type"`

	ID    uint64 `json:"id"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

// Discussion 帖子或是讨论
type Discussion struct {
	BaseDiscussion

	CommentCount     uint64    `json:"commentCount"`
	ParticipantCount int       `json:"participantCount"`
	CreatedAt        time.Time `json:"createdAt"`

	LastPostedAt   string `json:"lastPostedAt"`
	LastPostNumber int    `json:"lastPostNumber"`

	CanReply  bool `json:"canReply"`
	CanRename bool `json:"canRename"`
	CanDelete bool `json:"canDelete"`
	CanHide   bool `json:"canHide"`

	IsHidden bool   `json:"isHidden"`
	HiddenAt string `json:"hiddenAt"`

	LastReadAt         string `json:"lastReadAt"`
	LastReadPostNumber int    `json:"lastReadPostNumber"`
}

// DiscussionRelations 帖子具有的关系
type DiscussionRelations struct {
	User           RelationDict `json:"user"` // 创建帖子的用户
	LastPostedUser RelationDict `json:"lastPostedUser"`
	FirstPost      RelationDict `json:"firstPost"`

	Tags           RelationArray `json:"tags"`
	LatestViews    RelationArray `json:"latestViews"`
	RecipientUsers RelationArray `json:"recipientUsers"`

	OldRecipientUsers  RelationArray `json:"oldRecipientUsers"`
	RecipientGroups    RelationArray `json:"recipientGroups"`
	OldRecipientGroups RelationArray `json:"oldRecipientGroups"`
}

// DoInit 初始化一篇帖子
func (d *BaseDiscussion) DoInit() {
	d.Type = "discussions"
}

// GetDefaultAttributes 获取属性
func (d *BaseDiscussion) GetDefaultAttributes(obj interface{}) {
	// uObj := obj.(model.User)
	// fmt.Println(uObj)
}

// GetType 获取类型
func (d *BaseDiscussion) GetType() string {
	return d.Type
}

// GetID 获取ID信息
func (d *BaseDiscussion) GetID() uint64 {
	return d.ID
}
