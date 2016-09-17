package websockets

import "net/http"

type Hub struct {
	users     map[string]User
	broadcast chan []byte
	name      string
}

func newHub(name string) *Hub {
	hub := new(Hub)
	hub.broadcast = make(chan []byte, 50)
	hub.users = make(map[string]User)
	hub.name = name
	return hub
}

func (hub *Hub) run() {

}

func (hub *Hub) addUser(name string, w http.ResponseWriter, r *http.Request) {

}

func (hub *Hub) unregister(name string)
