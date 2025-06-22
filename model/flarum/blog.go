package flarum

// "github.com/corvofeng/go-flarum/model/flarum"

type FlarumBlogMeta struct {
	BaseResources

	// featuredImage
	// summary
	FeaturedImage string `json:"featuredImage,omitempty"`
	Summary       string `json:"summary,omitempty"`

	// isFeatured
	// isPendingReview
	// isSized
	IsFeatured      bool `json:"isFeatured,omitempty"`      // 是否是特色文章
	IsPendingReview bool `json:"isPendingReview,omitempty"` // 是否待审核
	IsSized         bool `json:"isSized,omitempty"`         // 是否已调整大小
}

type BlogMetaRelations struct {
	LastPostedDiscussion RelationDict   `json:"lastPostedDiscussion"`
	Children             []RelationDict `json:"children"`
}

// DoInit 初始化tags
func (t *FlarumBlogMeta) DoInit(id uint64) {
	t.setID(id)
	t.setType("blogMeta")
}

// GetType 获取类型
func (t *FlarumBlogMeta) GetType() string {
	return t.Type
}
