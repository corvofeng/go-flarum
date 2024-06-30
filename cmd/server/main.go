package main

import (
	"context"

	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/corvofeng/go-flarum/cronjob"
	"github.com/corvofeng/go-flarum/model"

	"github.com/corvofeng/go-flarum/router"
	"github.com/corvofeng/go-flarum/system"
	"github.com/corvofeng/go-flarum/util"

	ct "github.com/corvofeng/go-flarum/controller"

	_ "github.com/go-sql-driver/mysql"
	goji "goji.io"
	"goji.io/pat"
	// "github.com/go-redis/redis/v7"
)

var GitCommit string

func main() {
	configFile := flag.String("config", "config/config.yaml", "full path of config.yaml file")
	logLevel := flag.String("lvl", "INFO", "DEBUG LEVEL")

	flag.Parse()
	util.InitLogger(*logLevel)
	logger := util.GetLogger()
	if GitCommit == "" {
		GitCommit = "Development"
	}

	logger.Infof("Git version: %s", GitCommit)

	c := system.LoadConfig(*configFile)
	app := &system.Application{}

	app.Init(c, os.Args[0])
	model.RankMapInit(app.GormDB, app.RedisDB)

	// 验证码信息使用Redis保存
	model.SetCaptchaUseRedisStore(app.RedisDB)

	// cron job
	cr := cronjob.BaseHandler{App: app}
	if os.Getenv("type") == "cron" {
		logger.Info("Cron worker start !!!")
		go cr.MainCronJob()
	} else {
		logger.Info("This is not a cron worker")
	}

	root := goji.NewMux()

	mcf := app.Cf.Main

	// static file server
	staticPath := mcf.PubDir
	if len(staticPath) == 0 {
		staticPath = "static"
	}

	h := ct.BaseHandler{App: app}
	root.Handle(pat.New("/static/*"),
		h.OriginMiddleware(
			http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))),
		),
	)

	root.Handle(pat.New("/webpack/*"),
		h.OriginMiddleware(
			http.StripPrefix("/webpack/", http.FileServer(http.Dir(mcf.WebpackDir))),
		),
	)

	root.Handle(pat.New("/*"), router.NewRouter(app))

	// normal http
	// http.ListenAndServe(listenAddr, root)

	// graceful
	// subscribe to SIGINT signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	var srv *http.Server

	// http
	srv = &http.Server{Addr: ":" + strconv.Itoa(mcf.HTTPPort), Handler: root}
	// srv = &http.Server{Addr: ":" + *httpPort, Handler: root}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	logger.Debug("Web server Listen port", strconv.Itoa(mcf.HTTPPort))

	<-stopChan // wait for SIGINT
	logger.Notice("Shutting down server...")

	// refer to https://medium.com/honestbee-tw-engineer/gracefully-shutdown-in-go-http-server-5f5e6b83da5a
	// shut down gracefully, but wait no longer than 10 seconds before halting
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		app.Close()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("Server shutdown error: %+v", err)
	}
	logger.Notice("Server gracefully stopped")
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	target := "https://" + r.Host + r.URL.Path
	if len(r.URL.RawQuery) > 0 {
		target += "?" + r.URL.RawQuery
	}
	// consider HSTS if your clients are browsers
	w.Header().Set("Connection", "close")
	http.Redirect(w, r, target, 302)
}

func stlAge(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// add max-age to get A+
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
