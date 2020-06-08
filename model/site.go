package model

import (
	"time"

	"github.com/go-redis/redis/v7"
)

const (
	// FlarumAPIPath flarum 的api位置
	FlarumAPIPath = "/api/v1/flarum"
)

// MainConf 主配置
type MainConf struct {
	HTTPPort  int
	HttpsOn   bool   // TODO: 废弃
	HttpsPort int    // TODO: 废弃
	Domain    string // 若启用https 则该domain 为注册的域名，eg: domain.com、www.domain.com

	BaseURL string

	// 数据库地址
	MySQLURL string
	MongoURL string
	RedisURL string

	PubDir         string
	LocaleDir      string
	ViewDir        string
	Youdb          string // TODO: 废弃
	ServerStyle    string // 选择使用的样式
	ServerName     string
	CookieSecure   bool
	CookieHttpOnly bool
	OldSiteDomain  string
	TLSCrtFile     string
	TLSKeyFile     string

	// secure cookie 初始化时需要
	SCHashKey  string
	SCBlockKey string
}

// SiteConf 站点配置
type SiteConf struct {
	GoVersion  string
	MD5Sums    string
	Name       string
	Desc       string
	AdminEmail string
	MainDomain string // 上传图片后添加网址前缀, eg: http://domian.com 、http://234.21.35.89:8082

	CDNBaseURL string // 静态文件cdn地址

	MainNodeIDs       string
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

// AppConf 应用配置文件
type AppConf struct {
	Main *MainConf
	Site *SiteConf
}

// SiteInfo 存储在数据库中的信息
type SiteInfo struct {
	Days     uint64
	UserNum  uint64
	NodeNum  uint64
	TagNum   uint64
	PostNum  uint64
	ReplyNum uint64
}

// GetDays 获取从建站开始, 到目前的天数, 用于主页中的显示
func GetDays(redisDB *redis.Client) uint64 {

	siteCreateTime, err := redisDB.Get("site_create_time").Uint64()
	if err != nil {
		siteCreateTime = 1557585456 // 2019-05-11 22:37:36 +0800 HKT
	}
	then := time.Unix(int64(siteCreateTime), 0)
	diff := time.Now().UTC().Sub(then)
	return uint64(diff.Hours()/24) + 1
}

// GetSiteInfo 直接获取网站信息
func GetSiteInfo(redisDB *redis.Client) SiteInfo {
	si := SiteInfo{}
	si.Days = GetDays(redisDB)
	// si.UserNum = db.Hsequence("user")
	// si.NodeNum = db.Hsequence("category")
	// si.TagNum = db.Hsequence("tag")
	// si.PostNum = db.Hsequence("article")
	// si.ReplyNum = db.Hget("count", []byte("comment_num")).Uint64()

	return si
}
