package model

import (
	"goyoubbs/model/flarum"
	"goyoubbs/util"
)

// FlarumCreateForumInfo 从SiteInfo创建ForumInfo
// tags 当前站点具有的标签集合, TODO: 缓存
func FlarumCreateForumInfo(
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
	data.AdminURL = "/admin"
	// data.BasePath = "http://192.168.101.35:8082"
	data.BaseURL = mainConf.BaseURL
	data.CanStartDiscussion = true
	data.AllowSignUp = true
	data.WelcomeMessage = "这是一个简单的小站"
	data.WelcomeTitle = "用作测试"
	data.Debug = true
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

	data.Color = "#B72A2A"
	data.Icon = "fas fa-wrench"

	// please refer to #11
	if cat.ParentID == 0 {
		obj.Relationships = flarum.TagRelations{}
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
	data.LastPostID = article.LastPostID
	data.LastPostNumber = lastComment.Number
	data.FirstPostID = article.FirstPostID
	data.LastPostedAt = lastComment.AddTimeFmt

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

// // FlarumCreateUser 创建用户资源
// func FlarumCreateUser(article ArticleListItem) flarum.Resource {
// 	obj := flarum.NewResource(flarum.EBaseUser, article.UID)
// 	data := obj.Attributes.(*flarum.BaseUser)
// 	data.Username = article.Cname
// 	data.Displayname = article.Cname
// 	data.AvatarURL = article.Avatar

// 	obj.BindRelations(
// 		"Groups",
// 		flarum.RelationArray{
// 			Data: []flarum.BaseRelation{},
// 		},
// 	)
// 	return obj
// }

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
	data.IsEmailConfirmed = true
	data.JoinTime = util.TimeFmt(user.RegTime, util.TIME_FMT, 0)
	obj.BindRelations(
		"Groups",
		flarum.RelationArray{
			Data: []flarum.BaseRelation{},
		},
	)
	return obj
}

// FlarumCreateUserFromComments 通过评论信息创建用户资源
func FlarumCreateUserFromComments(comment CommentListItem) flarum.Resource {

	obj := flarum.NewResource(flarum.ECurrentUser, comment.UID)
	data := obj.Attributes.(*flarum.CurrentUser)
	data.Displayname = comment.UserName
	data.Username = comment.UserName
	data.AvatarURL = comment.Avatar

	// data.LastSeenAt = "2020-06-02T04:56:23+00:00"
	// data.CommentCount = 20
	// data.DiscussionCount = 3

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
func FlarumCreatePost(comment CommentListItem) flarum.Resource {
	obj := flarum.NewResource(flarum.EPost, comment.ID)
	data := obj.Attributes.(*flarum.Post)

	data.Number = comment.Number
	data.ContentType = "comment"
	data.Content = comment.Content
	data.ContentHTML = comment.ContentFmt
	data.CreatedAt = comment.AddTimeFmt
	data.CanLike = true
	data.IPAddress = "1.2.3.4"
	data.CanHide = true
	data.IsApproved = true

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
		flarum.RelationArray{
			Data: []flarum.BaseRelation{},
		},
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

	return obj
}

// FlarumCreatePostRelations 创建关系结构
func FlarumCreatePostRelations(postArr []flarum.Resource) flarum.IRelation {
	var obj flarum.RelationArray
	for _, p := range postArr {
		obj.Data = append(
			obj.Data,
			flarum.InitBaseResources(p.GetID(), p.Type),
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
