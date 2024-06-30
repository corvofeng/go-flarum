package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/corvofeng/go-flarum/model"
	"github.com/corvofeng/go-flarum/model/flarum"
	"github.com/corvofeng/go-flarum/system"
	"github.com/corvofeng/go-flarum/util"
)

func main() {
	configFile := flag.String("config", "config/config.yaml", "full path of config.yaml file")
	logLevel := flag.String("lvl", "INFO", "DEBUG LEVEL")
	initDB := flag.Bool("initdb", false, "init db")

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

	logger.Info(app.GormDB, app.RedisDB)
	model.RankMapInit(app.GormDB, app.RedisDB)
	model.SQLTopicGetByTag(
		app.GormDB, app.RedisDB, 0, 1, 10,
		app.Cf.Site.TimeZone,
	)

	u, err := model.SQLUserRegister(
		app.GormDB,
		"corvofeng", "corvofeng@gmail.com",
		"NoPassWord",
	)
	fmt.Println(u, err)

	if *initDB {
		tag := model.Tag{
			Name:    "root",
			URLName: "r_root",
			Color:   "#000000",
		}
		tag.CreateFlarumTag(app.GormDB)

		u, _ := model.SQLUserRegister(app.GormDB, "root", "", "NoPassword")

		topic := model.Topic{
			Title:  "test",
			UserID: u.ID,
		}

		tags, _ := model.SQLGetTags(app.GormDB)
		tagsArray := flarum.RelationArray{}
		for _, tag := range tags {
			tagID := tag.ID
			tagsArray.Data = append(
				tagsArray.Data,
				flarum.InitBaseResources(uint64(tagID), "tags"),
			)
		}

		topic.CreateFlarumTopic(app.GormDB, tagsArray)
	}

	logger.Info("Migrate the db")
}
