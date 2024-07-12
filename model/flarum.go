package model

import (
	"encoding/json"
	"time"

	"github.com/corvofeng/go-flarum/model/flarum"
	"github.com/corvofeng/go-flarum/util"
)

// FlarumCreateForumInfo 从SiteInfo创建ForumInfo
// tags 当前站点具有的标签集合, TODO: 缓存
func FlarumCreateForumInfo(
	user *User,
	appConf AppConf,
	siteInfo SiteInfo,
	tags []flarum.Resource,
) flarum.Resource {
	obj := flarum.NewResource(flarum.EForum, 1) // flarum 默认为1
	data := (obj.Attributes).(*flarum.Forum)

	mainConf := appConf.Main
	siteConf := appConf.Site

	data.DefaultRoute = "/all"
	data.Title = siteConf.Name
	data.Description = siteConf.Desc
	data.APIURL = FlarumAPIPath
	if user.IsAdmin() {
		data.AdminURL = "/admin"
	}
	// data.BasePath = "http://192.168.101.35:8082"
	data.BaseURL = mainConf.BaseURL
	data.CanStartDiscussion = true
	data.AllowSignUp = siteConf.AllowSignup
	data.ShowLanguageSelector = true
	data.WelcomeMessage = siteConf.WelcomeMessage
	data.WelcomeTitle = siteConf.WelcomeTitle
	data.Debug = appConf.Main.Debug
	data.CanSearchUsers = true
	data.CanViewForum = true
	data.MaxPrimaryTags = 3
	data.MaxSecondaryTags = 3
	data.MinPrimaryTags = 1
	data.MinSecondaryTags = 0

	data.BasePath = ""
	// data.BaseURL = "/"
	obj.BindRelations(
		"Tags",
		FlarumCreateTagRelations(tags),
	)

	return obj
}

// FlarumCreateLocale 生成语言包配置
func FlarumCreateLocale(coreData *flarum.CoreData, locale string) {
	coreData.Locales = make(map[string]string)
	coreData.Locales["en"] = "English"
	coreData.Locales["zh"] = "中文"
	coreData.Locale = locale
}

// FlarumCreateTag 创建tag资源
func FlarumCreateTag(cat Tag) flarum.Resource {
	obj := flarum.NewResource(flarum.ETAG, cat.ID)

	data := obj.Attributes.(*flarum.Tag)
	data.Name = cat.Name
	data.DiscussionCount = 3
	data.IsHidden = cat.Hidden
	data.Slug = cat.URLName
	data.CanAddToDiscussion = true
	data.CanStartDiscussion = true
	data.LastPostedAt = "2020-06-10T01:20:37+00:00"
	data.Position = cat.Position

	data.IsChild = false
	data.IsRestricted = false

	data.Color = cat.Color
	data.Icon = cat.IconIMG

	// please refer to #11
	if cat.ParentID == 0 {
		obj.Relationships = flarum.TagRelations{
			LastPostedDiscussion: flarum.RelationDict{Data: flarum.InitBaseResources(1, "discussions")},
			// Children:             []flarum.RelationDict{},
		}
	} else {
		obj.Relationships = flarum.TagChildRelations{}
		obj.BindRelations(
			"Parent",
			flarum.RelationDict{Data: flarum.InitBaseResources(cat.ParentID, "tags")},
		)
	}

	return obj
}

// FlarumCreateDiscussion 创建帖子资源
func FlarumCreateDiscussion(topic Topic) flarum.Resource {
	obj := flarum.NewResource(flarum.EDiscussion, topic.ID)
	data := obj.Attributes.(*flarum.Discussion)
	data.Title = topic.Title
	data.CommentCount = topic.CommentCount
	data.CreatedAt = topic.CreatedAt.Format(time.RFC3339)
	data.FirstPostID = topic.FirstPostID
	data.CanReply = true
	data.CanHide = true
	data.Subscription = "follow"

	data.LastPostNumber = topic.LastPostID
	data.LastPostedAt = topic.UpdatedAt.Format(time.RFC3339)
	// data.LastPostID = article.LastPostID
	data.LastUserID = topic.LastPostUserID

	obj.BindRelations(
		"User",
		flarum.RelationDict{Data: flarum.InitBaseResources(topic.UserID, "users")},
	)
	var flarumTags []flarum.Resource
	for _, t := range topic.Tags {
		tag := FlarumCreateTag(t)
		flarumTags = append(flarumTags, tag)
	}

	obj.BindRelations("Tags",
		FlarumCreateTagRelations(flarumTags),
	)

	obj.BindRelations(
		"FirstPost",
		flarum.RelationDict{
			Data: flarum.InitBaseResources(data.FirstPostID, "posts"),
		},
	)

	obj.BindRelations(
		"Posts",
		flarum.RelationArray{
			Data: []flarum.BaseRelation{},
		},
	)

	if topic.LastPostID != 0 {
		// data.LastPostID = article.LastPostID
		data.LastUserID = topic.LastPostUserID
		obj.BindRelations(
			"LastPostedUser",
			flarum.RelationDict{
				Data: flarum.InitBaseResources(topic.LastPostUserID, "users"),
			},
		)
	}

	return obj
}

// FlarumCreateCurrentUser 创建用户资源
func FlarumCreateCurrentUser(user User) flarum.Resource {
	return FlarumCreateUser(user)
}

// FlarumCreateUser 创建用户资源
func FlarumCreateUser(user User) flarum.Resource {
	obj := flarum.NewResource(flarum.ECurrentUser, user.ID)
	data := obj.Attributes.(*flarum.CurrentUser)
	data.Username = user.Name
	data.Displayname = user.Name
	data.AvatarURL = user.Avatar
	data.Email = user.Email
	data.IsEmailConfirmed = true
	data.JoinTime = user.CreatedAt.String()
	data.Slug = user.Name
	if len(user.Preferences) > 0 {
		err := json.Unmarshal(user.Preferences, &data.Preferences)
		if err != nil {
			util.GetLogger().Errorf("Can't get preferences for user %s", user.ID, user.Preferences)
		}
	}

	obj.BindRelations(
		"Groups",
		flarum.RelationArray{
			Data: []flarum.BaseRelation{},
		},
	)
	return obj
}

// FlarumCreateGroup 创建组信息
func FlarumCreateGroup() flarum.Resource {
	obj := flarum.NewResource(flarum.EGroup, 1)
	data := obj.Attributes.(*flarum.Group)
	// 	color: "#B72A2A"
	// icon:
	// isHidden: 0
	// namePlural: "Admins"
	// nameSingular: "Admin"
	data.Color = "#B72A2A"
	data.Icon = "fas fa-wrench"
	data.NamePlural = "Admins"
	data.NameSingular = "Admin"
	return obj
}

// FlarumCreatePost 创建评论
func FlarumCreatePost(comment Comment, currentUser *User) flarum.Resource {
	obj := flarum.NewResource(flarum.EPost, comment.ID)
	data := obj.Attributes.(*flarum.Post)

	data.Number = comment.Number
	data.ContentType = "comment"
	data.Content = comment.Content
	data.ContentHTML = comment.ContentFmt
	data.CreatedAt = comment.CreatedAt.Format(time.RFC3339)

	if currentUser != nil {
		data.CanLike = true
	}

	if currentUser != nil && currentUser.IsAdmin() {
		data.CanEdit = true
		data.CanHide = true
		data.IsApproved = true
		data.IPAddress = comment.ClientIP
	}

	obj.BindRelations(
		"User",
		flarum.RelationDict{
			Data: flarum.InitBaseResources(comment.UID, "users"),
		},
	)
	obj.BindRelations(
		"Discussion",
		flarum.RelationDict{
			Data: flarum.InitBaseResources(comment.AID, "discussions"),
		},
	)
	obj.BindRelations(
		"Likes",
		FlarumCreateUserLikeRelations(comment.Likes),
	)
	obj.BindRelations(
		"Flags",
		flarum.RelationArray{
			Data: []flarum.BaseRelation{},
		},
	)
	obj.BindRelations(
		"MentionedBy",
		flarum.RelationArray{
			Data: []flarum.BaseRelation{},
		},
	)

	obj.BindRelations(
		"MentionsPosts",
		flarum.RelationArray{
			Data: []flarum.BaseRelation{},
		},
	)

	return obj
}

// FlarumCreateUserLikeRelations 点赞关系
func FlarumCreateUserLikeRelations(userList []uint64) flarum.RelationArray {
	userLikes := flarum.RelationArray{
		Data: []flarum.BaseRelation{},
	}

	for _, userID := range userList {
		userLikes.Data = append(userLikes.Data,
			flarum.InitBaseResources(userID, "users"),
		)
	}
	return userLikes
}

// FlarumCreatePostRelations 创建关系结构
func FlarumCreatePostRelations(postArr []flarum.Resource, comments []uint64) flarum.IRelation {
	var obj flarum.RelationArray
	for _, p := range postArr {
		obj.Data = append(
			obj.Data,
			flarum.InitBaseResources(p.GetID(), p.Type),
		)
	}

	p := flarum.Post{}
	p.DoInit(0)
	for _, cid := range comments {
		obj.Data = append(
			obj.Data,
			flarum.InitBaseResources(cid, p.Type),
		)
	}

	return obj
}

// FlarumCreateTagRelations 创建关系结构
func FlarumCreateTagRelations(tagArr []flarum.Resource) flarum.IRelation {
	var obj flarum.RelationArray
	for _, p := range tagArr {
		obj.Data = append(
			obj.Data,
			flarum.InitBaseResources(p.GetID(), p.Type),
		)
	}

	return obj
}
