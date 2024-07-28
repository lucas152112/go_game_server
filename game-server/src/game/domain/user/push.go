package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"game/domain/hall"
	"game/pb"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/glog"
)

type postCreateUserReq struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type postCreateUserRes struct {
	Entities []entitiesInfo `json:"entities"`
}

type entitiesInfo struct {
	UUID     string `json:"uuid"`
	Type     string `json:"type"`
	UserName string `json:"username"`
	Nickname string `json:"nickname"`
}

//CreatePushUser ...
func CreatePushUser(userID, userName, password, nickname string) (bool, string) {
	glog.Info("CreatePushUser in userID:", userID)
	url := "http://a61.easemob.com/1176220526118195/alive/users"

	reqData := &postCreateUserReq{}

	reqData.UserName = userName
	reqData.Password = password
	reqData.Nickname = nickname

	jsonStr, errT := json.Marshal(reqData)
	if errT != nil {
		glog.Info("Marshal error:", errT)
		return false, ""
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer YWMtRrTbxOWBEeyx5PWYWoe3NAMFDHoa1TGRlcHYk-ka9okiJ1MfYElNeIqyjWY_7u5XAgMAAAGBOIEbIDeeSAAbkkRbWxCvC-IImC659YqMyEK9KHtvbtX8DDmW5G-oCw")

	client := &http.Client{}
	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		glog.Info("do request error:", err)
		return false, ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		glog.Info("status:", resp.Status)
		return false, ""
	}

	body, okBody := ioutil.ReadAll(resp.Body)
	if okBody != nil {
		return false, ""
	}
	glog.Info("response Body:", string(body))
	res := &postCreateUserRes{}
	errUn := json.Unmarshal([]byte(body), res)
	if errUn != nil {
		fmt.Println("response Unmarshal error:", errUn)
		return false, ""
	}

	glog.Info("response res:", res)
	if len(res.Entities) > 0 {
		notify := &pb.PushUserInfoNotify{}
		notify.UUID = ""
		notify.Nickname = ""
		notify.Username = userName
		notify.Password = password
		GetPlayerManager().SendClientMsgNew(userID, int32(pb.MessageIDPushUserInfoNotify), notify)
		return true, res.Entities[0].UUID
	}

	return false, ""
}

type pushReq struct {
	Targets     []string        `json:"targets"`
	PushMessage pushMessageInfo `json:"pushMessage"`
	Async       bool            `json:"async"`
	Strategy    int             `json:"strategy"`
}

type pushMessageInfo struct {
	Titile   string `json:"title"`
	SubTitle string `json:"subTitle"`
	Content  string `json:"content"`
}

type pushRes struct {
	ID         string `json:"id"`
	PushStatus string `json:"pushStatus"`
	Desc       string `json:"desc"`
}

//SendPushToOneUser ...
func SendPushToOneUser(userID string, nickname string) bool {
	glog.Info("SendPushToOneUser userID:", userID)
	url := "http://a61.easemob.com/1176220526118195/alive/push/single"
	u, errU := FindByUserId(userID)
	if errU != nil {
		glog.Info("SendPushToOneUser not find userID:", userID)
		return false
	}

	if u.PushUserName == "" {
		glog.Info("SendPushToOneUser pushUserName = null userID:", userID)
		return false
	}

	var content string

	if u.Language == 1 {
		content = fmt.Sprintf("你关注的主播\"%v\"开播了，快去观看吧!", nickname)
	} else if u.Language == 3 {
		content = "วีเจที่คุณติดตามเปิดไลฟ์สดแล้ว  ไปดูได้แล้ว"
	} else if u.Language == 0 {
		content = "Live streamer in your follow now on air, go to watch it."
	} else {
		content = fmt.Sprintf("你關注的主播\"%v\"開播了，快去觀看吧!", nickname)
	}
	reqData := &pushReq{}
	reqData.Targets = append(reqData.Targets, u.PushUserName)
	reqData.PushMessage.Content = content
	reqData.PushMessage.Titile = "AALIVE"
	reqData.Async = false
	reqData.Strategy = 3

	jsonStr, errT := json.Marshal(reqData)
	if errT != nil {
		glog.Info("Marshal error:", errT)
		return false
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer YWMtRrTbxOWBEeyx5PWYWoe3NAMFDHoa1TGRlcHYk-ka9okiJ1MfYElNeIqyjWY_7u5XAgMAAAGBOIEbIDeeSAAbkkRbWxCvC-IImC659YqMyEK9KHtvbtX8DDmW5G-oCw")

	client := &http.Client{}
	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		glog.Info("do request error:", err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		glog.Info("status:", resp.Status)
		return false
	}

	body, okBody := ioutil.ReadAll(resp.Body)
	if okBody != nil {
		glog.Info("okBody = nil:", resp.Status)
		return false
	}
	glog.Info("response Body:", string(body))
	res := &pushRes{}
	errUn := json.Unmarshal([]byte(body), res)
	if errUn != nil {
		glog.Info("response Unmarshal error:", errUn)
		return false
	}

	glog.Info("response res:", res)

	return true
}

//PushLable ...
func PushLable(liverID string) {
	if !getPushLable(liverID, "english") {
		createPushLable(liverID, "englist")
	}
}

type getPushLableData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Count       int    `json:"count"`
	CreatedAt   int64  `json:"createdAt"`
}
type getPushLableRes struct {
	Data getPushLableData `json:"data"`
}

func getPushLable(liverID string, language string) bool {
	url := "http://a61.easemob.com/1176220526118195/alive/push/lable/"
	url += language
	url += liverID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("getPushLable NewRequest error:", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer YWMtRrTbxOWBEeyx5PWYWoe3NAMFDHoa1TGRlcHYk-ka9okiJ1MfYElNeIqyjWY_7u5XAgMAAAGBOIEbIDeeSAAbkkRbWxCvC-IImC659YqMyEK9KHtvbtX8DDmW5G-oCw")

	client := &http.Client{}
	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("getPushLable do error:", err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("getPushLable status:", resp.Status)
		return false
	}

	body, okBody := ioutil.ReadAll(resp.Body)
	if okBody != nil {
		fmt.Println("getPushLable ReadAll:", okBody)
		return false
	}
	fmt.Println("getPushLable Body:", string(body))
	res := &getPushLableRes{}
	errUn := json.Unmarshal([]byte(body), res)
	if errUn != nil {
		fmt.Println("getPushLable Unmarshal error:", errUn)
		return false
	}

	fmt.Println("getPushLable res:", res)

	if res.Data.Name != "" {
		return true
	}

	return false
}

type createPushLableReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func createPushLable(liverID string, language string) bool {
	url := "http://a61.easemob.com/1176220526118195/alive/push/lable/"
	reqData := &createPushLableReq{}
	name := language
	name += liverID
	reqData.Name = name

	jsonStr, errT := json.Marshal(reqData)
	if errT != nil {
		fmt.Println("createPushLable Marshal error:", errT)
		return false
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("do createPushLable get error:", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer YWMtRrTbxOWBEeyx5PWYWoe3NAMFDHoa1TGRlcHYk-ka9okiJ1MfYElNeIqyjWY_7u5XAgMAAAGBOIEbIDeeSAAbkkRbWxCvC-IImC659YqMyEK9KHtvbtX8DDmW5G-oCw")

	client := &http.Client{}
	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("createPushLable do request error:", err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("createPushLable status:", resp.Status)
		return false
	}

	body, okBody := ioutil.ReadAll(resp.Body)
	if okBody != nil {
		fmt.Println("createPushLable ReadAll:", okBody)
		return false
	}
	fmt.Println("createPushLable Body:", string(body))
	res := &getPushLableRes{}
	errUn := json.Unmarshal([]byte(body), res)
	if errUn != nil {
		fmt.Println("createPushLable Unmarshal error:", errUn)
		return false
	}

	fmt.Println("createPushLable res:", res)

	if res.Data.Name != "" {
		return true
	}

	return false
}

//LiverOnLineSendPush ...
func LiverOnLineSendPush(liverID string, nickname string) {
	time.Sleep(time.Minute)
	if !GetLiverManager().LiverIsOnline(liverID) {
		return
	}
	GetLiver(liverID)
	for i := 0; i < 100; i++ {
		userList := hall.GetFollowMe(100, i, liverID)
		if len(userList) == 0 {
			break
		}

		for j := 0; j < len(userList); j++ {
			userID := userList[j]
			SendPushToOneUser(userID, nickname)
		}
	}
}
