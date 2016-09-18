package main

import (
    "fmt"
    "net/http"
)

type Carrier struct {
    content []byte
    user string
}

type Hub struct {
    users     map[string]*User
    broadcast chan *Carrier
    name      string
}

func newHub(name string) *Hub {
    hub := new(Hub)
    hub.broadcast = make(chan *Carrier, 50)
    hub.users = make(map[string]*User)
    hub.name = name
    go hub.run()
    return hub
}

func (hub *Hub) run() {
    for {
        car := <-hub.broadcast
        msg := car.content
        storeJSON(msg, hub.name)
        for _, user := range hub.users {
            fmt.Print(string(msg))
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
