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
		Endpoint:     githuboauth.Endpoint,
	}
	return conf

}

// GithubOauthHandler github用户登录
func GithubOauthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h
	gOauth := githubOauth(h.App.Cf.Site.GithubClientID, h.App.Cf.Site.GithubClientSecret)
	url := gOauth.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusFound)
	return
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
	redisDB := h.App.RedisDB

	uobj, err := model.SQLUserGetByEmail(sqlDB, data.GetEmail())

	fmt.Println(uobj)
	if !uobj.IsValid() {
		uobj, err = model.SQLGithubRegister(sqlDB, data)
	}

	sessionid := xid.New().String()
	uobj.Session = sessionid

	uobj.CachedToRedis(redisDB)
	h.SetCookie(w, "SessionID", uobj.StrID()+":"+sessionid, 365)
	var retData string
	if uobj.IsValid() {
		retData = `<script>window.close(); window.opener.app.authenticationComplete({"loggedIn":true});</script>`

	} else {
		retData = `<script>window.close(); window.opener.app.authenticationComplete({"loggedIn":false});</script>`
	}
	w.Write([]byte(retData))
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