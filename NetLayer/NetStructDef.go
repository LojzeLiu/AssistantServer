package NetLayer

import (
	"fmt"
)

type WeatherRC struct {
	C int    `json:"c"`
	P string `json:"p"`
}

func (this WeatherRC) String() string {
	return fmt.Sprintf("c:%d;p:%s", this.C, this.P)
}

//城市信息
type CityInfo struct {
	CityId   int32  `json:"cityId"`   //城市ID
	Counname string `json:"counname"` //国家
	Name     string `json:"name"`     //区名称
	Pname    string `json:"pname"`    //城市名称
}

func (this CityInfo) String() string {
	return fmt.Sprintf("Citty id:%d; Counname:%s; Name:%s; Pname:%s;", this.CityId, this.Counname, this.Name, this.Pname)
}

//天气预警
type TodayAlertWeather struct {
	Content string `json:"content"`  //内容
	InfoID  int    `json:"infoid"`   //预警ID
	Level   string `json:"level"`    //等级
	Name    string `json:"name"`     //预警名称
	PubTime string `json:"pub_time"` //发布时间
	Title   string `json:"title"`    //标题
	Type    string `json:"type"`     //预警类型
}

func (this TodayAlertWeather) String() string {
	return fmt.Sprintf("Content:%s;InfoID:%s;Level:%s;Name:%s;PubTime:%s;Title:%s;Type:%s;",
		this.Content, this.InfoID, this.Level, this.Name, this.PubTime, this.Title, this.Type)
}

//天气预警数据
type TodayAlertData struct {
	Alert []TodayAlertWeather `json:"alert"`
	City  CityInfo            `json:"city"`
}

func (this TodayAlertData) String() string {
	return fmt.Sprintf("alert:{%s},city:{%s}", this.Alert, this.City.String())
}

//天气预警返回协议
type TodayAlertReturn struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"` //执行状态消息
	RC   WeatherRC      `json:"rc"`
	Data TodayAlertData `json:"data"`
}

func (this TodayAlertReturn) String() string {
	return fmt.Sprintf("Code:%d; Msg:%s; RC:{%s}; Data:{%s}", this.Code, this.Msg, this.RC.String(), this.Data.String())
}

//天气实况
type TodayCondition struct {
	Condition   string `json:"condition"`   //天气阴晴
	ConditionId string `json:"conditionId"` //实时天气 id
	Humidity    string `json:"humidity"`    //湿度
	Icon        string `json:"icon"`        //图标
	Pressure    string `json:"pressure"`    //气压
	RealFeel    string `json:"realFeel"`    //体感温度
	SunRise     string `json:"sunRise"`     //日出时间
	SunSet      string `json:"sunSet"`      //日落时间
	Temp        string `json:"temp"`        //温度
	Tips        string `json:"tips"`        //提示
	Updatetime  string `json:"updatetime"`  //更新时间
	Uvi         string `json:"uvi"`         //紫外线强度
	WindDir     string `json:"windDir"`     //风向
	WindLevel   string `json:"windLevel"`   //风级
	WindSpeed   string `json:"windSpeed"`   //风速
}

func (this *TodayCondition) String() string {
	return fmt.Sprintf("Condtion:%s; ConditionId:%s; Humidity:%s; icon:%s; Pressure:%s;RealFeel:%s; SunRise:%s; SunSet:%s; Temp:%s; Tips:%s; Updatetime:%s; Uvi:%s; WindDir:%s;windLevel:%s;WindSpeed:%s",
		this.Condition, this.ConditionId, this.Humidity, this.Icon, this.Pressure, this.RealFeel, this.SunRise, this.SunSet, this.Temp, this.Tips, this.Updatetime, this.Uvi,
		this.WindDir, this.WindLevel, this.WindSpeed)
}

//天气实况
type TodayWeather struct {
	City      CityInfo       `json:"city"`      //城市信息
	Condition TodayCondition `json:"condition"` //天气信息
}

func (this *TodayWeather) String() string {
	return fmt.Sprintf("City:{%s};Condition:{%s}", this.City.String(), this.Condition.String())
}

//天气实况API返回协议
type CondtionWeather struct {
	Code int          `json:"code"` //执行状态码
	Data TodayWeather `json:"data"` //天气数据
	Msg  string       `json:"msg"`  //执行状态消息
	RC   WeatherRC    `json:"rc"`
}

func (this CondtionWeather) String() string {
	return fmt.Sprintf("code:%d;msg:%s;Data:{%s};", this.Code, this.Msg, this.Data.String())
}

//空气质量
type TodayAQI struct {
	CityName string `json:"cityName"` //城市名称
	Co       string `json:"co"`       //一氧化碳指数
	No2      string `json:"no2"`      //二氧化氮指数
	O3       string `json:"o3"`       //臭氧指数
	Pm10     string `json:"pm10"`     //PM10 指数
	Pm25     string `json:"pm25"`     //PM2.5 指数
	Pubtime  string `json:"pubtime"`  //发布时间戳
	Rank     string `json:"rank"`     //全国排名
	So2      string `json:"so2"`      //二氧化硫浓度
	Value    string `json:"value"`    //空气质量指数值
}

func (this TodayAQI) String() string {
	return fmt.Sprintf("CityName:%s;Co:%s;No2:%s;O3:%s;Pm10:%s;Pm25:%s;Pubtime:%s;Rank:%s;So2:%s;Value:%s",
		this.CityName, this.Co, this.No2, this.O3, this.Pm10, this.Pm25, this.Pubtime, this.Rank, this.So2, this.Value)
}

//空气质量数据
type AQIData struct {
	City CityInfo `json:"city"` //城市信息
	AQI  TodayAQI `json:"aqi"`  //空气信息
}

func (this AQIData) String() string {
	return fmt.Sprintf("City:{%s},AQI:{%s}", this.City.String(), this.AQI.String())
}

//空气质量API返回协议
type AQIReturn struct {
	Code int       `json:"code"`
	Msg  string    `json:"msg"` //执行状态消息
	RC   WeatherRC `json:"rc"`
	Data AQIData   `json:"data"`
}

func (this AQIReturn) String() string {
	return fmt.Sprintf("Code:%d;Msg:%s;RC:{%s};Data:{%s}", this.Code, this.Msg, this.RC, this.Data.String())
}

//限行信息
type DayLimit struct {
	Date   string `json:"date"`
	Prompt string `json:"Prompt"`
}

func (this DayLimit) String() string {
	return fmt.Sprintf("Date:%s; Prompt:%s", this.Date, this.Prompt)
}

//限行数据
type LimitData struct {
	City  CityInfo   `json:"city"`
	Limit []DayLimit `json:"limit"`
}

func (this LimitData) String() string {
	return fmt.Sprintf("City:{%s};Limit:{%v}", this.City.String(), this.Limit)
}

//限行返回协议
type LimitReturn struct {
	Code int       `json:"code"`
	Msg  string    `json:"msg"` //执行状态消息
	RC   WeatherRC `json:"rc"`
	Data LimitData `json:"data`
}

func (this LimitReturn) String() string {
	return fmt.Sprintf("code:%d;msg:%s;data{%s};rc:{%s}", this.Code, this.Msg, this.Data.String(), this.RC.String())
}

//向客户端返回今天天气简要
type TodayWeatherBrief struct {
	Title     string `json:"title`      //标题
	Date      string `json:"date"`      //日期
	Humidity  string `json:"humidity"`  //湿度
	Temp      string `json:"temp"`      //温度
	WindLevel string `json:"windLevel"` //风级
	Icon      string `json:"icon"`      //图标
	Tips      string `json:"tips"`      //提示
	Uvi       string `json:"uvi"`       //紫外线强度
	Value     string `json:"value"`     //空气质量指数值
	Limit     string `json:"limit"`     //汽车限行
}

func (this *TodayWeatherBrief) String() string {
	return fmt.Sprintf("Title:%s; Date:%s;Humidity:%s;Temp:%s;WindLevel:%s;Icon:%s;Tips:%s;Uvi:%s;Value:%s;Limit:%s;",
		this.Title, this.Date, this.Humidity, this.Temp, this.WindLevel, this.Icon, this.Tips, this.Uvi, this.Value, this.Limit)
}
