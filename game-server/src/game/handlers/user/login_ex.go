package user

import (
	"game/config"
	"crypto/md5"
	"fmt"
	"io"
)

const (
	retryTimes       = 5
	retryMillisecond = 500
	AppID            = "100026251"
	IOS_CHANNEL      = "178"
	CAOHUA_CHANNEL   = "184"
	CaoHuaAppId      = "144"
	CaoHuaAppKey     = "AB2CBFCF1FE60C5900504D9FECA2A868"
	LESHI_CHANNEL    = "173"
	LeshiAppId       = "500098"
	LeshiAppKey      = "d4d7855c08de4ceead84ff28e132f1c3"

	XUNLEI_CHANNEL = "186"
	XunleiAppId    = "050283"
	XunleiAppKey   = "zdJKkuCNsd5luaurCSIJTi2e52nunMPA"

	HAIMA_IOS_CHANNEL = "187"
	HaimaIosAppId     = "a022b0bda1fbba6b65f8ca2d291e2f1b"
	HaimaIosAppKey    = "017e2413864f7c4e807afee9f2435dc0"

	JINLI_CHANNEL  = "212"
	jinliApikey    = "918AD85BD7214741B04B0E2ABA1AD3AD"
	jinliSecretKey = "437D75979E9F4F75882453CBEA67BA52"
	jinliHost      = "id.gionee.com"
	jinliPort      = "443"
	jinliMethod    = "POST"
	jinliUrl       = "/account/verify.do"

	MEIZU_CHANNEL  = "213"
	MeizuApiID     = "2906651"
	MeizuAppSecret = "xIF0GIyisU3E1j3GeyKDNI864friYVQ6"
	MeizuHost      = "https://api.game.meizu.com/game/security/checksession"

	LIANXIANG_CHANNEL = "214"

	KUPAI_CHANNEL  = "215"
	KupaiApiID     = "5000002955"
	KupaiSecretKey = "52886cb546764689b87de1ad747cc6a3"
	KupaiHost      = "https://openapi.coolyun.com/oauth2/token"

	IOS51_CHANNEL = "216"
	IOS51_APPID   = "100001045"
	IOS51Host     = "http://api.51pgzs.com/passport/checkLogin.php"

	LENOVO_CHANNEL = "222"
	LenovoHost     = "http://passport.lenovo.com/interserver/authen/1.2/getaccountid?"
	LenovoAppId    = "1602181616922.app.ln"

	PAPA_CHANNEL  = "223"
	PapaHost      = "http://sdkapi.papa91.com/auth/check_token"
	PapaAppkey    = "16000050"
	PapaSecretkey = "e2b64359bcc6d810128cf6267d7befb39434d4eb0afae75a7ffa03d56511f425"

	NDUO_CHANNEL        = "225"
	HTC_CHANNEL         = "226"
	ANDROID_360_CHANNEL = "227"
	ANDROID_360_HOST    = "https://openapi.360.cn/user/me"

	UC9GAME_CHANNEL = "218"
	UC_HOSTTEST     = "http://sdk.test4.g.uc.cn/cp/account.verifySession"
	UC_HOST         = "http://sdk.g.uc.cn/cp/account.verifySession"
	UC_APIKEY       = "94ef47f22e5ed23826f65975c5653c1a"

	UUCUN_CHANNEL = "229"
	UUCUN_HOST    = "http://uavapi.uuserv20.com/checkAccessToken.do"

	LE8_CHANNEL = "230"
	LE8_APPID   = "086293f3974e9deb3710653af5da72b1"
	LE8_HOST    = "http://api.le890.com/index.php?m=api&a=validate_token"

	SKY_CHANNEL = "235"
	SKY_APPID   = "086293f3974e9deb3710653af5da72b1"
	SKY_HOST    = "http://111.1.17.152:10015/skyppa/index!check.action"

	YIYOU_CHANNEL = "236"
	YIYOU_HOST    = "http://uavapi.uuserv20.com/checkAccessToken.do"

	BAZHANG_CHANNEL = "237"
	BAZHANG_APPKEY  = "ac63001222bb02f316aa3f316fabaa94"

	SHUOWAN_CHANNEL = "246"
	SHUOWAN_APPKEY  = "bf886f12e2ff4a72b2d871da3260f8f1"

	WANDOUJIA_CHANNEL = "245"
	WANDOUJIA_APPID   = "100040378"

	QISHIZHUSHOU_CHANNEL = "247"
	QISHIZHUSHOU_APPKEY  = "d87d6fbc5b4612dc3e884c61d258dd825b4612dc3e884c61"
)

func genPassword(pwd string) string {
	h := md5.New()
	io.WriteString(h, config.ControlKey+pwd)
	return fmt.Sprintf("%x", h.Sum(nil))
}
