package user

import (
	domainUser "game/domain/user"
)

func onCreatePlayer(p *domainUser.GamePlayer, robotWinTimes, robotLoseTimes, robotCurDayEarnGold, robotCurWeekEarnGold int32, robotMaxCards []int) {
	p.NewPlayer = true
}

//func onLogout(userId string) {
//	go domainDZ.GetGameManager().OffLine(userId)
//	go domainDZ.GetMatchManager().OffLine(userId, 0)
//}
