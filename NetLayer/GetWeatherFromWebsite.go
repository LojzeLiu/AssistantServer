package NetLayer

import (
	"Common"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
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

//限行信息
type DayLimit struct {
	Date   string `json:"date"`
	Prompt string `json:"Prompt"`
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
	mUpdateTime   time.Time           //更新时间
	mTodayBrief   TodayWeatherBrief   //今日天气简要
	mTodayAQI     TodayAQI            //今日空气质量
	mTodayWeather TodayWeather        //今日天气实况
	mTodayAlert   *TodayAlertWeather  //今日预警
	mLimitCar     map[string]DayLimit //限行信息
}

type WeatherCrawler struct {
	mCrawlerConf   Common.Configer    //配置信息
	mBufferTimeOut time.Duration      //缓存更新间隔
	mWeatherBuff   *WeatherInfoBuffer //天气情况缓存
}

//初始化
func (this *WeatherCrawler) Init(conf *Common.Configer) error {
	if this.mCrawlerConf, err = conf.GetConf("CRAWLER"); err != nil {
		return err
	}

	if len(this.mCrawlerConf) <= 0 {
		return errors.New("Crawler configur not exit.")
	}

	gap := strconv.Atoi(this.mCrawlerConf["BufferTimeout"])
	if gap <= 0 {
		this.mBufferTimeOut = time.Duration(gap)
	} else {
		this.mBufferTimeOut = Common.Hour * 2
	}
	Common.DEBUG("Buffer timeout:", gap)

	return nil
}

//获取今日天气简要
func (this *WeatherCrawler) GetTodayBrief(cityId int) (TodayWeatherBrief, TodayAlertWeather, error) {
	//检测缓存
	currTime := time.Now()
	if this.mWeatherBuff.mUpdateTime.Sub(currTime) >= this.mBufferTimeOut {
		//超时更新缓存
		if err := this.UpdateWeatherBuff(cityId); err != nil {
			return nil, nil, err
		}
	}

	//返回信息
	return this.mWeatherBuff.mTodayBrief, this.mWeatherBuff.mTodayAlert, nil
}

//更新天气信息
func (this *WeatherCrawler) UpdateWeatherBuff(cityId int) error {
	//获取天气实况
	if err := this.updateConditionWeather(cityId); err != nil {
		return err
	}
	//获取空气质量
	//获取限行
	//获取预警
	//更新天气简要
	Today := fmt.Sprint(time.Now().Year(), "-", time.Now().Month(), "-", time.Now().Day())
	this.mWeatherBuff.mTodayBrief.Date = this.mWeatherBuff.mTodayWeather.Updatetime
	this.mWeatherBuff.mTodayBrief.Humidity = this.mWeatherBuff.mTodayWeather.Humidity
	this.mWeatherBuff.mTodayBrief.Temp = this.mWeatherBuff.mTodayWeather.Temp
	this.mWeatherBuff.mTodayBrief.WindLevel = this.mWeatherBuff.mTodayWeather.WindLevel
	this.mWeatherBuff.mTodayBrief.Icon = this.mWeatherBuff.mTodayWeather.Icon
	this.mWeatherBuff.mTodayBrief.Tips = this.mWeatherBuff.mTodayWeather.Tips
	this.mWeatherBuff.mTodayBrief.Uvi = this.mWeatherBuff.mTodayWeather.Uvi
	this.mWeatherBuff.mTodayBrief.Value = this.mWeatherBuff.mTodayAQI.Value
	this.mWeatherBuff.mTodayBrief.Limit = this.mWeatherBuff.mLimitCar[Today]

	return nil
}

//更新天气实况
func (this *WeatherCrawler) updateConditionWeather(cityId int) error {
	//设置请求参数
	var query string
	fmt.Fprintf(query, "cityId=%v&token=%v", cityId, this.mCrawlerConf["token"])
	req, err := http.NewRequest("POST", this.mCrawlerConf["ConditionURL"], query)
	if err != nil {
		return err
	}

	//设置请求头
	APPCODE := "APPCODE"
	APPCODE = APPCODE + this.mCrawlerConf["Appcode"]
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Authorization", APPCODE)

	//发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var TodayCondition TodayWeather
	var RetData AlicityWeather
	RetData.Data.(TodayWeather)
	//分解返回信息
	switch resp.StatusCode {
	case http.StatusOK:
		//请求成功 200 OK
		//解析Json
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(body, &RetData); err != nil {
			return err
		}
		Common.DEBUG("Body:", string(body))
		//更新缓存
		this.mWeatherBuff.mTodayWeather = TodayCondition
	default:
		//失败
		errMsg := fmt.Sprint("HTTP Post request failed, status code:", resp.StatusCode, "; Status:", resp.Status)
		return errors.New(errMsg)
	}

	return nil
}
