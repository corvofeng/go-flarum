package system

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	"net/url"
	"strings"

	"database/sql"

	"goyoubbs/util"

	"github.com/ego008/youdb"
	"github.com/gorilla/securecookie"
	logging "github.com/op/go-logging"
	"github.com/qiniu/api.v7/storage"
	"github.com/weint/config"
)

type MainConf struct {
	HttpPort       int
	HttpsOn        bool
	Domain         string // 若启用https 则该domain 为注册的域名，eg: domain.com、www.domain.com
	HttpsPort      int
	PubDir         string
	ViewDir        string
	Youdb          string
	CookieSecure   bool
	CookieHttpOnly bool
	OldSiteDomain  string
	TLSCrtFile     string
	TLSKeyFile     string

	// secure cookie 初始化时需要
	SCHashKey  string
	SCBlockKey string
}

type SiteConf struct {
	GoVersion         string
	MD5Sums           string
	Name              string
	Desc              string
	AdminEmail        string
	MainDomain        string // 上传图片后添加网址前缀, eg: http://domian.com 、http://234.21.35.89:8082
	MainNodeIds       string
	TimeZone          int
	HomeShowNum       int
	PageShowNum       int
	TagShowNum        int
	CategoryShowNum   int
	TitleMaxLen       int
	ContentMaxLen     int
	PostInterval      int
	CommentListNum    int
	CommentInterval   int
	Authorized        bool
	RegReview         bool
	CloseReg          bool
	AutoDataBackup    bool
	AutoGetTag        bool
	GetTagApi         string
	QQClientID        int
	QQClientSecret    string
	WeiboClientID     int
	WeiboClientSecret string // eg: "jpg,jpeg,gif,zip,pdf"
	UploadSuffix      string
	UploadImgOnly     bool
	UploadImgResize   bool
	UploadMaxSize     int
	UploadMaxSizeByte int64
	QiniuAccessKey    string
	QiniuSecretKey    string
	QiniuDomain       string
	QiniuBucket       string
	UpyunDomain       string
	UpyunBucket       string
	UpyunUser         string
	UpyunPw           string
}

type AppConf struct {
	Main *MainConf
	Site *SiteConf
}

type Application struct {
	Cf      *AppConf
	Db      *youdb.DB
	MySQLdb *sql.DB
	Sc      *securecookie.SecureCookie
	QnZone  *storage.Zone
	Logger  *logging.Logger
	Rand    *rand.Rand // 负责处理随机数
}

func LoadConfig(filename string) *config.Engine {
	c := &config.Engine{}
	c.Load(filename)
	return c
}

func (app *Application) Init(c *config.Engine, currentFilePath string, sqlDb *sql.DB) {

	mcf := &MainConf{}
	c.GetStruct("Main", mcf)

	app.Rand = rand.New(rand.NewSource(time.Now().Unix()))

	// check domain
	if strings.HasPrefix(mcf.Domain, "http") {
		dm, err := url.Parse(mcf.Domain)
		if err != nil {
			log.Fatal("domain fmt err", err)
		}
		mcf.Domain = dm.Host
	} else {
		mcf.Domain = strings.Trim(mcf.Domain, "/")
	}

	scf := &SiteConf{}
	c.GetStruct("Site", scf)
	scf.GoVersion = runtime.Version()
	fMd5, _ := util.HashFileMD5(currentFilePath)
	scf.MD5Sums = fMd5
	scf.MainDomain = strings.Trim(scf.MainDomain, "/")
	// log.Println("MainDomain:", scf.MainDomain)
	if scf.TimeZone < -12 || scf.TimeZone > 12 {
		scf.TimeZone = 0
	}
	if scf.UploadMaxSize < 1 {
		scf.UploadMaxSize = 1
	}
	scf.UploadMaxSizeByte = int64(scf.UploadMaxSize) << 20

	app.Cf = &AppConf{mcf, scf}
	db, err := youdb.Open(mcf.Youdb)
	if err != nil {
		log.Fatalf("Connect Error: %v", err)
	}
	app.Db = db
	app.MySQLdb = sqlDb
	app.Logger = util.GetLogger()

	// set main node
	db.Hset("keyValue", []byte("main_category"), []byte(scf.MainNodeIds))

	app.Sc = securecookie.New(
		[]byte(app.Cf.Main.SCHashKey),
		[]byte(app.Cf.Main.SCBlockKey),
		// securecookie.GenerateRandomKey(64),
		// securecookie.GenerateRandomKey(32),
	)
	//app.Sc.SetSerializer(securecookie.JSONEncoder{})

	app.Logger.Debug("youdb Connect to", mcf.Youdb)
}

func (app *Application) Close() {
	app.Db.Close()
	app.MySQLdb.Close()
	app.Logger.Debug("db cloded")
}
