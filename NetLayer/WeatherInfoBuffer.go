package NetLayer

import (
	"Common"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//天气信息缓存
type WeatherInfoBuffer struct {
	mUpdateTime   int64              //更新时间
	mTodayBrief   TodayWeatherBrief  //今日天气简要
	mTodayAQI     AQIData            //今日空气质量
	mTodayWeather TodayWeather       //今日天气实况
	mTodayAlert   *TodayAlertWeather //今日预警
	mLimitCar     map[string]string  //限行信息
	mLimitLastDay string             //限行数据最后一天
}

func (this *WeatherInfoBuffer) String() string {
	return fmt.Sprintf("Update time:%d; Today brief:%s; AQI:%s; Weather:%s;Alert:%v;limit Car:{%s},last day:%s;",
		this.mUpdateTime, this.mTodayBrief.String(), this.mTodayAQI.String(), this.mTodayWeather.String(), this.mTodayAlert, this.mLimitCar, this.mLimitLastDay)
}

//更新天气信息
func (this *WeatherInfoBuffer) UpdateWeatherBuff(cityId int) error {
	if gCrawlerConf == nil {
		return errors.New("Error is not initilization.")
	}
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

	if this.NeedUpdateLimit(cityId) {
		Common.DEBUG("WeatherCrawler::UpdateWeatherBuff Update limit, Last date:", this.mLimitLastDay, ";")
		//获取限行
		if err := this.updateLimit(cityId); err != nil {
			Common.ERROR("updateLimit failed, Reason:", err)
			return err
		}
	}

	//获取预警
	//更新天气简要
	CurrTime := time.Now()
	Today := fmt.Sprintf("%04d-%02d-%02d", CurrTime.Year(), int(CurrTime.Month()), CurrTime.Day())
	this.mTodayBrief.Date = this.mTodayWeather.Condition.Updatetime
	this.mTodayBrief.Humidity = this.mTodayWeather.Condition.Humidity
	this.mTodayBrief.Temp = this.mTodayWeather.Condition.Temp
	this.mTodayBrief.WindLevel = this.mTodayWeather.Condition.WindLevel
	this.mTodayBrief.Icon = this.mTodayWeather.Condition.Icon
	this.mTodayBrief.Tips = this.mTodayWeather.Condition.Tips
	this.mTodayBrief.Uvi = this.mTodayWeather.Condition.Uvi
	this.mTodayBrief.Value = this.mTodayAQI.AQI.Value
	this.mTodayBrief.Limit = this.mLimitCar[Today]

	return nil
}

func (this *WeatherInfoBuffer) NeedUpdateLimit(cityId int) bool {
	if len(this.mLimitLastDay) > 0 {
		//分解
		times := strings.Split(this.mLimitLastDay, "-")
		if len(times) != 3 {
			Common.ERROR("Error The times not three.")
			return true
		}
		var LastDay []int
		for _, t := range times {
			i, err := strconv.Atoi(t)
			if err != nil {
				Common.ERROR("Atoi failed, Reason:", err)
				return true
			}
			LastDay = append(LastDay, i)
		}
		var Todays []int
		Curr := time.Now()
		Todays = append(Todays, Curr.Year())
		Todays = append(Todays, int(Curr.Month()))
		Todays = append(Todays, Curr.Day())

		//判断是否超出最后数据日
		for i := 0; i < 3; i++ {
			if Todays[i] > LastDay[i] {
				Common.DEBUG("WeatherCrawler::NeedUpdateLimit Need update limit, LastDay:", this.mLimitLastDay)
				return true
			}
		}
	} else {
		Common.DEBUG("WeatherCrawler::NeedUpdateLimit This is need update.")
		return true
	}

	return false
}

func (this *WeatherInfoBuffer) PostAPIrequest(cityId int, url, token, appcode string, cb func(*http.Response) error) error {
	//设置请求参数
	query := fmt.Sprintf("cityId=%v&token=%v", cityId, token)
	var err error
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

	//发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Common.ERROR("Client Do failed. Reason:", err)
		return err
	}
	defer resp.Body.Close()

	//分解返回信息
	switch resp.StatusCode {
	case http.StatusOK:
		//请求成功 200 OK
		return cb(resp)
	default:
		//失败
		errMsg := fmt.Sprint("HTTP Post Request failed, Status Code:", resp.StatusCode, "; Status:", resp.Status)
		Common.ERROR("Unknown HTTP Code: ", resp.StatusCode, "; Status:", resp.Status)
		return errors.New(errMsg)
	}

	return nil
}

//更新天气实况
func (this *WeatherInfoBuffer) updateConditionWeather(cityId int) error {
	Common.DEBUG("WeatherCrawler::updateConditionWeather Update Condition city id:", cityId)
	return this.PostAPIrequest(cityId, gCrawlerConf["ConditionURL"], gCrawlerConf["ConditionToken"],
		gCrawlerConf["Appcode"], func(resp *http.Response) error {
			var RetData CondtionWeather
			//分解返回信息
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
				this.mTodayWeather = RetData.Data
			case 1:
				Common.ERROR("Error is the Token invalid.code:", RetData.Code, "; msg:", RetData.Msg)
			case 2:
				Common.ERROR("Error is the Sign invalied.code:", RetData.Code, "; msg:", RetData.Msg)
			case 10:
				Common.ERROR("Error is the location invalied.code:", RetData.Code, "; msg:", RetData.Msg)
			default:
				Common.ERROR("Error unknown code:", RetData.Code, "; msg:", RetData.Msg)
			}
			return nil
		})
}

//更新空气情况
func (this *WeatherInfoBuffer) updateAQI(cityId int) error {
	return this.PostAPIrequest(cityId, gCrawlerConf["AQIurl"], gCrawlerConf["AQItoken"],
		gCrawlerConf["Appcode"], func(resp *http.Response) error {
			var RetData AQIReturn
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
				this.mTodayAQI = RetData.Data
			case 1:
				Common.ERROR("Error is the Token invalid.code:", RetData.Code, "; msg:", RetData.Msg)
			case 2:
				Common.ERROR("Error is the Sign invalied.code:", RetData.Code, "; msg:", RetData.Msg)
			case 10:
				Common.ERROR("Error is the location invalied.code:", RetData.Code, "; msg:", RetData.Msg)
			default:
				Common.ERROR("Error unknown code:", RetData.Code, "; msg:", RetData.Msg)
			}

			return nil
		})
}

//更新限行信息
func (this *WeatherInfoBuffer) updateLimit(cityId int) error {
	return this.PostAPIrequest(cityId, gCrawlerConf["LimitURL"], gCrawlerConf["LimitToken"], gCrawlerConf["Appcode"], func(resp *http.Response) error {
		var RetData LimitReturn
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			Common.ERROR("updateLimit Read All failed, Reason:", err)
			return err
		}
		if err := json.Unmarshal(body, &RetData); err != nil {
			Common.ERROR("updateLimit Unmarshal failed, Reason:", err)
			return err
		}

		switch RetData.Code {
		case 0:
			regx := regexp.MustCompile(`^W$`)
			//成功 更新缓存
			for _, Day := range RetData.Data.Limit {
				var InfoMsg string
				if regx.MatchString(Day.Prompt) {
					InfoMsg = "今日不限行"
				} else {
					//限行日
					limits := strings.Split(Day.Prompt, "")
					if len(limits) != 2 {
						Msg := fmt.Sprint("Error the prompt invalid.")
						Common.ERROR(Msg)
						return errors.New(Msg)
					}
					InfoMsg = fmt.Sprintf("今日限行尾号：%s 和 %s", limits[0], limits[1])
				}
				if this.mLimitCar == nil {
					this.mLimitCar = make(map[string]string)
				}
				this.mLimitCar[string(Day.Date)] = InfoMsg
				this.mLimitLastDay = Day.Date
			}
		case 1:
			Common.ERROR("Error is the Token invalid.code:", RetData.Code, "; msg:", RetData.Msg)
		case 2:
			Common.ERROR("Error is the Sign invalied.code:", RetData.Code, "; msg:", RetData.Msg)
		case 10:
			Common.ERROR("Error is the location invalied.code:", RetData.Code, "; msg:", RetData.Msg)
		default:
			Common.ERROR("Error unknown code:", RetData.Code, "; msg:", RetData.Msg)
		}
		return nil
	})
}
