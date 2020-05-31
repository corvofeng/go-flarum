package flarum

/**
 * 与topic行为一致
 *	refer to:
 *		view/flarum/src/Api/Serializer/DiscussionSerializer.php
 */

import "time"

// BaseDiscussion 基础类
type BaseDiscussion struct {
	Type string `json:"type"`

	Title string `json:"title"`
	Slug  string `json:"slug"`
}

// Discussion 帖子或是讨论
type Discussion struct {
	BaseDiscussion

	CommentCount     int       `json:"commentCount "`
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

// Init 初始化一篇帖子
func (d *BaseDiscussion) Init() {
	d.Type = "discussions"
}
