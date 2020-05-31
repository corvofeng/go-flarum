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
	obj := flarum.NewResource(flarum.EForum)
	data := (*obj.Attributes).(*flarum.Forum)
	data.DefaultRoute = "/all"
	data.Title = siteConf.Name
	data.Description = siteConf.Desc
	data.APIURL = FlarumAPIPath
	data.AdminURL = "/admin"
	data.BasePath = ""
	data.BaseURL = "/"

	obj.Relationships = flarum.ForumRelations{}
	obj.ID = 1

	return obj
}

// FlarumCreateTag 创建tag资源
func FlarumCreateTag(cat Category) flarum.Resource {
	obj := flarum.NewResource(flarum.ETAG)
	data := (*obj.Attributes).(*flarum.Tag)
	data.Name = cat.Name
	data.ID = cat.ID
	data.DiscussionCount = cat.Articles
	data.IsHidden = cat.Hidden
	obj.ID = data.ID
	return obj
}

// FlarumCreateDiscussion 创建帖子资源
func FlarumCreateDiscussion(article ArticleListItem) flarum.Resource {
	obj := flarum.NewResource(flarum.EDiscussion)
	data := (*obj.Attributes).(*flarum.Discussion)
	data.ID = article.ID
	data.Title = article.Title
	data.CommentCount = article.Comments
	obj.ID = data.ID
	obj.Relationships = flarum.DiscussionRelations{
		User: flarum.RelationDict{
			Data: flarum.BaseRelation{
				ID:   article.CID,
				Type: "users",
			}},
	}

	return obj
}

// FlarumCreateUser 创建用户资源
func FlarumCreateUser(article ArticleListItem) flarum.Resource {
	obj := flarum.NewResource(flarum.EBaseUser)
	data := (*obj.Attributes).(*flarum.BaseUser)
	data.ID = article.CID
	data.Username = article.Cname
	data.AvatarURL = article.Avatar

	obj.ID = data.ID
	return obj
}

// FlarumCreatPost 创建评论
func FlarumCreatPost(comment *Comment) flarum.Resource {
	obj := flarum.NewResource(flarum.EPost)
	data := (*obj.Attributes).(*flarum.Post)
	data.ID = comment.ID
	obj.ID = data.ID

	data.IPAddress = comment.ClientIp
	data.Type = "comment"
	data.Number = comment.Number
	data.Content = comment.Content

	obj.Relationships = flarum.PostRelations{
		User: flarum.RelationDict{
			Data: flarum.BaseRelation{
				Type: "user",
				ID:   comment.UID,
			},
		},
	}

	return obj
}
