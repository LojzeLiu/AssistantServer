package NetLayer

import (
	"Common"
	"golang.org/x/net/websocket"
	"time"
)

var gClients []*WeatherServer

func WSAccept(ws *websocket.Conn) {
	Common.DEBUG("On Accept, Addr:", ws.RemoteAddr())
	CurrClient := &WeatherServer{}
	CurrClient.Init()
	go timedWork()
	gClients = append(gClients, CurrClient)

	CurrClient.Listener(ws)
	Common.DEBUG("Deal end...")
}

func SendAllMsg(msg CmdInfo) {
	for num, curr := range gClients[0:] {
		if err := curr.SendMsg(msg); err != nil {
			Common.DEBUG("Send msg failed, Reason:", err, "; num:", num)
			continue
		}
	}
}

func timedWork() {
	for {
		var msg CmdInfo
		var data CmdMsg
		data.Msg = "Hello this is golang websocket server."
		msg.Cmd = 0
		msg.Version = "1.0.0.0"
		msg.HandleData = data
		SendAllMsg(msg)

		time.Sleep(time.Second * 10)
	}
}
