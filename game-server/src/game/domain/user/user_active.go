package user

import (
	"game/domain/core"
	"game/domain/hall"
	"game/util"
)

// SaveUserRegisterLog ...
func SaveUserRegisterLog(userID string, channel string, isTest bool) {
	channel = hall.GetChannelManager().GetChannel(channel)
	var log core.UserRegisterLog
	log.UserID = userID
	log.TimeStr = util.GetDbTime()
	log.Channel = channel
	log.IsTest = isTest

	db := core.DB
	db.Save(&log)

	initUserActiveLog(userID, channel)
	return
}

//initUserActiveLog ...
func initUserActiveLog(userID string, channel string) {
	var log core.UserActivityLog
	log.UserID = userID
	log.AddDay = util.GetDbTime()
	log.Channel = channel
	log.OneDay = 0
	log.ThreeDay = 0
	log.SevenDay = 0
	db := core.DB
	db.Save(&log)
	return
}

//SaveLoginLog ...
func SaveLoginLog(userID string, channel string) {
	var log core.UserLoginLog
	log.UserID = userID
	log.TimeStr = util.GetDbTime()
	log.Channel = channel

	db := core.DB
	db.Save(&log)
	return
}

//SaveChannelActiveCount ...
func SaveChannelActiveCount(channel string) {
	db := core.DB
	timeDay := util.GetCurrentDate2()
	var log core.UserChannelActiveLog
	if db.Where("channel = ? and time_str = ?", channel, timeDay).First(&log).RecordNotFound() {
		log.Channel = channel
		log.TimeStr = timeDay
		log.UserCount = 0
	}

	log.UserCount++
	db.Save(&log)
}

//UpdateUserActiveLog ...
func UpdateUserActiveLog(userID string, createTime string) {
	db := core.DB
	var log core.UserActivityLog
	if db.Where("user_id = ?", userID).First(&log).RecordNotFound() {
		return
	}

	createDay := util.GetYearDayByTimeStr(createTime)
	curDay := util.GetCurYearDay()
	day := curDay - createDay
	if day == 1 {
		log.OneDay = 1
		db.Save(&log)
	} else if day == 3 {
		log.ThreeDay = 1
		db.Save(&log)
	} else if day == 7 {
		log.SevenDay = 1
		db.Save(&log)
	}

	return
}

//GetOneChannelDay ...
func GetOneChannelDay(day string, channel string) (int, int, int, int) {
	row := core.DB.Model(&core.UserActivityLog{}).Where("DATE(created_at) = ? and channel = ?", day, channel).Select("count(*) as new_user, sum(one_day) as one_day, sum(three_day) as three_day, sum(seven_day) as seven_day").Row()
	var newUser int
	var oneDay int
	var threeDay int
	var sevenDay int
	row.Scan(&newUser, &oneDay, &threeDay, &sevenDay)

	return newUser, oneDay, threeDay, sevenDay
}
