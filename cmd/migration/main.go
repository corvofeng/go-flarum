package main

import (
	"flag"
	"os"
	"zoe/model"
	"zoe/system"
	"zoe/util"
)

func main() {
	configFile := flag.String("config", "config/config.yaml", "full path of config.yaml file")
	logLevel := flag.String("lvl", "INFO", "DEBUG LEVEL")

	flag.Parse()
	util.InitLogger(*logLevel)
	logger := util.GetLogger()

	c := system.LoadConfig(*configFile)
	app := &system.Application{}

	app.Init(c, os.Args[0])
	// app.GormDB.AutoMigrate(flarum.Preferences{})
	app.GormDB.AutoMigrate(model.User{})
	var user model.User

	result := app.GormDB.First(&user, 9999)
	logger.Info(user, result.Error)

	logger.Info("Migrate the db")
}
