package util

import "strings"

const (
	LANGUAGE_EN   = 0   // 英文
	LANGUAGE_CHS  = 1   // 简体
	LANGUAGE_CHT  = 2   // 繁体
)

type Lang struct {
	lanName int
}

func L() *Lang  {
	return &Lang{lanName:LANGUAGE_CHS}
}

func (this *Lang) To( lan int  ) *Lang {
	this.lanName = lan
	return this
}

func (this *Lang) Msg(msg string ) string {
	return lan{}.translate( msg,this.lanName )
}


type lan struct {}

func (this lan) translate( msg string ,language int  ) string  {
	message := ""
	msg = strings.ToLower(msg)  //小写
	switch language {
	case LANGUAGE_CHT: message = this.chs(msg)
	case LANGUAGE_CHS: message = this.chs(msg)
	default:
		message = msg
	}
	return message
}

func (this lan) chs(msg string) string   {

	if message,ok := Langage_CHS[msg];!ok{
		return msg
	}else{
		return message
	}
}
func (this lan) cht(msg string ) string   {
	return msg
}