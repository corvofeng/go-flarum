package flarum

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
	Position           uint64      `json:"position"`
	DefaultSort        interface{} `json:"defaultSort"`
	IsChild            bool        `json:"isChild"`
	IsHidden           bool        `json:"isHidden"`
	LastPostedAt       string      `json:"lastPostedAt"`
	CanStartDiscussion bool        `json:"canStartDiscussion"`
	CanAddToDiscussion bool        `json:"canAddToDiscussion"`
	IsRestricted       bool        `json:"isRestricted"`
}

// TagChildRelations 标签具有的关系
// 子节点需要携带父节点的信息
type TagChildRelations struct {
	LastPostedDiscussion RelationDict `json:"lastPostedDiscussion"`
	Parent               RelationDict `json:"parent"`
}

// TagRelations 标签具有的关系
type TagRelations struct {
	LastPostedDiscussion RelationDict   `json:"lastPostedDiscussion"`
	Children             []RelationDict `json:"children"`
}

// DoInit 初始化tags
func (t *Tag) DoInit(id uint64) {
	t.setID(id)
	t.setType("tags")
}

// GetType 获取类型
func (t *Tag) GetType() string {
	return t.Type
}

// // GetAttributes 获取属性
// func (t *Tag) GetAttributes() map[string]interface{} {
// 	return nil
// }
