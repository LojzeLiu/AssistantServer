package NetLayer

import (
	"Common"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
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
	fmt.Sprintf("Code:%d;Msg:%s;RC:{%s};Data:{%s}", this.Code, this.Msg, this.RC, this.Data.String())
}

//限行信息
type DayLimit struct {
	Date   string `json:"date"`
	Prompt string `json:"Prompt"`
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

//天气信息缓存
type WeatherInfoBuffer struct {
	mUpdateTime   int64               //更新时间
	mTodayBrief   TodayWeatherBrief   //今日天气简要
	mTodayAQI     AQIData             //今日空气质量
	mTodayWeather TodayWeather        //今日天气实况
	mTodayAlert   *TodayAlertWeather  //今日预警
	mLimitCar     map[string]DayLimit //限行信息
}

type WeatherCrawler struct {
	mCrawlerConf   map[string]string //配置信息
	mBufferTimeOut time.Duration     //缓存更新间隔
	mWeatherBuff   WeatherInfoBuffer //天气情况缓存
}

//初始化
func (this *WeatherCrawler) Init(conf *Common.Configer) error {
	var err error
	if this.mCrawlerConf, err = conf.GetConf("CRAWLER"); err != nil {
		return err
	}

	if len(this.mCrawlerConf) <= 0 {
		return errors.New("Crawler configur not exit.")
	}

	gap, err := strconv.Atoi(this.mCrawlerConf["BufferTimeout"])
	if err != nil {
		return err
	}
	if gap <= 0 {
		this.mBufferTimeOut = time.Duration(gap) * Common.Hour
	} else {
		this.mBufferTimeOut = Common.Hour * 2
	}
	Common.DEBUG("Buffer timeout:", gap)

	return nil
}

//获取今日天气简要
func (this *WeatherCrawler) GetTodayBrief(cityId int) (*TodayWeatherBrief, *TodayAlertWeather, error) {
	//检测缓存
	currTime := time.Now().Unix()
	if this.mWeatherBuff.mUpdateTime <= 0 || (currTime-this.mWeatherBuff.mUpdateTime) >= int64(this.mBufferTimeOut*60*60) {
		//超时更新缓存
		if err := this.UpdateWeatherBuff(cityId); err != nil {
			return nil, nil, err
		}
		Common.DEBUG("Update weather buffer.")
		this.mWeatherBuff.mUpdateTime = time.Now().Unix()
	}

	//返回信息
	return &this.mWeatherBuff.mTodayBrief, this.mWeatherBuff.mTodayAlert, nil
}

//更新天气信息
func (this *WeatherCrawler) UpdateWeatherBuff(cityId int) error {
	//获取天气实况
	if err := this.updateConditionWeather(cityId); err != nil {
		Common.ERROR("updateConditionWeather failed. Reason:", err)
		return err
	}

	//获取空气质量
	if err := this.updateAQI(cityId); err != nil {
		Common.ERROR("updateAQI failed. Reason:", err)
		return err
	}

	//获取限行
	//获取预警
	//更新天气简要
	Today := fmt.Sprint(time.Now().Year(), "-", time.Now().Month(), "-", time.Now().Day())
	this.mWeatherBuff.mTodayBrief.Date = this.mWeatherBuff.mTodayWeather.Condition.Updatetime
	this.mWeatherBuff.mTodayBrief.Humidity = this.mWeatherBuff.mTodayWeather.Condition.Humidity
	this.mWeatherBuff.mTodayBrief.Temp = this.mWeatherBuff.mTodayWeather.Condition.Temp
	this.mWeatherBuff.mTodayBrief.WindLevel = this.mWeatherBuff.mTodayWeather.Condition.WindLevel
	this.mWeatherBuff.mTodayBrief.Icon = this.mWeatherBuff.mTodayWeather.Condition.Icon
	this.mWeatherBuff.mTodayBrief.Tips = this.mWeatherBuff.mTodayWeather.Condition.Tips
	this.mWeatherBuff.mTodayBrief.Uvi = this.mWeatherBuff.mTodayWeather.Condition.Uvi
	this.mWeatherBuff.mTodayBrief.Value = this.mWeatherBuff.mTodayAQI.Value
	this.mWeatherBuff.mTodayBrief.Limit = this.mWeatherBuff.mLimitCar[Today].Prompt

	return nil
}

func (this *WeatherCrawler) CreateAPIpost(cityID int, token, url, appcode string) (*http.Request, error) {
	//设置请求参数
	query := fmt.Sprintf("cityId=%v&token=%v", cityID, token)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(query)))
	if err != nil {
		Common.ERROR("New Request Failed, Reason:", err)
		return err
	}

	//设置请求头
	APPCODE := "APPCODE "
	APPCODE = APPCODE + appcode
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Authorization", APPCODE)

	return req
}

//更新天气实况
func (this *WeatherCrawler) updateConditionWeather(cityId int) error {
	Common.DEBUG("Update Condition city id:", cityId)

	//生成头
	req, err := this.CreateAPIpost(cityId, this.mCrawlerConf["ConditionToken"], this.mCrawlerConf["ConditionURL"], this.mCrawlerConf["Appcode"])
	if err != nil {
		Common.ERROR("Create API Post head failed. Reason:", err)
		return err
	}

	var RetData CondtionWeather
	//分解返回信息
	switch resp.StatusCode {
	case http.StatusOK:
		//请求成功 200 OK
		//解析Json
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			Common.ERROR("Read All filed. Reason:", err)
			return err
		}
		if err := json.Unmarshal(body, &RetData); err != nil {
			Common.ERROR("Unmarshal Filed. Reason:", err)
			return err
		}
		switch RetData.Code {
		case 0:
			//成功
			//更新缓存
			this.mWeatherBuff.mTodayWeather = RetData.Data
		case 1:
			Common.ERROR("Error is the Token invalid.code:", RetData.Code, "; msg:", RetData.Msg)
		case 2:
			Common.ERROR("Error is the Sign invalied.code:", RetData.Code, "; msg:", RetData.Msg)
		case 10:
			Common.ERROR("Error is the location invalied.code:", RetData.Code, "; msg:", RetData.Msg)
		default:
			Common.ERROR("Error unknown code:", RetData.Code, "; msg:", RetData.Msg)
		}
	default:
		//失败
		errMsg := fmt.Sprint("HTTP Post Request failed, Status Code:", resp.StatusCode, "; Status:", resp.Status)
		Common.ERROR("Unknown HTTP Code: ", resp.StatusCode, "; Status:", resp.Status)
		return errors.New(errMsg)
	}

	return nil
}

//更新空气情况
func (this *WeatherCrawler) updateAQI(cityId int) error {
	//生成头
	req, err := this.CreateAPIpost(cityId, this.mCrawlerConf["AQItoken"], this.mCrawlerConf["AQIurl"], this.mCrawlerConf["Appcode"])
	if err != nil {
		Common.ERROR("Create API Post head failed. Reason:", err)
		return err
	}

	//发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Common.ERROR("Client Do failed. Reason:", err)
		return err
	}
	defer resp.Body.Close()

	var RetData AQIReturn
	switch resp.StatusCode {
	case http.StatusOK:
		//200 OK
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			Common.ERROR("updateAQI Read all failed. Reason:", err)
			return err
		}
		//解析到JSON
		if err := json.Unmarshal(body, &RetData); err != nil {
			Common.ERROR("updateAQI Unmarshal failed, Reason:", err)
			return err
		}
		switch RetData.Code {
		case 0:
			//成功更新缓存
			this.mWeatherBuff.mTodayAQI = AQIReturn.Data
		case 1:
			Common.ERROR("Error is the Token invalid.code:", RetData.Code, "; msg:", RetData.Msg)
		case 2:
			Common.ERROR("Error is the Sign invalied.code:", RetData.Code, "; msg:", RetData.Msg)
		case 10:
			Common.ERROR("Error is the location invalied.code:", RetData.Code, "; msg:", RetData.Msg)
		default:
			Common.ERROR("Error unknown code:", RetData.Code, "; msg:", RetData.Msg)
		}
	default:
		//失败
		errMsg := fmt.Sprint("HTTP Post Request failed, Status Code:", resp.StatusCode, "; Status:", resp.Status)
		Common.ERROR("Unknown HTTP Code: ", resp.StatusCode, "; Status:", resp.Status)
		return errors.New(errMsg)
	}

	return nil
}
