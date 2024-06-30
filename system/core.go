package system

import (
	"math/rand"
	"runtime"
	"time"

	"net/url"
	"strings"

	"database/sql"

	"github.com/corvofeng/go-flarum/model"
	"github.com/corvofeng/go-flarum/util"

	"github.com/gorilla/securecookie"
	logging "github.com/op/go-logging"
	"github.com/weint/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/go-redis/redis/v7"
)

// Application 应用数据库以及外部服务
type Application struct {
	Cf      *model.AppConf
	RedisDB *redis.Client
	MySQLdb *sql.DB
	// MongoDB *mongo.Client
	Sc *securecookie.SecureCookie
	// QnZone  *storage.Zone
	Logger *logging.Logger
	Rand   *rand.Rand // 负责处理随机数
	GormDB *gorm.DB
}

// LoadConfig 从文件中初始化程序配置
func LoadConfig(filename string) *config.Engine {
	c := &config.Engine{}
	err := c.Load(filename)
	logger := util.GetLogger()
	if err != nil {
		logger.Error("读取配置文件失败:", err)
	}
	return c
}

// Init ， 连接数据库
func (app *Application) Init(c *config.Engine, currentFilePath string) {
	// .. version_changed: 2019-11-09
	// 添加 redis, 目前redis只用于缓存数据，理论上不能包含数据结构

	mcf := &model.MainConf{}
	c.GetStruct("Main", mcf)
	logger := util.GetLogger()
	app.Logger = logger

	app.Rand = rand.New(rand.NewSource(time.Now().Unix()))

	// check domain
	if strings.HasPrefix(mcf.Domain, "http") {
		dm, err := url.Parse(mcf.Domain)
		if err != nil {
			logger.Fatal("domain fmt err", err)
		}
		mcf.Domain = dm.Host
	} else {
		mcf.Domain = strings.Trim(mcf.Domain, "/")
	}

	scf := &model.SiteConf{}
	c.GetStruct("Site", scf)
	scf.GoVersion = runtime.Version()
	fMd5, _ := util.HashFileMD5(currentFilePath)
	scf.MD5Sums = fMd5
	scf.MainDomain = strings.Trim(scf.MainDomain, "/")
	if scf.TimeZone < -12 || scf.TimeZone > 12 {
		scf.TimeZone = 0
	}
	if scf.UploadMaxSize < 1 {
		scf.UploadMaxSize = 1
	}
	scf.UploadMaxSizeByte = int64(scf.UploadMaxSize) << 20

	app.Cf = &model.AppConf{mcf, scf}

	logger.Debugf("Get redis db url: %s", mcf.RedisURL)
	opt, err := redis.ParseURL(mcf.RedisURL)
	if err != nil {
		panic(err)
	}
	rdsClient := redis.NewClient(opt)
	pong, err := rdsClient.Ping().Result()
	if err != nil {
		logger.Errorf("Connect redis error, %s", err)
		return
	}
	logger.Debug(pong, err)

	logger.Debugf("Get mongo db url: %s", mcf.MongoURL)
	// mongoClient, err := mongo.NewClient(options.Client().ApplyURI(mcf.MongoURL))
	// if err != nil {
	// 	logger.Errorf("Connect mongo error, %s", err)
	// 	return
	// }

	logger.Debugf("Get mysql db url: %s", mcf.MySQLURL)
	sqlDb, err := sql.Open("mysql", mcf.MySQLURL)
	sqlDb.SetConnMaxLifetime(time.Minute * 10)
	if err != nil {
		logger.Errorf("Connect mysql error, %s", err)
		return
	}

	app.MySQLdb = sqlDb
	app.RedisDB = rdsClient
	app.GormDB, err = gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDb,
	}), &gorm.Config{})
	util.CheckError(err, "gorm open error")
	app.Sc = securecookie.New(
		[]byte(app.Cf.Main.SCHashKey),
		[]byte(app.Cf.Main.SCBlockKey),
	)
}

// IsFlarum 当前论坛是否为flarum风格
func (app *Application) IsFlarum() bool {
	return app.Cf.Main.ServerStyle == "flarum"
}

func (app *Application) CanServeAdmin() bool {
	return true
}

// Close 清理程序连接
func (app *Application) Close() {
	if app.MySQLdb != nil {
		app.MySQLdb.Close()
		app.MySQLdb = nil
	}
	if app.RedisDB != nil {
		app.RedisDB.Close()
		app.RedisDB = nil
	}
	app.Logger.Debug("db cloded")
}
