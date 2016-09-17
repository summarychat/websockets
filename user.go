package websockets

import "github.com/gorilla/websocket"

type User struct {
	conn         *websocket.Conn
	name         string
	hub          *Hub
	broadcast    chan []byte
	toSend       chan []byte
	disconnect   chan string
	unregistered bool
}

func makeUser(hub *Hub, name string, conn *websocket.Conn) *User {
	user := new(User)
	user.conn = conn
	user.hub = hub
	user.broadcast = hub.broadcast
	user.toSend = make(chan []byte, 35)
	user.disconnect = make(chan string, 2)
	user.unregistered = false
	user.run()
	return user
}

func (user *User) run() {
	go user.read()
	go user.write()
}

func (user *User) read() {
	for {
		_, msg, err := user.conn.ReadMessage()
		if err != nil {
			break
		}
		user.broadcast <- msg
	}
	user.logout()
}

func (user *User) write() {
	var msg []byte
	for {
		msg = <-user.toSend
		err := user.conn.WriteMessage(websocket.BinaryMessage, msg)
		if err != nil {
			break
		}
	}
	user.logout()
}

func (user *User) logout() {
	if !user.unregistered {
		user.hub.unregister(user.name)
		user.disconnect <- "disconnect"
		user.unregistered = true
	}
}
