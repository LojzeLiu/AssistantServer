package NetLayer

import (
	"Common"
	"encoding/json"
	"errors"
	"golang.org/x/net/websocket"
	"strconv"
	"time"
)

type CmdInfo struct {
	Cmd        int16
	Version    string
	HandleData interface{}
}

type WeatherServer struct {
	mLastLive   time.Time //最后心跳时间
	mFuncHandle map[int16]func(interface{})
}

func (this *WeatherServer) Init() {
	ws.mFuncHandle[0] = ws.liveHandle
}

func (this *WeatherServer) Listener(ws *websocket.Conn) {
	for {
		var RecvBuff [2048]byte
		if err := websocket.Message.Receive(ws, RecvBuff); err != nil {
			Common.DEBUG("Error:", err)
			break
		}

		//分解信息
		var ci CmdInfo
		if err := json.Unmarshal(RecvBuff, &ci); err != nil {
			Common.ERROR("Error:", err)
			continue
		}

		//执行响应的操作
		if err := this.cmdHandle(ci.Cmd, ci.HandleData); err != nil {
			Common.ERROR("Error:", err)
			continue
		}
	}
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
}
