package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var hubs map[string]*Hub

func main() {
	r := gin.Default()

	r.GET("/:chat/:name", func(c *gin.Context) {
		chat := c.Param("chat")
		name := c.Param("name")
		coordinate(chat, name, c.Writer, c.Request)
	})

	r.Run("0.0.0.0:80")
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

func coordinate(chat string, name string, w http.ResponseWriter, r *http.Request) {
	hub, ok := hubs[chat]
	if !ok {
		hub = newHub(chat)
		hubs[chat] = hub
	}
	hub.addUser(name, w, r)
}
