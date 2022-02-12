package main

import (
	"flag"
	"fmt"
	"os"
	"zoe/model"
	"zoe/system"
	"zoe/util"
)

func main() {
	configFile := flag.String("config", "../../config/config.yaml", "full path of config.yaml file")
	logLevel := flag.String("lvl", "INFO", "DEBUG LEVEL")

	flag.Parse()
	util.InitLogger(*logLevel)
	logger := util.GetLogger()

	c := system.LoadConfig(*configFile)
	app := &system.Application{}

	app.Init(c, os.Args[0])
	defer app.Close()
	model.RankMapInit(app.GormDB, app.MySQLdb, app.RedisDB)
	// app.GormDB.AutoMigrate(flarum.Preferences{})

	// app.GormDB.AutoMigrate(model.User{})
	// app.GormDB.AutoMigrate(model.CommentBase{})
	// app.GormDB.AutoMigrate(model.ArticleBase{})

	article, _ := model.SQLArticleGetByID(app.GormDB, app.MySQLdb, app.RedisDB, 938)
	fmt.Println(article.GetCommentIDList(app.RedisDB))
	// cmt, _ := model.SQLCommentByID(app.GormDB, app.MySQLdb, app.RedisDB, 158, 1)
	pageInfo := model.SQLCommentListByTopic(
		app.GormDB, app.MySQLdb, app.RedisDB, article.ID, 100, app.Cf.Site.TimeZone)
	for _, c := range pageInfo.Items {
		fmt.Println(c.ID, c.CreatedAt.UTC().String())
	}

	// article.CleanCache()
	// app.RedisDB.Del(article.toKeyForComments())
	// pageInfo := model.SQLArticleGetByCID(
	// 	app.GormDB, app.MySQLdb, app.RedisDB, 0, 1, 10,
	// 	app.Cf.Site.TimeZone,
	// )
	// logger.Info(pageInfo)

	// var user model.User
	// result := app.GormDB.First(&user, 9999)
	// logger.Info(user, result.Error)

	logger.Info("Func utils finished")
}
