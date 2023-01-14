package controllers

import (
	"WBA/config"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
	oauthGoogleUrlAPI string
)

type LoginController struct {
}

func NewGoogleLoginController(config *config.Config) (LoginController, error) {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  config.Oauth2["google"]["redirecturl"].(string),
		ClientID:     config.Oauth2["google"]["clientid"].(string),
		ClientSecret: config.Oauth2["google"]["clientsecret"].(string),
		Scopes:       []string{config.Oauth2["google"]["scopes"].(string)},
		Endpoint:     google.Endpoint,
	}
	oauthGoogleUrlAPI = config.Oauth2["google"]["oauthgoogleurlapi"].(string)
	return LoginController{}, nil
}

func (lc *LoginController) GoogleLoginHandler(ctx *gin.Context) {
	state := lc.generateStateOauthCookie(ctx.Writer)
	url := googleOauthConfig.AuthCodeURL(state)
	ctx.Redirect(http.StatusMovedPermanently, url)
}

// cookie 에 일회용 비번을 저장해서 검증
func (lc *LoginController) generateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(1 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, cookie)
	return state
}

func (lc *LoginController) GoogleAuthCallback(ctx *gin.Context) {
	oauthstate, _ := ctx.Cookie("oauthstate")

	// 쿠키값과 폼데이터의 state 가 같은지 비교
	if ctx.Query("state") != oauthstate {
		log.Printf("invalid google oauth state cookie:%s state:%s\n", oauthstate, ctx.Query("state"))
		ctx.Redirect(400, "/")
		return
	}

	email, err := lc.getGoogleUserInfo(ctx.Query("code"), ctx)
	if err != nil {
		log.Println(err.Error())
		ctx.Redirect(400, "/")
		return
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{"isLogined": true, "id": email})
}

func (lc *LoginController) getGoogleUserInfo(code string, ctx *gin.Context) (string, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange %s", err.Error())
	}

	// RefreshToken 은 AccessToken 이 expired 된 경우, 다시 받을때 사용
	resp, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)

	if err != nil {
		return "", fmt.Errorf("failed to get userInfo %s", err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read request Body %s", err.Error())
	}
	var data map[string]interface{} // TopTracks
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}
	return data["email"].(string), nil
}
