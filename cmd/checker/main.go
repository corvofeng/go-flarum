package main

import (
	"flag"
	"os"

	"github.com/corvofeng/go-flarum/model"
	"github.com/corvofeng/go-flarum/system"
	"github.com/corvofeng/go-flarum/util"
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
	model.RankMapInit(app.GormDB,    , app.RedisDB)
	// app.GormDB.AutoMigrate(flarum.Preferences{})

	// app.GormDB.AutoMigrate(model.User{})
	// app.GormDB.AutoMigrate(model.CommentBase{})
	// app.GormDB.AutoMigrate(model.ArticleBase{})

	pageInfo := model.SQLTopicGetByTag(
		app.GormDB,    , app.RedisDB, 0, 1, 10,
		app.Cf.Site.TimeZone,
	)
	logger.Info(pageInfo)

	// var user model.User
	// result := app.GormDB.First(&user, 9999)
	// logger.Info(user, result.Error)

	logger.Info("Func check finished")
}
