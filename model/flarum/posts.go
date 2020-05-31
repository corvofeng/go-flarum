package flarum

import "time"

// Post flarum 评论信息
type Post struct {
	Type string `json:"type"`
	ID   uint64 `json:"id"`

	Number      uint64    `json:"number"`
	CreatedAt   time.Time `json:"createdAt"`
	ContentType string    `json:"contentType"`
	ContentHTML string    `json:"contentHtml"`
	Content     string    `json:"content"`
	IPAddress   string    `json:"ipAddress"`
	CanEdit     bool      `json:"canEdit"`
	CanDelete   bool      `json:"canDelete"`
	CanHide     bool      `json:"canHide"`
	IsApproved  bool      `json:"isApproved"`
	CanApprove  bool      `json:"canApprove"`
	CanFlag     bool      `json:"canFlag"`
	CanLike     bool      `json:"canLike"`
}

// PostRelations 评论具有的关系
type PostRelations struct {
	User        RelationDict  `json:"user"`
	Discussion  RelationDict  `json:"discussion"`
	Flags       RelationArray `json:"flags"`
	Likes       RelationArray `json:"likes"`
	MentionedBy RelationArray `json:"mentionedBy"`
}

// DoInit 初始化评论数据
func (p *Post) DoInit() {
	p.Type = "posts"
}

// GetType 获取类型
func (p *Post) GetType() string {
	return p.Type
}

// GetDefaultAttributes 获取属性
func (p *Post) GetDefaultAttributes(obj interface{}) {
	// uObj := obj.(model.User)
	// fmt.Println(uObj)
}
