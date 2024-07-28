package user

import (
	"bufio"
	"flag"
	"fmt"
	"game/pb"
	"math/rand"
	mRand "math/rand"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/golang/glog"
)

//RankingBalance ...
const (
	RankingBalance string = "rank_gold"  //财富排行
	RankingShare   string = "rank_share" //分享
	RankingLiver   string = "rank_liver" //主播
	RankingGift    string = "rank_gift"  //打赏
)

//RankingManager ...
//type RankingManager struct {
//	RankingName string
//}

//var rankingManager *RankingManager

var redisAddr string
var nicknameList []string
var defaultBalanceRank []pb.RankingItemDef
var defaultShareRank []pb.RankingItemDef
var defaultLiverRank []pb.RankingItemDef
var defaultGiftRank []pb.RankingItemDef

func init() {
	getnickNameList()
}

//InitRank ...
func InitRank() {
	redisAddr = flag.Lookup("redis").Value.String()
	go runInit()
	initDefaultRank()
}

//GetRankingManager ...
//func GetRankingManager() *RankingManager {
//	return rankingManager
//}

func runInit() {
	for {
		time.Sleep(time.Second * 600)
		if GetRankUpdateTime(RankingShare+"_time") == "" {
			RemoveRank(RankingShare)
			SetRankUpdateTime(RankingShare+"_time", "123")
		}

		if GetRankUpdateTime(RankingLiver+"_time") == "" {
			RemoveRank(RankingLiver)
			SetRankUpdateTime(RankingLiver+"_time", "123")
		}

		if GetRankUpdateTime(RankingGift+"_time") == "" {
			RemoveRank(RankingGift)
			SetRankUpdateTime(RankingGift+"_time", "123")
		}
	}
}

//GetRankingListByType ...
func GetRankingListByType(userID string, rankType int, page int, pageSize int) *pb.RankingListRes {
	res := &pb.RankingListRes{}
	res.RankType = rankType
	res.List = []pb.RankingItemDef{}

	if rankType < 0 || rankType > 3 {
		rankType = 0
	}

	rankTypeStr := RankingBalance
	if rankType == 0 {
		rankTypeStr = RankingBalance
	} else if rankType == 1 {
		rankTypeStr = RankingShare
	} else if rankType == 2 {
		rankTypeStr = RankingGift
	} else {
		rankTypeStr = RankingLiver
	}

	myRank, myScore := GetRankAndScore(rankTypeStr, userID)
	res.MyRank.Rank = myRank
	res.MyRank.Score = myScore
	if myScore == 0 {
		res.MyRank.Rank = 0
	}

	start := page * pageSize
	end := start + pageSize - 1

	if end > 65 {
		return res
	}

	items := GetRankRange(rankTypeStr, start, end)
	for i := 0; i < len(items); i++ {
		item := items[i]

		backItem := pb.RankingItemDef{}
		backItem.Rank = item.Rank
		backItem.UserID = item.UserID
		backItem.Score = item.Score
		backItem.Nickname = ""
		backItem.HeadURL = ""

		u, err := FindByUserId(item.UserID)
		if err == nil && u != nil {
			backItem.Nickname = u.Nickname
			backItem.HeadURL = u.PhotoUrl
		}

		res.List = append(res.List, backItem)
	}

	lt := len(res.List)
	need := 4 - lt

	if page != 0 && lt == 0 {
		need = 0
	}

	if need > 0 {
		for i := 0; i < need; i++ {
			if rankType == 0 {
				if len(defaultBalanceRank) > need {
					itemT := defaultBalanceRank[i]
					itemT.Rank = lt + i + 1 + start
					res.List = append(res.List, itemT)
				}
			} else if rankType == 1 {
				if len(defaultShareRank) > need {
					itemT := defaultShareRank[i]
					itemT.Rank = lt + i + 1 + start
					res.List = append(res.List, itemT)
				}
			} else if rankType == 2 {
				if len(defaultGiftRank) > need {
					itemT := defaultGiftRank[i]
					itemT.Rank = lt + i + 1 + start
					res.List = append(res.List, itemT)
				}
			} else {
				if len(defaultLiverRank) > need {
					itemT := defaultLiverRank[i]
					itemT.Rank = lt + i + 1 + start
					res.List = append(res.List, itemT)
				}

			}
		}
	}

	return res
}

//getnickNameList ...
func getnickNameList() {
	f, e := os.Open("nick.txt")
	if e != nil {
		fmt.Println("File error.")
	} else {
		buf := bufio.NewScanner(f)
		for {
			if !buf.Scan() {
				break
			}
			line := buf.Text()
			line = strings.TrimSpace(line)
			nicknameList = append(nicknameList, line)
		}
	}

	out := []string{}
	temp := nicknameList
	for i := len(temp) - 1; i >= 0; i-- {
		randIndex := mRand.Intn(i + 1)
		t := temp[randIndex]
		out = append(out, t)
		temp[randIndex] = temp[i]
		temp = append(temp[0:], temp[:i]...)
	}

	nicknameList = out

	glog.Info("init robot name list len:", len(nicknameList))
}

func GetOneRandNickname(indexIn int) string {
	nicknameLen := len(nicknameList)
	if indexIn == -1 || indexIn > nicknameLen-1 {
		index := rand.Int() % nicknameLen
		if index < 0 || index >= nicknameLen {
			return ""
		}
		return nicknameList[index]
	}

	return nicknameList[indexIn]
}

func initDefaultRank() {
	count := 10
	nicknameLen := len(nicknameList)
	if nicknameLen == 0 {
		return
	}

	for rankType := 0; rankType < 4; rankType++ {
		begin := 0
		end := 100
		list := []pb.RankingItemDef{}

		if rankType == 0 {
			begin = 10000
			end = 20000
		} else if rankType == 1 {
			begin = 1
			end = 2
		} else if rankType == 2 {
			begin = 1
			end = 2
		} else if rankType == 3 {
			begin = 1
			end = 2
		}

		if rankType != 3 {
			if rankType == 0 {
				count = 13
			}
			for i := 0; i < count && i < nicknameLen; i++ {
				item := pb.RankingItemDef{}
				item.Score = rand.Int()%(end-begin) + begin
				index := rand.Int() % nicknameLen
				if index < 0 || index > nicknameLen {
					return
				}
				item.Nickname = nicknameList[index]

				if rand.Float64() < 0.5 {
					item.HeadURL = fmt.Sprintf("%v", rand.Int()%4)
				} else {
					item.HeadURL = fmt.Sprintf("%v", 3+rand.Int()%4)
				}

				list = append(list, item)
			}
		} else {
			liverIDs := []string{}
			_, liverList := GetLivers(0, 40)
			for i := 0; i < len(liverList); i++ {
				liverIDs = append(liverIDs, liverList[i].LiverID)
			}

			if len(liverIDs) != 0 {
				for i := 0; i < len(liverIDs) && i < count; i++ {
					item := pb.RankingItemDef{}
					item.Score = 1

					index := rand.Int() % len(liverIDs)
					item.UserID = liverIDs[index]
					item.Nickname = ""
					item.HeadURL = ""
					u, err := FindByUserId(item.UserID)
					if err == nil && u != nil {
						item.Nickname = u.Nickname
						item.HeadURL = u.PhotoUrl
					}

					list = append(list, item)
				}
			}
		}

		glog.Info("default list:", list)

		if rankType == 0 {
			defaultBalanceRank = append(defaultBalanceRank, list...)
		} else if rankType == 1 {
			defaultShareRank = append(defaultShareRank, list...)
		} else if rankType == 2 {
			defaultGiftRank = append(defaultGiftRank, list...)
		} else {
			defaultLiverRank = append(defaultLiverRank, list...)
		}
	}
}

//RestRank ...
func RestRank(rankName string, userID string, value int) {
	cli := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	defer cli.Close()

	cli.ZRem(rankName, userID)
	cli.ZIncrBy(rankName, float64(value), userID)
}

//AddRank ...
func AddRank(rankName string, userID string, value int) {
	cli := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	defer cli.Close()

	cli.ZIncrBy(rankName, float64(value), userID)
}

//GetRank ...
func GetRank(rankName string, userID string) int {
	cli := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	defer cli.Close()
	val, err := cli.ZRank(rankName, userID).Result()
	fmt.Println("getRank ", userID, ",value:", val, ",error:", err)
	if err != nil {
		fmt.Println("not ranking")
		return -1
	}
	return int(val)
}

//GetRankAndScore ...
func GetRankAndScore(rankName string, userID string) (int, int) {
	cli := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	defer cli.Close()
	//val := cli.ZRank(rankName, userID)
	val := cli.ZRevRank(rankName, userID)
	score := cli.ZScore(rankName, userID)

	glog.Info("GetRankAndScore userID", userID, ",rank:", val, ",score:", score)

	return int(val.Val()) + 1, int(score.Val())
}

//RankItem ...
type RankItem struct {
	Rank   int
	UserID string
	Score  int
}

//GetRankRange ...
func GetRankRange(rankName string, start int, end int) []*RankItem {
	cli := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	defer cli.Close()

	list := []*RankItem{}

	//items := cli.ZRevRange(rankName, int64(start), int64(end))
	items := cli.ZRevRangeWithScores(rankName, int64(start), int64(end)).Val()
	fmt.Println("items :", items)

	for i := 0; i < len(items); i++ {
		fmt.Println("item member:", items[i].Member, ", score:", items[i].Score)
		item := &RankItem{}
		item.Rank = i + 1 + start
		item.UserID = string(items[i].Member.(string))
		item.Score = int(items[i].Score)

		list = append(list, item)
	}

	return list
}

//RemoveRank ...
func RemoveRank(rankName string) {
	cli := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	defer cli.Close()

	cli.Del(rankName)
}

func diffTime1Hour(t time.Time) time.Duration {
	midnight := time.Date(t.Year(), t.Month(), t.Day(), 1, 0, 0, 0, time.Local) //今日的凌晨5
	if t.Hour() > 1 {
		nextDay := midnight.AddDate(0, 0, 1) //明天的凌晨5
		tt := nextDay.Sub(t)
		return tt
	}

	tt := midnight.Sub(t)

	return tt
}

//SetRankUpdateTime ...
func SetRankUpdateTime(key string, day string) {
	cli := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	defer cli.Close()

	cli.Set(key, day, diffTime1Hour(time.Now()))
}

//GetRankUpdateTime ...
func GetRankUpdateTime(key string) string {
	cli := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	defer cli.Close()

	val, err := cli.Get(key).Result()
	if err != nil {
		return ""
	}

	return val
}
