package NetLayer

import (
	"Common"
	"encoding/json"
	"errors"
	"golang.org/x/net/websocket"
	"time"
)

//请求类型
type ASK_TYPE int32

const (
	QueryGetWeather      ASK_TYPE = iota //获取天气信息
	QueryGetDrivingLimit                 //获取限行
)

//协议类型
type CMD_TYEP int16

const (
	CmdAlive CMD_TYEP = iota
	CmdLogin
)

//消息协议
type CmdMsg struct {
	Msg string `json:"msg"`
}

//协议头
type CmdInfo struct {
	Cmd        CMD_TYEP    `json:"cmd"`
	Version    string      `json:"version"`
	HandleData interface{} `json:"handle_data"`
}

type WeatherServer struct {
	mWSconn     *websocket.Conn
	mLastLive   time.Time //最后心跳时间
	mFuncHandle map[CMD_TYEP]func(interface{}) error
}

func (this *WeatherServer) Init() {
	this.mFuncHandle = make(map[CMD_TYEP]func(interface{}) error)
	this.mFuncHandle[0] = this.liveHandle
}

func (this *WeatherServer) Listener(ws *websocket.Conn) {
	this.mWSconn = ws
	for {
		var RecvBuff string
		if err := websocket.Message.Receive(ws, &RecvBuff); err != nil {
			Common.DEBUG("Receive failed. Reason:", err)
			break
		}

		//分解信息
		var ci CmdInfo
		if err := json.Unmarshal([]byte(RecvBuff), &ci); err != nil {
			Common.ERROR("Unmarshal failed. Reason:", err, ";Recv:", RecvBuff)
			continue
		}

		//执行响应的操作
		if err := this.cmdHandle(ci.Cmd, ci.HandleData); err != nil {
			Common.ERROR("cmdHandle failed. Reason:", err)
			continue
		}

	}
}

func (this *WeatherServer) SendMsg(cmd CmdInfo) error {
	if this.mWSconn == nil {
		return errors.New("Not initialization WS conn. ")
	}
	msg, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	sendData := string(msg)
	if err := websocket.Message.Send(this.mWSconn, sendData); err != nil {
		return err
	}
	return nil
}

func (this *WeatherServer) cmdHandle(cmd CMD_TYEP, hd interface{}) error {
	f := this.mFuncHandle[cmd]
	if f == nil {
		return errors.New("Not the handle.")
	}
	return f(hd)
}

//心跳处理
func (this *WeatherServer) liveHandle(hd interface{}) error {
	Common.DEBUG("This is live tick.")
	this.mLastLive = time.Now()

	return nil
}
