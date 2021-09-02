package model

import (
	"zoe/model/flarum"
	"zoe/util"
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
	data.AllowSignUp = true
	data.ShowLanguageSelector = true
	data.WelcomeMessage = siteConf.WelcomeMessage
	data.WelcomeTitle = siteConf.WelcomeTitle
	data.Debug = appConf.Main.Debug
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
func FlarumCreateTag(cat Category) flarum.Resource {
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
func FlarumCreateDiscussion(article ArticleListItem) flarum.Resource {
	lastComment := article.LastComment
	obj := flarum.NewResource(flarum.EDiscussion, article.ID)
	data := obj.Attributes.(*flarum.Discussion)
	data.Title = article.Title
	data.CommentCount = article.Comments
	data.CreatedAt = article.AddTimeFmt
	data.LastPostID = article.LastPostID
	data.FirstPostID = article.FirstPostID
	data.CanReply = true
	if lastComment != nil {
		data.LastPostNumber = lastComment.Number
		data.LastPostedAt = lastComment.AddTimeFmt
		data.LastUserID = lastComment.UID
	}

	obj.BindRelations(
		"User",
		flarum.RelationDict{Data: flarum.InitBaseResources(article.UID, "users")},
	)
	obj.BindRelations(
		"Tags",
		flarum.RelationArray{
			Data: []flarum.BaseRelation{},
		},
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

	if article.LastPostID != 0 {
		data.LastPostID = article.LastPostID
		data.LastUserID = lastComment.UID
		obj.BindRelations(
			"LastPostedUser",
			flarum.RelationDict{
				Data: flarum.InitBaseResources(lastComment.UID, "users"),
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
	data.Displayname = user.Nickname
	data.AvatarURL = user.Avatar
	data.Email = user.Email
	data.IsEmailConfirmed = true
	data.JoinTime = util.TimeFmt(user.RegTime, util.TIME_FMT, 0)
	data.Slug = user.Name
	if user.Preferences != nil {
		data.Preferences = *user.Preferences
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
func FlarumCreatePost(comment CommentListItem, currentUser *User) flarum.Resource {
	obj := flarum.NewResource(flarum.EPost, comment.ID)
	data := obj.Attributes.(*flarum.Post)

	data.Number = comment.Number
	data.ContentType = "comment"
	data.Content = comment.Content
	data.ContentHTML = comment.ContentFmt
	data.CreatedAt = comment.AddTimeFmt

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
