package NetLayer

import (
	"Common"
	"golang.org/x/net/websocket"
)

func WSAccept(ws *websocket.Conn) {
	Common.DEBUG("On Accept ")
	for {
		var reply string
		if err := websocket.Message.Receive(ws, &reply); err != nil {
			Common.ERROR("Error:", err)
			break
		}
		Common.DEBUG("Recv:", reply)
		if err := websocket.Message.Send(ws, reply); err != nil {
			Common.ERROR("Error:", err)
			break
		}
	}
}

type WeatherServer struct {
}
