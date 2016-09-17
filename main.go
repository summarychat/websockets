package websockets

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var hubs map[string]Hub

func main() {
	r := gin.Default()

	r.GET("/:chat/:name", func(c *gin.Context) {
		chat := c.Param("chat")
		name := c.Param("name")
		coordinate(chat, name, c.Writer, c.Request)
	})

	r.Run("localhost:80")
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(t, msg)
	}
}

func coordinate(chat string, name string, w http.ResponseWriter, r *http.Request) {
	hub, ok := hubs[chat]
	if !ok {
		hub = *newHub(chat)
		hubs[chat] = hub
	}
	hub.addUser(name, w, r)
}
