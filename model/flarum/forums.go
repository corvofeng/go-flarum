package flarum

/**
 * refer to:
 *   view/flarum/src/Api/Serializer/ForumSerializer.php
 */

// Forum 论坛的基础信息
type Forum struct {
	BaseResources

	Title                string      `json:"title"`
	Description          string      `json:"description"`
	ShowLanguageSelector bool        `json:"showLanguageSelector"`
	BaseURL              string      `json:"baseUrl"`
	JsChunksBaseUrl      string      `json:"jsChunksBaseUrl"`
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
	CanSearchUsers       bool        `json:"canSearchUsers"` // Please refer to /view/extensions/forum.js
	CanViewForum         bool        `json:"canViewForum"`
	AdminURL             string      `json:"adminUrl"`
	Version              string      `json:"version"`
	CanViewFlags         bool        `json:"canViewFlags"`
	FlagCount            int         `json:"flagCount"`
	GuidelinesURL        interface{} `json:"guidelinesUrl"`
	MinPrimaryTags       int         `json:"minPrimaryTags"`
	MaxPrimaryTags       int         `json:"maxPrimaryTags"`
	MinSecondaryTags     int         `json:"minSecondaryTags"`
	MaxSecondaryTags     int         `json:"maxSecondaryTags"`

	// (new Extend\Settings)
	//     ->serializeToForum('flarum-emoji.cdn', 'flarum-emoji.cdn')
	//     ->default('flarum-emoji.cdn', 'https://cdn.jsdelivr.net/gh/twitter/twemoji@[version]/assets/'),
	FlarumEmojiCDN string `json:"flarum-emoji.cdn,omitempty"` // Emoji CDN URL, optional
}

// ForumRelations 站点关系
type ForumRelations struct {
	Groups    RelationArray `json:"groups"`
	Tags      RelationArray `json:"tags"`
	Links     RelationArray `json:"links"`
	Reactions RelationArray `json:"reactions"`
}

// DoInit 初始化forum
func (f *Forum) DoInit(id uint64) {
	f.setID(id)
	f.Type = "forums"
}

// GetType 获取类型
func (f *Forum) GetType() string {
	return f.Type
}

// // GetAttributes 获取属性
// func (f *Forum) GetAttributes() map[string]interface{} {
// 	return nil
// }
