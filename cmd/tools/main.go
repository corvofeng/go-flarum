package main

import (
	"flag"
	"fmt"
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

	article, _ := model.SQLArticleGetByID(app.GormDB,    , app.RedisDB, 938)
	fmt.Println(article.GetCommentIDList(app.RedisDB))
	// cmt, _ := model.SQLCommentByID(app.GormDB,    , app.RedisDB, 158, 1)
	pageInfo := model.SQLCommentListByTopic(
		app.GormDB,    , app.RedisDB, article.ID, 100, app.Cf.Site.TimeZone)
	for _, c := range pageInfo.Items {
		fmt.Println(c.ID, c.CreatedAt.UTC().String())
	}

	// 调试tags
	// fmt.Println(model.SQLGetTags(app.GormDB))
	// fmt.Println(model.SQLGetTagByUrlName(app.GormDB, "r_funny"))

	// 测试通过tag查找帖子
	tag, err := model.SQLGetTagByUrlName(app.GormDB, "r_funny")
	fmt.Println(tag, err)
	var topics []model.Topic
	// fmt.Println(app.GormDB.Debug().Model(&topics).Association("Tags").Find(&[]model.Tag{tag}))
	fmt.Println(app.GormDB.Debug().Limit(10).Offset(0).Model(&tag).Association("Topics").Find(&topics))
	fmt.Println(topics)

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
