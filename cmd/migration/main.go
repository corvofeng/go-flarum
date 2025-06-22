package main

import (
	"flag"
	"os"

	"github.com/corvofeng/go-flarum/model"
	"github.com/corvofeng/go-flarum/system"
	"github.com/corvofeng/go-flarum/util"
)

func main() {
	configFile := flag.String("config", "config/config.yaml", "full path of config.yaml file")
	logLevel := flag.String("lvl", "INFO", "DEBUG LEVEL")
	initDB := flag.Bool("initdb", false, "init db")
	queryBlog := flag.Uint64("queryBlog", 0, "query blog meta by topic id")

	flag.Parse()
	util.InitLogger(*logLevel)
	logger := util.GetLogger()

	c := system.LoadConfig(*configFile)
	app := &system.Application{}

	app.Init(c, os.Args[0])
	defer app.Close()
	// app.GormDB.AutoMigrate(flarum.Preferences{})

	util.CheckError(app.GormDB.AutoMigrate(model.User{}), "migrate user")
	util.CheckError(app.GormDB.AutoMigrate(model.Tag{}), "migrate tag")
	util.CheckError(app.GormDB.AutoMigrate(model.Topic{}), "migrate topic")
	util.CheckError(app.GormDB.AutoMigrate(model.Reply{}), "migrate reply")
	util.CheckError(app.GormDB.AutoMigrate(model.ReplyLikes{}), "migrate reply likes")
	util.CheckError(app.GormDB.AutoMigrate(model.UserTopic{}), "migrate user topic")
	util.CheckError(app.GormDB.AutoMigrate(model.ActionRecord{}), "migrate action record")
	util.CheckError(app.GormDB.AutoMigrate(model.UserFiles{}), "migrate userfiles record")
	util.CheckError(app.GormDB.AutoMigrate(model.BlogMeta{}), "migrate blog meta record")

	logger.Info(app.GormDB, app.RedisDB)
	model.RankMapInit(app.GormDB, app.RedisDB)

	if *initDB {
		tag := model.Tag{
			Name:    "root",
			URLName: "r_root",
			Color:   "#000000",
		}

		if t, err := model.SQLGetTagByUrlName(app.GormDB, tag.URLName); t.ID == 0 {
			err = tag.CreateFlarumTag(app.GormDB)
			util.CheckError(err, "create tag")
		} else {
			tag = t
			logger.Warningf("Tag %s already exists, skip creating", tag.Name)
		}

		// root:1234
		user := model.User{
			Name:     "root",
			Password: "81dc9bdb52d04dc20036dbd8313ed055",
			Admin:    true,
		}
		if u, err := model.SQLUserGetByName(app.GormDB, user.Name); u.ID == 0 {
			err = user.CreateFlarumUser(app.GormDB)
			util.CheckError(err, "create user")
		} else {
			user = u
			logger.Warningf("User %s already exists, skip creating", u.Name)
		}

		tags, _ := model.SQLGetTags(app.GormDB)
		topic := model.Topic{
			Title:  "Welcome to Go Flarum",
			UserID: user.ID,
			Content: `# Welcome to Flarum

> The default user is root:1234
You can login with this user, or create your own user.
	
This topic is created by migration tool, please delete it if don't want to see this
			`,
			Tags: tags,

			CommentCount: 1,
		}
		if top, err := model.SQLGetTopicGetByTitle(app.GormDB, topic.Title); top.ID == 0 {
			_, err = topic.CreateFlarumTopic(app.GormDB)
			util.CheckError(err, "create topic")
		} else {
			topic = top
			logger.Warningf("Topic %s already exists, skip creating", topic.Title)
		}

		// blogMeta := model.BlogMeta{
		// 	TopicID: topic.ID,
		// 	Summary: "topic summary",
		// }

		// b, err := blogMeta.CreateFlarumBlogMeta(app.GormDB)
		// if util.CheckError(err, "create blog meta") {
		// 	return
		// }
		// logger.Info("Create blog meta", b)
	}

	if *queryBlog != 0 {
		article, err := model.SQLArticleGetByID(app.GormDB, app.RedisDB, *queryBlog)
		if util.CheckError(err, "get blog") {
			return
		}
		logger.Info("Get blog %+v %+v", article, article.BlogMetaData)
	}

	logger.Info("Migrate the db")
}
