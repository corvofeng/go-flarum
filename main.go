package main

import (
	"context"
	"crypto/tls"

	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"goyoubbs/cronjob"
	"goyoubbs/getold"
	"goyoubbs/model"

	"goyoubbs/router"
	"goyoubbs/system"
	"goyoubbs/util"

	ct "goyoubbs/controller"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xi2/httpgzip"
	goji "goji.io"
	"goji.io/pat"

	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/http2"
	// "github.com/go-redis/redis/v7"
)

func main() {
	util.InitLogger()
	logger := util.GetLogger()
	configFile := flag.String("config", "config/config.yaml", "full path of config.yaml file")
	getOldSite := flag.String("getoldsite", "0", "get or not old site, 0 or 1, 2")

	flag.Parse()

	c := system.LoadConfig(*configFile)
	app := &system.Application{}

	app.Init(c, os.Args[0])
	model.RankMapInit(app.MySQLdb, app.RedisDB)

	// 验证码信息使用Redis保存
	model.SetCaptchaUseRedisStore(app.RedisDB)

	if *getOldSite == "1" || *getOldSite == "2" {
		bh := &getold.BaseHandler{
			App: app,
		}
		if *getOldSite == "1" {
			bh.GetRemote()
		} else if *getOldSite == "2" {
			bh.GetLocal()
		}
		app.Close()
		return
	}

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
	scf := app.Cf.Site

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

	if mcf.HttpsOn {
		// https
		logger.Debug("Register sll for domain:", mcf.Domain)
		logger.Debug("TLSCrtFile : ", mcf.TLSCrtFile)
		logger.Debug("TLSKeyFile : ", mcf.TLSKeyFile)

		root.Use(stlAge)

		tlsCf := &tls.Config{
			NextProtos: []string{http2.NextProtoTLS, "http/1.1"},
		}

		if mcf.Domain != "" && mcf.TLSCrtFile == "" && mcf.TLSKeyFile == "" {

			domains := strings.Split(mcf.Domain, ",")
			certManager := autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(domains...),
				Cache:      autocert.DirCache("certs"),
				Email:      scf.AdminEmail,
			}
			tlsCf.GetCertificate = certManager.GetCertificate
			//tlsCf.ServerName = domains[0]

			go func() {
				// 必须是 80 端口
				// log.Fatal(http.ListenAndServe(":http", certManager.HTTPHandler(nil)))
			}()

		} else {
			// rewrite
			go func() {
				if err := http.ListenAndServe(":"+strconv.Itoa(mcf.HTTPPort), http.HandlerFunc(redirectHandler)); err != nil {
					logger.Debug("Http2https server failed ", err)
				}
			}()
		}

		srv = &http.Server{
			Addr:           ":" + strconv.Itoa(mcf.HttpsPort),
			Handler:        httpgzip.NewHandler(root, nil),
			TLSConfig:      tlsCf,
			MaxHeaderBytes: int(app.Cf.Site.UploadMaxSizeByte),
		}

		go func() {
			// 如何获取 TLSCrtFile、TLSKeyFile 文件参见 https://www.youbbs.org/t/2169
			logger.Fatal(srv.ListenAndServeTLS(mcf.TLSCrtFile, mcf.TLSKeyFile))
		}()

		logger.Debug("Web server Listen port", mcf.HttpsPort)
		logger.Debug("Web server URL", "https://"+mcf.Domain)

	} else {
		// http
		srv = &http.Server{Addr: ":" + strconv.Itoa(mcf.HTTPPort), Handler: root}
		// srv = &http.Server{Addr: ":" + *httpPort, Handler: root}
		go func() {
			log.Fatal(srv.ListenAndServe())
		}()

		logger.Debug("Web server Listen port", strconv.Itoa(mcf.HTTPPort))
	}

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
