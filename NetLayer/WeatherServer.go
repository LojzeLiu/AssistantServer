package NetLayer

import (
	"Common"
	"encoding/json"
	"errors"
	"golang.org/x/net/websocket"
	"time"
)

type CmdMsg struct {
	msg string
}

type CmdInfo struct {
	Cmd        int16       `json:cmd`
	Version    string      `json:version`
	HandleData interface{} `json:handle_data`
}

type WeatherServer struct {
	mWSconn     *websocket.Conn
	mLastLive   time.Time //最后心跳时间
	mFuncHandle map[int16]func(interface{}) error
}

func (this *WeatherServer) Init() {
	this.mFuncHandle = make(map[int16]func(interface{}) error)
	this.mFuncHandle[0] = this.liveHandle
}

func (this *WeatherServer) Listener(ws *websocket.Conn) {
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

		this.mWSconn = ws
	}
}

func (this *WeatherServer) SendMsg(cmd CmdInfo) error {
	if this.mWSconn == nil {
		return nil
	}
	msg, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	if err := websocket.Message.Send(this.mWSconn, msg); err != nil {
		return err
	}
	return nil
}

func (this *WeatherServer) cmdHandle(cmd int16, hd interface{}) error {
	f := this.mFuncHandle[cmd]
	if f == nil {
		return errors.New("Not the handle.")
	}
	return f(hd)
}

func (this *WeatherServer) liveHandle(hd interface{}) error {
	Common.DEBUG("This is live tick.")
	this.mLastLive = time.Now()

	return nil
}
