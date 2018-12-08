package NetLayer

import (
	"Common"
	"errors"
	"strconv"
	"time"
)

var gCrawlerConf map[string]string //配置信息

//天气爬虫
type WeatherCrawler struct {
	mBufferTimeOut time.Duration              //缓存更新间隔
	mWeatherBuff   map[int]*WeatherInfoBuffer //天气情况缓存
}

//初始化
func (this *WeatherCrawler) Init(conf *Common.Configer) error {
	var err error
	if gCrawlerConf, err = conf.GetConf("CRAWLER"); err != nil {
		return err
	}

	if len(gCrawlerConf) <= 0 {
		return errors.New("Crawler configur not exit.")
	}

	gap, err := strconv.Atoi(gCrawlerConf["BufferTimeout"])
	if err != nil {
		return err
	}
	if gap <= 0 {
		this.mBufferTimeOut = time.Duration(gap) * Common.Hour
	} else {
		this.mBufferTimeOut = Common.Hour * 2
	}
	this.mWeatherBuff = make(map[int]*WeatherInfoBuffer)
	Common.DEBUG("WeatherCrawler::Init Buffer timeout:", gap)

	return nil
}

//获取今日天气简要
func (this *WeatherCrawler) GetTodayBrief(cityId int) (*TodayWeatherBrief, *TodayAlertWeather, error) {
	//检测缓存
	Common.DEBUG("Get city id:", cityId)
	if this.mWeatherBuff[cityId] == nil {
		this.mWeatherBuff[cityId] = &WeatherInfoBuffer{}
	}
	lpBuff := this.mWeatherBuff[cityId]
	currTime := time.Now().Unix()
	if lpBuff.mUpdateTime <= 0 || (currTime-lpBuff.mUpdateTime) >= int64(this.mBufferTimeOut*60*60) {
		//超时更新缓存
		if err := lpBuff.UpdateWeatherBuff(cityId); err != nil {
			return nil, nil, err
		}
		Common.DEBUG("WeatherCrawler::GetTodayBrif Update weather buffer.")
		lpBuff.mUpdateTime = time.Now().Unix()
	}

	//返回信息
	return &lpBuff.mTodayBrief, lpBuff.mTodayAlert, nil
}
