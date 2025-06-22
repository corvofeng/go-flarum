package main

import (
	"flag"
	"os"
	"strings"

	"github.com/corvofeng/go-flarum/model"
	"github.com/corvofeng/go-flarum/system"
	"github.com/corvofeng/go-flarum/util"
)

func main() {
	configFile := flag.String("config", "./config/config.yaml", "full path of config.yaml file")
	logLevel := flag.String("lvl", "INFO", "DEBUG LEVEL")
	export := flag.String("export", "", "export data to file, e.g. export=article,comment,user")
	importer := flag.String("import", "", "import data to file, e.g. impoert=article,comment,user")

	flag.Parse()
	util.InitLogger(*logLevel)
	logger := util.GetLogger()

	c := system.LoadConfig(*configFile)
	app := &system.Application{}

	app.Init(c, os.Args[0])
	defer app.Close()
	model.RankMapInit(app.GormDB, app.RedisDB)
	// app.GormDB.AutoMigrate(flarum.Preferences{})

	// app.GormDB.AutoMigrate(model.User{})
	// app.GormDB.AutoMigrate(model.CommentBase{})
	// app.GormDB.AutoMigrate(model.ArticleBase{})

	article, _ := model.SQLArticleGetByID(app.GormDB, app.RedisDB, 4)
	article.CleanCache()
	// cmt, _ := model.SQLCommentByID(app.GormDB,    , app.RedisDB, 158, 1)
	// pageInfo := model.SQLCommentListByTopic(
	// 	app.GormDB, app.RedisDB, article.ID, 100, app.Cf.Site.TimeZone)
	// for _, c := range pageInfo.Items {
	// 	fmt.Println(c.ID, c.CreatedAt.UTC().String())
	// }

	for _, e := range strings.Split(*export, ",") {
		switch e {
		case "article":
			if err := model.ExportArticles(app.GormDB, "./export/articles.json"); err != nil {
				logger.Error("Export articles failed:", err)
			} else {
				logger.Info("Export articles successfully")
			}
		case "comment":
			if err := model.ExportComments(app.GormDB, "./export/comments.json"); err != nil {
				logger.Error("Export comments failed:", err)
			} else {
				logger.Info("Export comments successfully")
			}
		case "user":
			if err := model.ExportUsers(app.GormDB, "./export/users.json"); err != nil {
				logger.Error("Export users failed:", err)
			} else {
				logger.Info("Export users successfully")
			}
		case "tag":
			if err := model.ExportTags(app.GormDB, "./export/tags.json"); err != nil {
				logger.Error("Export tags failed:", err)
			} else {
				logger.Info("Export tags successfully")
			}
		default:
			logger.Warningf("Unknown export type: %s", e)
		}
	}

	for _, i := range strings.Split(*importer, ",") {
		switch i {
		case "article":
			if err := model.ImportArticles(app.GormDB, "./export/articles.json"); err != nil {
				logger.Error("Import articles failed:", err)
			} else {
				logger.Info("Import articles successfully")
			}
		case "comment":
			if err := model.ImportComments(app.GormDB, "./export/comments.json"); err != nil {
				logger.Error("Import comments failed:", err)
			} else {
				logger.Info("Import comments successfully")
			}
		case "user":
			if err := model.ImportUsers(app.GormDB, "./export/users.json"); err != nil {
				logger.Error("Import users failed:", err)
			} else {
				logger.Info("Import users successfully")
			}
		case "tag":
			if err := model.ImportTags(app.GormDB, "./export/tags.json"); err != nil {
				logger.Error("Import tags failed:", err)
			} else {
				logger.Info("Import tags successfully")
			}

		default:
			logger.Warningf("Unknown import type: %s", i)
		}
	}

	// 调试tags
	// fmt.Println(model.SQLGetTags(app.GormDB))
	// fmt.Println(model.SQLGetTagByUrlName(app.GormDB, "r_funny"))

	// 测试通过tag查找帖子
	// tag, err := model.SQLGetTagByUrlName(app.GormDB, "r_funny")
	// fmt.Println(tag, err)
	// var topics []model.Topic
	// // fmt.Println(app.GormDB.Debug().Model(&topics).Association("Tags").Find(&[]model.Tag{tag}))
	// fmt.Println(app.GormDB.Debug().Limit(10).Offset(0).Model(&tag).Association("Topics").Find(&topics))
	// fmt.Println(topics)

	// article.CleanCache()
	// app.RedisDB.Del(article.toKeyForComments())
	// pageInfo := model.SQLArticleGetByCID(
	// 	app.GormDB,    , app.RedisDB, 0, 1, 10,
	// 	app.Cf.Site.TimeZone,
	// )
	// logger.Info(pageInfo)

	// var user model.User
	// result := app.GormDB.First(&user, 9999)
	// logger.Info(user, result.Error)

	logger.Info("Func utils finished")
}
