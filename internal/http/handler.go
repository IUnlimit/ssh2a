package http

import (
	"github.com/IUnlimit/ssh2a/internal/cache"
	"github.com/IUnlimit/ssh2a/internal/model"
	"github.com/IUnlimit/ssh2a/tools"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func handleAuth(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "auth.tmpl", gin.H{})
}

func handleAPIAuth(ctx *gin.Context) {
	var token model.AccessToken
	ip := ctx.RemoteIP()
	if err := ctx.BindJSON(&token); err != nil {
		log.Warnf("Remote IP %s requests auth api with error: %v", ip, err)
		return
	}

	g := tools.GoogleAuthenticator2FaSha1{
		Base32NoPaddingEncodedSecret: testSecret,
		ExpireSecond:                 30,
		Digits:                       6,
	}
	totp, err := g.Totp()
	if err != nil {
		t.Error(err)
		return
	}
	if token.Token != "765743073" {
		ctx.JSON(http.StatusOK, model.Response{StatusCode: -1, StatusMsg: "Invalid token. Please try again."})
		return
	}
	cache.SetAccessCache(ip)
	cache.InvalidTempCache(ip)
	ctx.JSON(http.StatusOK, model.Response{StatusCode: 0})
}
