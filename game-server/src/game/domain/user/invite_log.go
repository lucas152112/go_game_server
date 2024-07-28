package user

import (
	"game/domain/core"
	"game/util"
)

//SaveInviteLog ...
func SaveInviteLog(userID, inviteUserID, nickname string, reason int, reward float32, rewardType int, status int, payAmount int) error {
	log := core.InviteLog{}
	log.InvitorID = inviteUserID
	log.UserID = userID
	log.Nickname = nickname
	log.Reward = reward
	log.RewardType = rewardType
	log.Status = status
	log.Reason = reason
	log.TimeStr = util.GetDbTime()
	log.PayAmount = payAmount
	DB.Save(&log)
	return nil
}

//GetInviteLog ...
func GetInviteLog(userID string, limit int, offSet int) []*core.InviteLog {
	db := core.DB
	var logs []*core.InviteLog

	db.Order("id").Where("invitor_id = ?", userID).Limit(limit).Offset(offSet * limit).Find(&logs)
	return logs
}

//GetAllInviteLog ...
func GetAllInviteLog(limit int, offSet int) []*core.InviteLog {
	db := core.DB
	var logs []*core.InviteLog

	db.Order("id desc").Limit(limit).Offset(offSet * limit).Find(&logs)
	return logs
}
