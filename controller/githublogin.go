package controller

import (
	"context"
	"fmt"
	"goyoubbs/model"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/rs/xid"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

var (
	oauthStateString = "pseudo-random"
)

// GithubOauth github验证
func githubOauth(clientID, clientSecret string) *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"user", "email"},
		// RedirectURL:  "https://flarum.yjzq.fun/auth/github/callback",
		Endpoint: githuboauth.Endpoint,
	}
	return conf

}

// GithubOauthHandler github用户登录
func GithubOauthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h
	gOauth := githubOauth(h.App.Cf.Site.GithubClientID, h.App.Cf.Site.GithubClientSecret)
	url := gOauth.AuthCodeURL(oauthStateString)

	// http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	// if r.FormValue("state") != "" && r.FormValue("code") != "" {
	// 	gOauth := githubOauth(h.App.Cf.Site.GithubClientID, h.App.Cf.Site.GithubClientSecret)
	// 	data, err := getUserInfo(gOauth, r.FormValue("state"), r.FormValue("code"))
	// 	if data == nil || err != nil {
	// 		fmt.Println(err.Error())
	// 		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 		return
	// 	}
	// 	fmt.Println(data)
	// } else {
	// https: //github.com/login/oauth/authorize?scope=user%3Aemail&
	// state=ac79b976fe05010f6e569717a14e75f9
	// redirect_uri=https%3A%2F%2Fdiscuss.flarum.org.cn%2Fauth%2Fgithub&client_id=810332e645869bef14af&display=popup
	fmt.Println(url)
	url += "&display=popup&response_type=code&approval_prompt=auto"
	w.Header()["Content-Type"] = []string{"text/html; charset=UTF-8"}
	w.Header()["x-csrf-token"] = []string{"Xz2zpiS2RYANX8Va05aPTB85C7Qo53mVe7QpX9Xg"}
	w.Header().Del("Content-Length")
	fmt.Println(w.Header())

	http.Redirect(w, r, url, 302)
	return
	// }
}

// GithubOauthCallbackHandler github用户登录回调
func GithubOauthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h
	gOauth := githubOauth(h.App.Cf.Site.GithubClientID, h.App.Cf.Site.GithubClientSecret)
	data, err := getUserInfo(gOauth, r.FormValue("state"), r.FormValue("code"))
	if data == nil || err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sqlDB := h.App.MySQLdb
	// redisDB := h.App.RedisDB
	uobj, err := model.SQLUserGetByName(sqlDB, *data.Name)
	sessionid := xid.New().String()
	uobj.Session = sessionid

	// uobj.CachedToRedis(redisDB)
	// h.SetCookie(w, "SessionID", strconv.FormatUint(uobj.ID, 10)+":"+sessionid, 365)
	// uobj.RefreshCSRF(redisDB)

	rsp := response{}
	rsp.Retcode = 200
	rsp.Retmsg = "登录成功"
	fmt.Println(data)
	w.WriteHeader(200)
	h.jsonify(w, rsp)

	// http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

func getUserInfo(_oauth *oauth2.Config, state string, code string) (*github.User, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := _oauth.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	oauthClient := _oauth.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)

	ctx := context.Background()
	user, _, err := client.Users.Get(ctx, "")

	return user, nil
}
