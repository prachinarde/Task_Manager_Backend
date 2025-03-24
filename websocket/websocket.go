package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]bool)

func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	clients[conn] = true

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			delete(clients, conn)
			conn.Close()
			break
		}

		for client := range clients {
			client.WriteMessage(messageType, msg)
		}
	}
}
