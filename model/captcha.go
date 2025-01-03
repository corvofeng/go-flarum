package model

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/corvofeng/go-flarum/util"

	"github.com/dchest/captcha"
	"github.com/go-redis/redis/v7"
)

// redisStore is an internal store for captcha ids and their values.
type redisStore struct {
	redisDB *redis.Client
}

// SetCaptchaUseRedisStore use redis to store captcha
func SetCaptchaUseRedisStore(redisDB *redis.Client) {
	captcha.SetCustomStore(&redisStore{
		redisDB: redisDB,
	})
}

// NewCaptcha 产生新的验证码图片
func NewCaptcha(path string) string {
	captchaID := captcha.New()
	SaveImage(path, captchaID)
	return captchaID
}

// Set sets the digits for the captcha id.
func (rs *redisStore) Set(id string, digits []byte) {
	logger := util.GetLogger()
	logger.Debugf("Set captcha %s with %v", id, digits)

	rs.redisDB.Set(id, string(digits), time.Second*200) // 3 minutes
}

// Get returns stored digits for the captcha id. Clear indicates
// whether the captcha must be deleted from the store.
func (rs *redisStore) Get(id string, clear bool) (digits []byte) {
	// logger := util.GetLogger()
	// logger.Debugf("wanna get captcha %s", id)

	rlt, err := rs.redisDB.Get(id).Result()
	if err == redis.Nil {
		return []byte("")
	}

	return []byte(rlt)
}

// SaveImage to static dir
func SaveImage(path, id string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
	savePath := filepath.Join(path, fmt.Sprintf("%s.png", id))

	f, err := os.Create(savePath)

	if util.CheckError(err, "保存验证码") {
		return
	}

	captcha.WriteImage(f, id, captcha.StdWidth, captcha.StdHeight)
	f.Sync()
}
