package NetLayer

import (
	"Common"
	"golang.org/x/net/websocket"
	"sync"
	"time"
)

var gMutex sync.Mutex
var gClientID int64 = 80000
var gClients map[int64]*WeatherServer = make(map[int64]*WeatherServer)

func WSAccept(ws *websocket.Conn) {
	Common.DEBUG("On Accept, Addr:", ws.RemoteAddr())
	CurrClient := &WeatherServer{}
	CurrClient.Init()
	go timedWork()
	gMutex.Lock()
	gClientID++
	CurrClientID := gClientID
	gClients[CurrClientID] = &CurrClient
	gMutex.Unlock()

	CurrClient.Listener(ws)
	gMutex.Lock()
	delete(gClients, CurrClientID)
	gMutex.Unlock()
	Common.DEBUG("Deal end...")
}

func SendAllMsg(msg CmdInfo) {
	for id, Client := range gClients {
		if err := curr.SendMsg(msg); err != nil {
			Common.DEBUG("Send msg failed, Reason:", err, "; ID:", id)
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
