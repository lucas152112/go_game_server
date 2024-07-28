package user

import (
	"encoding/json"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
)

const WEIXIN_APP_ID = "wx05e5e26a69f5d78e"
const WEIXIN_SECRET = "bdb862bfd62586c859492c7a3f913309"

const WEIXIN_APP_ID_APPLE = "wxb2b27e362ee47973"
const WEIXIN_SECRET_APPLE = "e5f3b7589af0c80d7dba727950f099b2"

type ACCESS_TOKEN_RESULT struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionId      string `json:"unionid"`
}

func GetAccessToken(code string, channel string) (bool, *ACCESS_TOKEN_RESULT) {
	tempAppId := ""
	tempSecret := ""

	tempAppId = WEIXIN_APP_ID
	tempSecret = WEIXIN_SECRET

	if channel == "apple" {
		tempAppId = WEIXIN_APP_ID_APPLE
		tempSecret = WEIXIN_SECRET_APPLE
	}

	sendStr := "https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + tempAppId
	sendStr = sendStr + "&secret=" + tempSecret + "&code=" + code + "&grant_type=authorization_code"
	res, err := http.Get(sendStr)
	if err != nil {
		glog.Info("GetAccessToken err:", err)
		return false, nil
	}

	resbody, errRes := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if errRes != nil {
		glog.Info("GetAccessToken readall err:", err)
		return false, nil
	}

	result := &ACCESS_TOKEN_RESULT{}
	err = json.Unmarshal(resbody, result)
	if err != nil {
		glog.Info("GetAccessToken Unmarshal err:", err)
		return false, nil
	}

	return true, result
}

func FreshAccessToken(freshToken string, channel string) (bool, *ACCESS_TOKEN_RESULT) {
	tempAppId := ""
	tempAppId = WEIXIN_APP_ID

	if channel == "apple" {
		tempAppId = WEIXIN_APP_ID_APPLE
	}

	sendStr := "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=" + tempAppId
	sendStr = sendStr + "&grant_type=refresh_token&refresh_token=" + freshToken
	res, err := http.Get(sendStr)
	if err != nil {
		glog.Info("FreshAccessToken err:", err)
		return false, nil
	}

	resbody, errRes := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if errRes != nil {
		glog.Info("FreshAccessToken readall err:", err)
		return false, nil
	}

	result := &ACCESS_TOKEN_RESULT{}
	err = json.Unmarshal(resbody, result)
	if err != nil {
		glog.Info("FreshAccessToken Unmarshal err:", err)
		return false, nil
	}

	return true, result
}

type WEIXIN_USER_INFO struct {
	OpenId     string `json:"openid"`
	Nickname   string `json:"nickname"`
	Sex        int    `json:"sex"` //1男，2女
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	HeadImgUrl string `json:"headimgurl"`
	UnionId    string `json:"unionid"`
}

func GetUserInfo(openId string, accessToken string) (bool, *WEIXIN_USER_INFO) {
	sendStr := "https://api.weixin.qq.com/sns/userinfo?access_token=" + accessToken + "&openid=" + openId
	res, err := http.Get(sendStr)
	if err != nil {
		glog.Info("GetUserInfo err:", err)
		return false, nil
	}

	resbody, errRes := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if errRes != nil {
		glog.Info("GetUserInfo readall err:", err)
		return false, nil
	}

	result := &WEIXIN_USER_INFO{}
	err = json.Unmarshal(resbody, result)
	if err != nil {
		glog.Info("GetUserInfo Unmarshal err:", err)
		return false, nil
	}

	return true, result
}
