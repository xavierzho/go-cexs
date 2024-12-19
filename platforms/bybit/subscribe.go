package bybit

import (
	"github.com/gorilla/websocket"
	"github.com/xavierzho/go-cexs/utils"
)

type Message struct {
	RequestId string   `json:"req_id"`
	Operation string   `json:"op"`
	Args      []string `json:"args"`
}

func (c *Connector) subscribe(message Message) {
	bytesMsg, err := utils.Json.Marshal(&message)
	if err != nil {
	}
	c.Channels["spot"].WriteMessage(websocket.TextMessage, bytesMsg)
}
