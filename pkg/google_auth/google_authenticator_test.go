package google_auth

import (
	"github.com/skip2/go-qrcode"
	"testing"
)

var Secret = Base32NoPaddingEncodedSecretFunc("")

func TestTotp(t *testing.T) {
	g := GoogleAuthenticator2FaSha1{
		Base32NoPaddingEncodedSecret: Secret,
		ExpireSecond:                 30,
		Digits:                       6,
	}
	totp, err := g.Totp()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(totp)
}

func TestQr(t *testing.T) {
	g := GoogleAuthenticator2FaSha1{
		Base32NoPaddingEncodedSecret: Secret,
		ExpireSecond:                 30,
		Digits:                       6,
	}
	qrString := g.QrString("SSH2A:github.com/IUnlimit/ssh2a", "IllTamer")
	t.Log(qrString)
	code, err := qrcode.New(qrString, qrcode.Medium)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(code.ToSmallString(false))
}
