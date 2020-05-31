package system

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"net/url"
	"strings"

	"database/sql"

	"goyoubbs/model"
	"goyoubbs/util"

	"github.com/ego008/youdb"
	"github.com/gorilla/securecookie"
	logging "github.com/op/go-logging"
	"github.com/qiniu/api.v7/storage"
	"github.com/weint/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-redis/redis/v7"
)

// Application 应用数据库以及外部服务
type Application struct {
	Cf      *model.AppConf
	Db      *youdb.DB
	RedisDB *redis.Client
	MySQLdb *sql.DB
	MongoDB *mongo.Client
	Sc      *securecookie.SecureCookie
	QnZone  *storage.Zone
	Logger  *logging.Logger
	Rand    *rand.Rand // 负责处理随机数
}

// LoadConfig 从文件中初始化程序配置
func LoadConfig(filename string) *config.Engine {
	c := &config.Engine{}
	c.Load(filename)
	return c
}

// Init ， 连接数据库
func (app *Application) Init(c *config.Engine, currentFilePath string) {
	// .. version_changed: 2019-11-09
	// 添加 redis, 目前redis只用于缓存数据，理论上不能包含数据结构

	mcf := &model.MainConf{}
	c.GetStruct("Main", mcf)
	logger := util.GetLogger()

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
	// db, err := youdb.Open(mcf.Youdb)
	// if err != nil {
	// 	logger.Fatalf("Connect Error: %v", err)
	// }
	dbStr := app.Cf.Main.MySQLURL //  fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", mcf.MySQL_USER, mcf.MySQL_PASS, mcf.MySQL_HOST, mcf.MySQL_PORT, mcf.MySQL_DB)
	logger.Debugf("Get mysql db url: %s", dbStr)

	rdsClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", mcf.RedisHost, mcf.RedisPort),
		Password: mcf.RedisPass,
		DB:       mcf.RedisDB,
	})
	pong, err := rdsClient.Ping().Result()
	if err != nil {
		logger.Errorf("Connect redis error, %s", err)
		return
	}
	logger.Debug(pong, err)

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(mcf.MongoURL))
	logger.Debugf("Get mongo db url: %s", mcf.MongoURL)
	if err != nil {
		logger.Errorf("Connect mongo error, %s", err)
		return
	}

	sqlDb, err := sql.Open("mysql", dbStr)
	sqlDb.SetConnMaxLifetime(time.Minute * 10)
	if err != nil {
		logger.Errorf("Connect mysql error, %s", err)
		return
	}
	app.Db = nil
	app.MySQLdb = sqlDb
	app.RedisDB = rdsClient
	app.MongoDB = mongoClient
	app.Logger = util.GetLogger()

	// set main node
	// db.Hset("keyValue", []byte("main_category"), []byte(scf.MainNodeIDs))

	app.Sc = securecookie.New(
		[]byte(app.Cf.Main.SCHashKey),
		[]byte(app.Cf.Main.SCBlockKey),
		// securecookie.GenerateRandomKey(64),
		// securecookie.GenerateRandomKey(32),
	)
	//app.Sc.SetSerializer(securecookie.JSONEncoder{})

	app.Logger.Debug("youdb Connect to", mcf.Youdb)
}

// IsFlarum 当前论坛是否为flarum风格
func (app *Application) IsFlarum() bool {
	return app.Cf.Main.ServerStyle == "flarum"
}

// Close 清理程序连接
func (app *Application) Close() {
	if app.Db != nil {
		app.Db.Close()
	}
	app.MySQLdb.Close()
	app.RedisDB.Close()
	app.Logger.Debug("db cloded")
}
