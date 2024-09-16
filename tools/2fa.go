package tools

import (
	"fmt"
	"github.com/IUnlimit/ssh2a/pkg/google_auth"
	log "github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
)

const (
	ExpireSecond = 30
	Digits       = 6
)

func PrintQRCode(privateSecret string) error {
	var Secret = google_auth.Base32NoPaddingEncodedSecretFunc(privateSecret)
	g := google_auth.GoogleAuthenticator2FaSha1{
		Base32NoPaddingEncodedSecret: Secret,
		ExpireSecond:                 ExpireSecond,
		Digits:                       Digits,
	}
	qrString := g.QrString("SSH2A:github.com/IUnlimit/ssh2a", "IllTamer")
	log.Infof("Base32NoPaddingEncodedSecretFunc: %s", qrString)
	code, err := qrcode.New(qrString, qrcode.Medium)
	if err != nil {
		return err
	}
	fmt.Println(code.ToSmallString(false))
	return nil
}

func TotP(privateSecret string) (string, error) {
	var Secret = google_auth.Base32NoPaddingEncodedSecretFunc(privateSecret)
	g := google_auth.GoogleAuthenticator2FaSha1{
		Base32NoPaddingEncodedSecret: Secret,
		ExpireSecond:                 ExpireSecond,
		Digits:                       Digits,
	}
	totp, err := g.Totp()
	if err != nil {
		return "", err
	}
	return totp, err
}
