package NetLayer

import (
	"Common"
	"fmt"
	"golang.org/x/net/websocket"
	"time"
)

var gClients []WeatherServer

func WSAccept(ws *websocket.Conn) {
	Common.DEBUG("On Accept, Addr:", ws.RemoteAddr())
	CurrClient := WeatherServer{}
	CurrClient.Init()
	go timedWork()
	gClients = append(gClients, CurrClient)

	CurrClient.Listener(ws)
	Common.DEBUG("Deal end...")
}

func SendAllMsg(msg CmdInfo) {
	for _, curr := range gClients[0:] {
		fmt.Println("Send a client msg.")
		if err := curr.SendMsg(msg); err != nil {
			Common.DEBUG("Send msg failed, Reason:", err)
			continue
		}
	}
}

func timedWork() {
	for {
		var msg CmdInfo
		var data CmdMsg
		data.msg = "Hello this is golang websocket server."
		msg.Cmd = 0
		msg.Version = "1.0.0.0"
		msg.HandleData = data
		SendAllMsg(msg)
		fmt.Println("Send all client msg.")

		time.Sleep(time.Second * 10)
	}
}
