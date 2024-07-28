package user

import (
	"errors"
	"fmt"
	"game/domain/core"
	"game/util"
	"strconv"

	"game/domain/hall"

	"github.com/golang/glog"
)

func GetUserPayInfoLog(userID string) (error, *core.UserPayInfoLog) {
	db := core.DB

	var log core.UserPayInfoLog
	if db.Where("user_id = ?", userID).First(&log).RecordNotFound() {
		log.UserID = userID
		log.Amount = 0
		log.PayTimes = 0
		return errors.New("user not pay"), &log
	}

	return nil, &log
}

func SaveUserPayInfoLog(log *core.UserPayInfoLog) {
	db := core.DB
	go db.Save(log)
	return
}

func GetTodayUsersPay() int {
	var total_money []int
	todayTime := util.GetTodyZeroStr()
	sqlStr := fmt.Sprintf("SELECT SUM(amount) AS total_money FROM user_pay_logs WHERE time_str >= '%s'", todayTime)
	fmt.Println("sqlstr :", sqlStr)
	err := core.DB.Raw(sqlStr).Pluck("SUM(amout) as total_money", &total_money).Error
	if err != nil {
		glog.Error("没有数据 %v", err)
		return 0
	}

	glog.Info("day:", todayTime, ",money:", total_money)
	if len(total_money) > 0 {
		return total_money[0]
	}

	return 0
}

//SaveUserPayLog ...
func SaveUserPayLog(userID string, amount int, data string, threeOrder string, productID int, orderID string) {
	u, errU := FindByUserId(userID)
	if errU == nil {
		db := core.DB
		var log core.UserPayLog
		log.UserID = userID
		log.Amount = amount
		if u != nil {
			log.Channel = u.ChannelId
		}
		log.PurchaseToken = data
		log.ThreeOrder = threeOrder
		log.ProductID = productID
		log.OrderID = orderID

		log.TimeStr = util.GetDbTime()
		go db.Save(&log)

		if u.InviteUserId != "" {
			inviteUserID := u.InviteUserId
			userID := u.UserId
			nickname := u.Nickname
			rate := 0.05
			userData := hall.GetUserData(inviteUserID)
			if userData.InviteRewardRate != 0 {
				rate = userData.InviteRewardRate
			}
			fenhong := float64(amount) * rate
			value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", fenhong), 64)

			glog.Info("pay fenhong userID:", userID, ",invitor:", inviteUserID, ",payAmout:", amount, ", fen:", value)

			go SaveInviteLog(userID, inviteUserID, nickname, 2, float32(value), 3, 1, amount)
		}
		return
	}
}

//GetUserPayLog ...
func GetUserPayLog(threeOrderID string) bool {
	db := core.DB

	var log core.UserPayLog
	if db.Where("three_order = ?", threeOrderID).First(&log).RecordNotFound() {
		return false
	}

	return true
}

//NoPayUserToday ...
func NoPayUserToday(userID string) bool {
	db := core.DB

	var log core.UserPayLog
	if db.Where("user_id = ? and time_str >= ?", userID, util.GetTodyZeroStr()).First(&log).RecordNotFound() {
		return true
	}

	return false
}

//GetUsersPayLog ...
func GetUsersPayLog(limit int, offSet int) []*core.UserPayLog {
	db := core.DB
	var logs []*core.UserPayLog

	db.Order("id desc").Limit(limit).Offset(offSet * limit).Find(&logs)
	return logs
}

//GetUsersPayLogByTime ...
func GetUsersPayLogByTime(begin, end string, limit int, page int) (int, []*core.UserPayLog) {
	db := core.DB
	var logs []*core.UserPayLog

	var count int
	db.Model(&core.UserPayLog{}).Where("created_at > ? and created_at < ?", begin, end).Count(&count)

	if count > 0 {
		db.Order("id desc").Where("created_at > ? and created_at < ?", begin, end).Limit(limit).Offset(limit * page).Find(&logs)
	}

	return count, logs
}
