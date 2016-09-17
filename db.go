package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "strconv"
    "time"

    _ "github.com/lib/pq"
)

var db *sql.DB

func initStorage() {
	var err error
    db, err = sql.Open("postgres", "postgres://root@localhost:26257?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    createTables(db)
}

func createTables(db *sql.DB) error {
    statements := []string{
        `CREATE DATABASE IF NOT EXISTS context`,
        `CREATE TABLE IF NOT EXISTS context.messages (
        channel STRING,
		namegi STRING,
        msg_id INT,
        message STRING,
        timestamp TIMESTAMP,
        PRIMARY KEY (msg_id)
    )`,
	    `CREATE TABLE IF NOT EXISTS context.events (
        event_id INT,
        event_type STRING,
        msg_id INT references context.messages(msg_id),
		index(msg_id),
        scores STRING,
        channel STRING,
        PRIMARY KEY (event_id)
    )`,
    }
    for _, stmt := range statements {
        if _, err := db.Exec(stmt); err != nil {
            panic(err)
        }
    }
    return nil
}

type Message struct {
    User string `json:"user"`
    Msg  string `json:"msg"`
}

func storeJSON(data []byte, channel string) {
    msg := parseJSON(data)
    const insertSQL = `
INSERT INTO context.messages VALUES ($1, $2, DEFAULT, $3, NOW());`
    if _, err := db.Exec(insertSQL, channel, msg.User, msg.Msg); err != nil {
        log.Printf("insert into messages failed: %s", err)
    }
    fmt.Printf(msg.User, msg.Msg)
}

func parseJSON(data []byte) *Message {
    msg := new(Message)
    json.Unmarshal(data, msg)
    fmt.Printf(string(data))
    fmt.Printf(msg.User, msg.Msg)
    return msg
}

func convertFromUnix(b []byte) time.Time {
    ts, err := strconv.Atoi(string(b))
    if err != nil {
        panic(err)
    }

    t := time.Time(time.Unix(int64(ts), 0))

    return t
}
