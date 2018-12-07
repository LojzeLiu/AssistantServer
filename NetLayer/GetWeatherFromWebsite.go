package NetLayer

import (
	//"fmt"
	//"github.com/gocolly/colly"
	"Common"
	"errors"
	"time"
)

//城市信息
type CityInfo struct {
	CityId   int32  `json:"cityId"`   //城市ID
	Counname string `json:"counname"` //国家
	Name     string `json:"name"`     //区名称
	Pname    string `json:"pname"`    //城市名称
}

//天气预警
type TodayAlertWeather struct {
	Content     string `json:"content"`     //内容
	InfoID      string `json:"infoid"`      //预警ID
	Level       string `json:"level"`       //等级
	Name        string `json:"name"`        //预警名称
	Pub_time    string `json:"pub_time"`    //发布时间
	Title       string `json:"title"`       //标题
	Type        string `json:"type"`        //预警类型
	Update_time string `json:"update_time"` //更新时间
}

//天气实况
type TodayWeather struct {
	City        CityInfo `json:"city"`        //城市信息
	Condition   string   `json:"condition"`   //天气阴晴
	ConditionId string   `json:"conditionId"` //实时天气 id
	Humidity    string   `json:"humidity"`    //湿度
	Icon        string   `json:"icon"`        //图标
	Pressure    string   `json:"pressure"`    //气压
	RealFeel    string   `json:"realFeel"`    //体感温度
	SunRise     string   `json:"sunRise"`     //日出时间
	SunSet      string   `json:"sunSet"`      //日落时间
	Temp        string   `json:"temp"`        //温度
	Tips        string   `json:"tips"`        //提示
	Updatetime  string   `json:"updatetime"`  //更新时间
	Uvi         string   `json:"uvi"`         //紫外线强度
	WindDir     string   `json:"windDir"`     //风向
	WindLevel   string   `json:"windLevel"`   //风级
	WindSpeed   string   `json:"windSpeed"`   //风速
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

type WeatherRC struct {
	C int    `json:"c"`
	P string `json:"p"`
}

//天气API返回协议
type AlicityWeather struct {
	Code int         `json:"code"` //执行状态码
	Data interface{} `json:"data"` //天气数据
	Msg  string      `json:"msg"`  //执行状态消息
	RC   WeatherRC   `json:"rc"`
}

//向客户端返回今天天气简要
type TodayWeatherBrief struct {
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

//天气信息缓存
type WeatherInfoBuffer struct {
	mUpdateTime   time.Time         //更新时间
	mTodayBrief   TodayWeatherBrief //今日天气简要
	mTodayAQI     TodayAQI          //今日空气质量
	mTodayWeather TodayWeather      //今日天气实况
}

type WeatherCrawler struct {
	mCrawlerConf Common.Configer    //配置信息
	mWeatherBuff *WeatherInfoBuffer //天气情况缓存
}

//初始化
func (this *WeatherCrawler) Init(conf *Common.Configer) error {
	if this.mCrawlerConf, err = conf.GetConf("CRAWLER"); err != nil {
		return err
	}

	if len(this.mCrawlerConf) <= 0 {
		return errors.New("Crawler configur not exit.")
	}

	return nil
}

//获取今日天气简要
func (this *WeatherCrawler) GetTodayBrief() (TodayWeatherBrief, error) {
	//检测缓存
	currTime := time.Now()
	if this.mWeatherBuff.mUpdateTime.Sub(currTime) >= Common.Hour*2 {
		//超时更新缓存
	}
}

//获取天气信息
