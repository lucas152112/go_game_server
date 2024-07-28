package conf

/**
目前想到的是
1.声网
2.短信
3.版本
4.版本服务器
*/

/**
声网配置
*/
type AgoraConfig struct {
	AppId   string `json:"appId"`   //声网API
	AppCert string `json:"appCert"` //声网Cert
}

/**
网宿云 短信
*/
type WangSuConfig struct {
	Template  string `json:"template"` //短信模板
	UserName  string `json:"userName"` //网宿云的 用户名
	UserKey   string `json:"userKey"`  //网宿云的 用户KEY
	KeySecret string `json:"keySecret"`
	ApiAddr   string `json:"addr"`
}

/**
版本配置
*/
type VerManageConfig struct {
	Addr string `json:"addr"`
}

/**
阿里云短息
*/
type AliYunSmsConfig struct {
	RegionId        string `json:"regionId"` // cn-hangzhou
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	SingName        string `json:"singName"`     // 签名
	TemplateCode    string `json:"templateCode"` // 模板Code
}

type Config struct {
	Version       string          `json:"version"`
	Agora         AgoraConfig     `json:"agora"`
	ChinaSms      WangSuConfig    `json:"chinaSms"`
	GlobalSms     WangSuConfig    `json:"globalSms"`
	VersionManage VerManageConfig `json:"versionManage"`
	AliSms        AliYunSmsConfig `json:"aliYunSms"`  //阿里云短信
	ChinaGuoDu    GuoDuConfig     `json:"guoduChina"` //
	ChuangLan     ChuangLanConfig `json:"chuangLan"`
	LiveGameTax   LiveTaxConfig   `json:"liveGameTax"` //直播场抽水配置
}

type GuoDuConfig struct {
	UserName string `json:"userName"` //国都短信 用户名
	Password string `json:"password"` //
	ApiAddr  string `json:"addr"`
}

type LiveTaxConfig struct {
	PumpPercent float32 `json:"pumpPercent"`
}

type ChuangLanConfig struct {
	UserName string `json:"userName"` //创蓝短信 用户名
	Password string `json:"password"` //
	ApiAddr  string `json:"addr"`
}
