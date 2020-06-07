package flarum

import "time"

// Tag flarum tag信息
type Tag struct {
	BaseResources

	Name               string      `json:"name"`
	Description        string      `json:"description"`
	Slug               string      `json:"slug"`
	Color              string      `json:"color"`
	BackgroundURL      string      `json:"backgroundUrl"`
	BackgroundMode     string      `json:"backgroundMode"`
	Icon               string      `json:"icon"`
	DiscussionCount    uint64      `json:"discussionCount"`
	Position           int         `json:"position"`
	DefaultSort        interface{} `json:"defaultSort"`
	IsChild            bool        `json:"isChild"`
	IsHidden           bool        `json:"isHidden"`
	LastPostedAt       time.Time   `json:"lastPostedAt"`
	CanStartDiscussion bool        `json:"canStartDiscussion"`
	CanAddToDiscussion bool        `json:"canAddToDiscussion"`
	IsRestricted       bool        `json:"isRestricted"`
}

// TagRelations 标签具有的关系
type TagRelations struct {
	LastPostedDiscussion RelationDict `json:"lastPostedDiscussion"`
	Parent               RelationDict `json:"parent"`
}

// DoInit 初始化tags
func (t *Tag) DoInit() {
	t.Type = "tags"
}

// GetType 获取类型
func (t *Tag) GetType() string {
	return t.Type
}

// // GetAttributes 获取属性
// func (t *Tag) GetAttributes() map[string]interface{} {
// 	return nil
// }
