package user

const (
	MGO_USER_USER_NICKNAME = "user_nickname"
)

type UserNickName struct {
	UserId     string  `bson:"userId"`
	NickName   string  `bson:"NickName"`
	ModifyNum  uint    `bson:"modifyNum"`
}
