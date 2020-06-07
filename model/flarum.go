package model

import (
	"goyoubbs/model/flarum"
)

// import "goyoubbs/model/flarum"

// // For flarum adapter

// // FlarumBaseAttributes flarum资源的属性
// type FlarumBaseAttributes struct {
// }

// // FlarumLocales flarum 语言信息
// type FlarumLocales struct {
// 	EN string `json:"en"`
// }

// FlarumCreateForumInfo 从SiteInfo创建ForumInfo
func FlarumCreateForumInfo(
	siteConf SiteConf,
	siteInfo SiteInfo,
) flarum.Resource {
	obj := flarum.NewResource(flarum.EForum, 1)
	data := (obj.Attributes).(*flarum.Forum)
	data.DefaultRoute = "/all"
	data.Title = siteConf.Name
	data.Description = siteConf.Desc
	data.APIURL = FlarumAPIPath
	data.AdminURL = "/admin"
	// data.BasePath = "http://192.168.101.35:8082"
	data.BaseURL = "http://192.168.101.35:8082"
	data.AllowSignUp = true

	data.BasePath = ""
	// data.BaseURL = "/"

	obj.Relationships = flarum.ForumRelations{}

	return obj
}

// FlarumCreateTag 创建tag资源
func FlarumCreateTag(cat Category) flarum.Resource {
	obj := flarum.NewResource(flarum.ETAG, cat.ID)
	data := obj.Attributes.(*flarum.Tag)
	data.Name = cat.Name
	data.DiscussionCount = cat.Articles
	data.IsHidden = cat.Hidden

	// data.Candd
	data.Color = "#B72A2A"
	data.Icon = "fas fa-wrench"

	return obj
}

// FlarumCreateDiscussion 创建帖子资源
func FlarumCreateDiscussion(article ArticleListItem) flarum.Resource {
	obj := flarum.NewResource(flarum.EDiscussion, article.ID)
	data := obj.Attributes.(*flarum.Discussion)
	data.Title = article.Title
	// data.CommentCount = article.Comments
	data.CommentCount = 10
	data.ParticipantCount = 5
	data.LastPostNumber = 1
	obj.Relationships = flarum.DiscussionRelations{
		User: flarum.RelationDict{
			Data: flarum.InitBaseResources(article.CID, "users"),
		},
	}

	return obj
}

// FlarumCreateDiscussionFromArticle 创建帖子资源
func FlarumCreateDiscussionFromArticle(article Article) flarum.Resource {
	obj := flarum.NewResource(flarum.EDiscussion, article.ID)
	data := obj.Attributes.(*flarum.Discussion)
	data.Title = article.Title
	// data.CommentCount = article.Comments
	data.CommentCount = 1
	data.ParticipantCount = 1
	data.LastPostNumber = 1
	data.CreatedAt = "2019-05-30T14:26:02+00:00"
	// data.LastReadAt = "2020-05-31T14:26:02+00:00"
	// data.LastPostedAt = "2019-06-29T11:20:01Z"
	data.LastPostedAt = "2020-05-31T12:49:51+00:00"
	// data.CanRename = true
	// data.CanReply = true
	// data.IsApproved = true
	// data.IsSticky = true

	obj.BindRelations(
		"User",
		flarum.RelationDict{Data: flarum.InitBaseResources(article.UID, "users")},
	)

	return obj
}

// FlarumCreateUser 创建用户资源
func FlarumCreateUser(article ArticleListItem) flarum.Resource {
	obj := flarum.NewResource(flarum.EBaseUser, article.UID)
	data := obj.Attributes.(*flarum.BaseUser)
	data.Username = article.Cname
	data.AvatarURL = article.Avatar

	return obj
}

// FlarumCreateCurrentUser 创建用户资源
func FlarumCreateCurrentUser(user User) flarum.Resource {
	obj := flarum.NewResource(flarum.ECurrentUser, user.ID)
	data := obj.Attributes.(*flarum.CurrentUser)
	data.Username = user.Name
	data.Displayname = user.Name
	data.AvatarURL = user.Avatar
	data.IsEmailConfirmed = true
	return obj
}

// FlarumCreateUserFromComments 通过评论信息创建用户资源
func FlarumCreateUserFromComments(comment CommentListItem) flarum.Resource {

	obj := flarum.NewResource(flarum.ECurrentUser, comment.UID)
	data := obj.Attributes.(*flarum.CurrentUser)
	data.Displayname = comment.UserName
	data.Username = comment.UserName
	data.AvatarURL = comment.Avatar

	data.LastSeenAt = "2020-06-02T04:56:23+00:00"

	data.CommentCount = 20
	data.DiscussionCount = 3

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
			Data: flarum.InitBaseResources(comment.Aid, "discussions"),
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
