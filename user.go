package websockets

import (
	"sync"

	"github.com/gorilla/websocket"
)

type User struct {
	conn       *websocket.Conn
	name       string
	hub        *Hub
	broadcast  chan []byte
	toSend     chan []byte
	disconnect chan string
	unregister sync.Mutex
}

func makeUser(hub *Hub, name string, conn *websocket.Conn) {
	user := new(User)
	user.conn = conn
	user.hub = hub
	user.broadcast = hub.broadcast
	user.toSend = make(chan []byte, 35)
	user.disconnect = make(chan string, 2)
}

func (user *User) run() {
	go user.read()
	go user.write()
}

func (user *User) read() {
	for {

	}
	user.logout()
	user.disconnect <- "disconnect"
}

func (user *User) write() {
	for {

	}

}

func (user *User) logout() {
	user.hub.unregister(user.name)
	user.disconnect <- "disconnect"
}
