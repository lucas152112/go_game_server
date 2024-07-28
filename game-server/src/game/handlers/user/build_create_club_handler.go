package user

import (
	"encoding/json"
	"errors"
	"game/domain/dzclub"
	"game/domain/user"
	"io/ioutil"
	"net/http"
	"game/pb"
)

/**
 创建 创建固定用户
 */

type BuildUserClubReq struct {
	UserId     string     `json:"userId"`
	ClubName   string     `json:"clubName"`
	Icon       string     `json:"icon"`
	Describe   string     `json:"describe"`
	UserCount  int        `json:"userCount"`
}

type BuildUserClubRes struct {
	pb.ProBaseResponse     
	ClubId     int       `json:"clubId"`
	ClubName   string    `json:"clubName"`
}

func BuildUserClub(w http.ResponseWriter, r *http.Request)  {
	req :=&BuildUserClubReq{}
	res :=&BuildUserClubRes{}
	token := r.FormValue("token")
	if token != "majiangwebtoken123" {
		res.Error("token Error",1)
		raw,_:=json.Marshal(res)
		w.Write(raw)
		return
	}

	raw,err:=ioutil.ReadAll(r.Body)
	if err!=nil{
		res.Error("read body error",2)
		raw,_:=json.Marshal(res)
		w.Write(raw)
		return
	}
	err =json.Unmarshal(raw,req)
	if err!=nil{
		res.Error("json parse error",3)
		raw,_:=json.Marshal(res)
		w.Write(raw)
		return
	}

	err,clubId :=CreateClub(req)
	if err!=nil{
		res.Error(err.Error(),4)
		raw,_:=json.Marshal(res)
		w.Write(raw)
		return
	}
	res.ClubId = clubId
	res.ClubName = req.ClubName
	res.Success("Success")
	raw,_ =json.Marshal(res)
	w.Write(raw)
	return
}

type BuildSetJoinClubUserReq struct {
	ClubId     int         `json:"clubId"`
	UserIds    []string    `json:"userIds"`
}

type BuildSetJoinClubUserRes struct {
	pb.ProBaseResponse
	SuccessIds    []string   `json:"successIds"`
	ErrorIds      []string   `json:"errorIds"`
}


func BuildSetJoinClub(w http.ResponseWriter, r *http.Request)  {
	req :=&BuildSetJoinClubUserReq{}
	res :=&BuildSetJoinClubUserRes{}
	token := r.FormValue("token")
	if token != "majiangwebtoken123" {
		res.Error("token Error",1)
		raw,_:=json.Marshal(res)
		w.Write(raw)
		return
	}
	raw,err:=ioutil.ReadAll(r.Body)
	if err!=nil{
		res.Error("read body error",2)
		raw,_:=json.Marshal(res)
		w.Write(raw)
		return
	}
	err =json.Unmarshal(raw,req)
	if err!=nil{
		res.Error("json parse error",3)
		raw,_:=json.Marshal(res)
		w.Write(raw)
		return
	}

	for _,userId := range req.UserIds {
		if EnterCreateClub(req.ClubId,userId){
			res.SuccessIds = append(res.SuccessIds,userId)
		}else{
			res.ErrorIds = append(res.ErrorIds,userId)
		}
	}
	res.Success("Success")
	raw,_ =json.Marshal(res)
	w.Write(raw)
	return

}

func CreateClub( req *BuildUserClubReq ) (error,int) {
	userInfo, err := user.FindByUserId( req.UserId )
	if err!=nil{
		return err,0
	}
	ClubUser :=&pb.DZClubUser{}
	ClubUser.UserId = userInfo.UserId
	ClubUser.Nickname = userInfo.Nickname
	ClubUser.HeadUrl = userInfo.PhotoUrl

	_, list := dzclub.GetUserOwnerClubs(userInfo.UserId)
	if len(list) > 0{
		return nil,0
	}

	errFind,_ :=dzclub.GetClubItemByName(req.ClubName)
	if errFind ==nil{
		return errors.New("名称被占用"),0
	}
	clubId := dzclub.GetClubManager().CreateClub(ClubUser, req.Icon, req.ClubName, req.Describe, "",req.UserCount)

	if clubId != 0{
		return nil ,clubId
	}
	return errors.New("创建失败"),0
}


func EnterCreateClub( clubId int ,userId string  ) bool {
	userInfo, err := user.FindByUserId( userId )
	if err!=nil{
		return false
	}
	ClubUser :=&pb.ClubUser{}
	ClubUser.UserId = userInfo.UserId
	ClubUser.Nickname = userInfo.Nickname
	ClubUser.HeadUrl = userInfo.PhotoUrl

	err =dzclub.GetClubManager().AddToClub(userId,clubId)

	if err==nil {
		return true
	}else{
		return false
	}
}