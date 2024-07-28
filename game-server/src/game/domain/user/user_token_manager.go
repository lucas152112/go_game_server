package user

import (
	"crypto/md5"
	"fmt"
	"github.com/golang/glog"
	"io"
	"sync"
	"time"
)

type UserTokenInfo struct {
	UserId string
	Time   int64
}

type UserTokenManager struct {
	sync.RWMutex
	tokens map[string]*UserTokenInfo
}

var tokenManager *UserTokenManager

func init() {
	tokenManager = &UserTokenManager{}
	tokenManager.tokens = make(map[string]*UserTokenInfo)
	go tokenManager.checkFunc()
}

func GetUserTokenManager() *UserTokenManager {
	return tokenManager
}

func (this *UserTokenManager) GetToken(userId string) string {
	this.Lock()
	defer this.Unlock()

	tm := time.Now()
	sign := "token" + userId + fmt.Sprintf("%d-%d-%d %02d:%02d:%02d", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
	token := Md5func(sign)

	token = string(token[0:30])

	tokenInfo := &UserTokenInfo{}
	tokenInfo.Time = time.Now().Unix()
	tokenInfo.UserId = userId

	this.tokens[token] = tokenInfo

	return token
}

func (this *UserTokenManager) CheckToken(token string) (bool, string) {
	this.Lock()
	defer this.Unlock()

	info, ok := this.tokens[token]
	if !ok {
		glog.Info("CheckToken filed, token=", token)
		return false, ""
	}

	userId := info.UserId

	delete(this.tokens, token)

	return true, userId
}

func (this *UserTokenManager) check() {
	this.Lock()
	defer this.Unlock()

	t := time.Now().Unix()
	for token, info := range this.tokens {
		if (t - info.Time) > 600 {
			delete(this.tokens, token)
		}
	}
}

func (this *UserTokenManager) checkFunc() {
	for {
		this.check()
		time.Sleep(5 * time.Second)
	}
}

func Md5func(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
