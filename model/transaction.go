package model

import (
	"encoding/json"
	"os"

	"github.com/corvofeng/go-flarum/util"
	"gorm.io/gorm"
)

func ExportTopicTags(gormDB *gorm.DB, fp string) (err error) {
	var tags []Tag
	logger := util.GetLogger()
	err = gormDB.Find(&tags).Error
	if err != nil {
		return
	}

	logger.Infof("Export %d tags to %s", len(tags), fp)
	data, err := json.Marshal(tags)

	if err != nil {
		logger.Error("Marshal tags to json failed", err)
		return
	}
	err = os.WriteFile(fp, data, 0644)

	return
}

// save to file path
func ExportArticles(gormDB *gorm.DB, fp string) (err error) {
	var topics []Topic
	logger := util.GetLogger()
	err = gormDB.Preload("Tags").Preload("BlogMetaData").Find(&topics).Error
	if err != nil {
		return
	}

	logger.Infof("Export %d articles to %s", len(topics), fp)
	data, err := json.Marshal(topics)

	if err != nil {
		logger.Error("Marshal articles to json failed", err)
		return
	}
	err = os.WriteFile(fp, data, 0644)

	return
}

func ImportArticles(gormDB *gorm.DB, fp string) (err error) {
	var topics []Topic
	logger := util.GetLogger()
	data, err := os.ReadFile(fp)
	if err != nil {
		logger.Error("Read file failed", err)
		return
	}
	err = json.Unmarshal(data, &topics)
	if err != nil {
		logger.Error("Unmarshal json to articles failed", err)
		return
	}
	logger.Infof("Import %d articles from %s", len(topics), fp)
	for _, topic := range topics {
		err = gormDB.Create(&topic).Error
		if err != nil {
			logger.Error("Import article failed", topic.ID, topic.Title, err)
			continue
		}
	}
	logger.Infof("Import %d articles successfully", len(topics))
	return
}

func ExportTags(gormDB *gorm.DB, fp string) (err error) {
	var tags []Tag
	logger := util.GetLogger()
	err = gormDB.Find(&tags).Error
	if err != nil {
		return
	}

	logger.Infof("Export %d tags to %s", len(tags), fp)
	data, err := json.Marshal(tags)

	if err != nil {
		logger.Error("Marshal tags to json failed", err)
		return
	}
	err = os.WriteFile(fp, data, 0644)

	return
}

func ImportTags(gormDB *gorm.DB, fp string) (err error) {
	var tags []Tag
	logger := util.GetLogger()
	data, err := os.ReadFile(fp)
	if err != nil {
		logger.Error("Read file failed", err)
		return
	}
	err = json.Unmarshal(data, &tags)
	if err != nil {
		logger.Error("Unmarshal json to tags failed", err)
		return
	}
	logger.Infof("Import %d tags from %s", len(tags), fp)
	for _, tag := range tags {
		err = gormDB.Create(&tag).Error
		if err != nil {
			logger.Error("Import tag failed", tag.ID, tag.Name, err)
			continue
		}
	}
	logger.Infof("Import %d tags successfully", len(tags))
	return
}

func ExportComments(gormDB *gorm.DB, fp string) (err error) {
	var replies []Reply
	logger := util.GetLogger()
	err = gormDB.Find(&replies).Error
	if err != nil {
		return
	}

	logger.Infof("Export %d comments to %s", len(replies), fp)
	data, err := json.Marshal(replies)

	if err != nil {
		logger.Error("Marshal comments to json failed", err)
		return
	}
	err = os.WriteFile(fp, data, 0644)

	return
}

func ImportComments(gormDB *gorm.DB, fp string) (err error) {
	var replies []Reply
	logger := util.GetLogger()
	data, err := os.ReadFile(fp)
	if err != nil {
		logger.Error("Read file failed", err)
		return
	}
	err = json.Unmarshal(data, &replies)
	if err != nil {
		logger.Error("Unmarshal json to comments failed", err)
		return
	}
	logger.Infof("Import %d comments from %s", len(replies), fp)
	for _, reply := range replies {
		err = gormDB.Create(&reply).Error
		if err != nil {
			logger.Error("Import comment failed", reply.ID, reply.Content, err)
			continue
		}
	}
	logger.Infof("Import %d comments successfully", len(replies))
	return
}

func ExportUsers(gormDB *gorm.DB, fp string) (err error) {
	var users []User
	logger := util.GetLogger()
	err = gormDB.Find(&users).Error
	if err != nil {
		return
	}

	logger.Infof("Export %d users to %s", len(users), fp)
	data, err := json.Marshal(users)

	if err != nil {
		logger.Error("Marshal users to json failed", err)
		return
	}
	err = os.WriteFile(fp, data, 0644)

	return
}

func ImportUsers(gormDB *gorm.DB, fp string) (err error) {
	var users []User
	logger := util.GetLogger()
	data, err := os.ReadFile(fp)
	if err != nil {
		logger.Error("Read file failed", err)
		return
	}
	err = json.Unmarshal(data, &users)
	if err != nil {
		logger.Error("Unmarshal json to users failed", err)
		return
	}
	logger.Infof("Import %d users from %s", len(users), fp)
	for _, user := range users {
		err = gormDB.Create(&user).Error
		if err != nil {
			logger.Error("Import user failed", user.ID, user.Name, err)
			continue
		}
	}
	logger.Infof("Import %d users successfully", len(users))
	return
}
