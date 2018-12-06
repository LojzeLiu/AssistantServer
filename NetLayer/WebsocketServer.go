package NetLayer

import (
	"Common"
	"golang.org/x/net/websocket"
)

var gClients = make([]*WeatherServer)

func WSAccept(ws *websocket.Conn) {
	Common.DEBUG("On Accept ")
	CurrClient := &WeatherServer{}
	go CurrClient.Listener(ws)

	gClients = append(gClients, CurrClient)
	Common.DEBUG("Deal end...")
}
