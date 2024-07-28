package user

import (
	"fmt"
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sort"
	"time"
	"game/util"
)

type UserWrapper struct {
	users []*User
	by    func(p, q *User) bool
}
type SortBy func(p, q *User) bool

func (uw UserWrapper) Len() int {
	return len(uw.users)
}
func (uw UserWrapper) Swap(i, j int) {
	uw.users[i], uw.users[j] = uw.users[j], uw.users[i]
}
func (uw UserWrapper) Less(i, j int) bool {
	return uw.by(uw.users[i], uw.users[j])
}

func SortUsers(users []*User, by SortBy) {
	sort.Sort(UserWrapper{users, by})
}

func GetUsersSortBalance(maxSortNum int) ([]*User, error) {
	users := []*User{}
	for i := 1; i <= 20; i++ {
		userTableC := "user_" + fmt.Sprintf("%v", i)
		tempUsers := []*User{}
		util.WithUserCollection(userTableC, func(c *mgo.Collection) error {
			return c.Find(bson.M{"balance": bson.M{"$gt": 0}}).Select(bson.M{"userId": 1, "nickname": 1, "balance": 1, "winGames": 1, "totalGames": 1, "createTime": 1}).Sort("-balance").Limit(maxSortNum).All(&tempUsers)
		})
		if len(tempUsers) != 0 {
			users = append(users, tempUsers...)
		}
	}
	SortUsers(users, func(p, q *User) bool {
		return p.Balance > q.Balance
	})
	if len(users) > maxSortNum {
		users = users[:maxSortNum]
	}
	return users, nil
}

func GetUsersSortTotalGames(maxSortNum int) ([]*User, error) {
	users := []*User{}
	for i := 1; i <= 20; i++ {
		userTableC := "user_" + fmt.Sprintf("%v", i)
		tempUsers := []*User{}
		util.WithUserCollection(userTableC, func(c *mgo.Collection) error {
			return c.Find(bson.M{"totalGames": bson.M{"$gt": 0}}).Select(bson.M{"userId": 1, "nickname": 1, "balance": 1, "winGames": 1, "totalGames": 1, "createTime": 1}).Sort("-totalGames").Limit(maxSortNum).All(&tempUsers)
		})
		if len(tempUsers) != 0 {
			users = append(users, tempUsers...)
		}
	}
	SortUsers(users, func(p, q *User) bool {
		return p.TotalGames > q.TotalGames
	})
	if len(users) > maxSortNum {
		users = users[:maxSortNum]
	}
	return users, nil
}

//func GetUsersSortJQK(maxSortNum int) ([]*User, error) {
//	users := []*User{}
//	for i := 1; i <= 20; i++ {
//		userTableC := "user_" + fmt.Sprintf("%v", i)
//		tempUsers := []*User{}
//		util.WithUserCollection(userTableC, func(c *mgo.Collection) error {
//			return c.Find(bson.M{"jQK": bson.M{"$gt": 0}}).Select(bson.M{"userId": 1, "nickname": 1, "jQK": 1, "jQKFrozen": 1}).Sort("-jQK").Limit(maxSortNum).All(&tempUsers)
//		})
//		if len(tempUsers) != 0 {
//			users = append(users, tempUsers...)
//		}
//	}
//	SortUsers(users, func(p, q *User) bool {
//		return p.JQK > q.JQK
//	})
//	if len(users) > maxSortNum {
//		users = users[:maxSortNum]
//	}
//	return users, nil
//}

func GetTotalBalance() (float32, float32, float32, float32, error) {
	var balance float32 = 0
	var frozenBalance float32 = 0
	var jQK float32 = 0
	var jQKFrozen float32 = 0
	p := []bson.M{
		bson.M{"$group": bson.M{"_id": 0, "balance": bson.M{"$sum": "$balance"}, "frozenBalance": bson.M{"$sum": "$frozenBalance"}, "coins": bson.M{"$sum": "$coins"}, "jQK": bson.M{"$sum": "$jQK"}, "jQKFrozen": bson.M{"$sum": "$jQKFrozen"}}},
	}
	for i := 1; i <= 20; i++ {
		userTableC := "user_" + fmt.Sprintf("%v", i)
		resp := []bson.M{}
		err := util.WithUserCollection(userTableC, func(c *mgo.Collection) error {
			return c.Pipe(p).All(&resp)
		})
		glog.Info("GetTotalBalance err:", err, resp)
		if err == nil && len(resp) > 0 {
			balance += float32(resp[0]["balance"].(float64))
			frozenBalance += float32(resp[0]["frozenBalance"].(float64))
			jQK += float32(resp[0]["jQK"].(float64))
			jQKFrozen += float32(resp[0]["jQKFrozen"].(float64))
		}
	}
	return balance, frozenBalance, jQK, jQKFrozen, nil
}

//用户JQK持有信息
type User_JQK_Info struct {
	UserId         string    `bson:"userId"`
	JQK            float32   `bson:"jQK"`
	JQKAge         float32   `bson:"jQKAge"`         //币龄
	JQKLastChgTime time.Time `bson:"jQKLastChgTime"` //权益币最近一次变化时间
}

//所有用户JQK持有信息
func GetUsersJQKInfo() ([]*User_JQK_Info, error) {
	items := []*User_JQK_Info{}
	for i := 1; i <= 20; i++ {
		tempItems := []*User_JQK_Info{}
		userTableC := "user_" + fmt.Sprintf("%v", i)
		err := util.WithUserCollection(userTableC, func(c *mgo.Collection) error {
			return c.Find(bson.M{"jQK": bson.M{"$gt": 0}}).Select(bson.M{"userId": 1, "jQK": 1, "jQKAge": 1, "jQKLastChgTime": 1}).All(&tempItems)
		})
		if err == nil {
			items = append(items, tempItems...)
		}
	}
	return items, nil
}
