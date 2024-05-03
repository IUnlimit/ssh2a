package tools

import (
	"encoding/base32"
	"testing"
)

// otpauth://totp/ACME%20Co:john@example.com?secret=HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ&issuer=ACME%20Co&algorithm=SHA1&digits=6&period=30
const testSecret = "HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ" //base32-no-padding-encoded-string

func TestTotp(t *testing.T) {
	// 	// 要编码的数据
	//	data := []byte("Hello, World!")
	//
	//	// 创建一个 base32 编码器，没有填充
	//	encoder := base32.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567").WithPadding(base32.NoPadding)
	//
	//	// 编码数据
	//	encodedData := encoder.EncodeToString(data)
	g := GoogleAuthenticator2FaSha1{
		Base32NoPaddingEncodedSecret: testSecret,
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
		Base32NoPaddingEncodedSecret: testSecret,
		ExpireSecond:                 30,
		Digits:                       6,
	}
	qrString := g.QrString("TechBlog:mojotv.cn", "Eric Zhou")
	t.Log(qrString)
}

func TestMakePhoneOrEmail2FACode(t *testing.T) {
	yourAppSecret := "server_config_secret"
	phoneOrEmailOrUserId := "13312345678" //neochau@gmail.com
	//base 32 no padding encode 是必须的不能省去
	googleAuthenticatorKey := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString([]byte(yourAppSecret + phoneOrEmailOrUserId))
	mfa := GoogleAuthenticator2FaSha1{
		Base32NoPaddingEncodedSecret: googleAuthenticatorKey,
		ExpireSecond:                 300, //五分钟之后失效
		Digits:                       4,   //数字
	}
	code, err := mfa.Totp()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("使用这个code发送给用户邮箱手机,", "也可以使用这个code来,验证码收到的code == 服务器收到的code 是否一直", code)
}
