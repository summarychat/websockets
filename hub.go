package main

import (
	"fmt"
	"net/http"
)

type Hub struct {
	users     map[string]*User
	broadcast chan []byte
	name      string
}

func newHub(name string) *Hub {
	hub := new(Hub)
	hub.broadcast = make(chan []byte, 50)
	hub.users = make(map[string]*User)
	hub.name = name
	go hub.run()
	return hub
}

func (hub *Hub) run() {
	for {
		msg := <-hub.broadcast
		storeJSON(msg, hub.name)
		for _, user := range hub.users {
			user.toSend <- msg
		}
	}
}

func (hub *Hub) addUser(name string, w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}
	hub.users[name] = makeUser(hub, name, conn)
}

func (hub *Hub) unregister(name string) {
	delete(hub.users, name)
}
