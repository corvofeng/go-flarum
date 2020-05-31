package flarum

/**
 * refer to:
 *   view/flarum/src/Api/Serializer/ForumSerializer.php
 */

// Forum 论坛的基础信息
type Forum struct {
	Type string `json:"type"`

	Title                string      `json:"title"`
	Description          string      `json:"description"`
	ShowLanguageSelector bool        `json:"showLanguageSelector"`
	BaseURL              string      `json:"baseUrl"`
	BasePath             string      `json:"basePath"`
	Debug                bool        `json:"debug"`
	APIURL               string      `json:"apiUrl"`
	WelcomeTitle         string      `json:"welcomeTitle"`
	WelcomeMessage       string      `json:"welcomeMessage"`
	ThemePrimaryColor    string      `json:"themePrimaryColor"`
	ThemeSecondaryColor  string      `json:"themeSecondaryColor"`
	LogoURL              string      `json:"logoUrl"`
	FaviconURL           string      `json:"faviconUrl"`
	HeaderHTML           string      `json:"headerHtml"`
	FooterHTML           string      `json:"footerHtml"`
	AllowSignUp          bool        `json:"allowSignUp"`
	DefaultRoute         string      `json:"defaultRoute"`
	CanViewDiscussions   bool        `json:"canViewDiscussions"`
	CanStartDiscussion   bool        `json:"canStartDiscussion"`
	CanViewUserList      bool        `json:"canViewUserList"`
	AdminURL             string      `json:"adminUrl"`
	Version              string      `json:"version"`
	CanViewFlags         bool        `json:"canViewFlags"`
	FlagCount            int         `json:"flagCount"`
	GuidelinesURL        interface{} `json:"guidelinesUrl"`
	MinPrimaryTags       string      `json:"minPrimaryTags"`
	MaxPrimaryTags       string      `json:"maxPrimaryTags"`
	MinSecondaryTags     string      `json:"minSecondaryTags"`
	MaxSecondaryTags     string      `json:"maxSecondaryTags"`
}

// ForumRelations 站点关系
type ForumRelations struct {
	Groups    RelationArray `json:"groups"`
	Tags      RelationArray `json:"tags"`
	Links     RelationArray `json:"links"`
	Reactions RelationArray `json:"reactions"`
}

// DoInit 初始化forum
func (f *Forum) DoInit() {
	f.Type = "forum"
}

// GetDefaultAttributes 获取属性
func (f Forum) GetDefaultAttributes(obj interface{}) {
}
