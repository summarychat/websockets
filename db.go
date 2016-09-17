package main

import (
	"database/sql"
	js "encoding/json"
	"log"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initStorage() {
	db, err := sql.Open("postgres", "postgres://root@localhost:26257?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	createTables(db)
}

func createTables(db *sql.DB) error {
	statements := []string{
		`CREATE DATABASE IF NOT EXISTS context`,
		`CREATE TABLE IF NOT EXISTS context.events (
        event_id INT,
        event_type STRING,
        msg_id INT,
		scores STRING
        PRIMARY KEY (event_id)
		FOREIGN KEY (msg_id) REFERENCES messages(msg_id)
    )`,
		`CREATE TABLE IF NOT EXISTS context.messages (
        channel STRING,
        msg_id INT,
        message STRING,
        timestamp TIMESTAMP
        PRIMARY KEY (msg_id)
    )`,
	}
	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}

type Message struct {
	User string
	Msg  string
}

func storeJSON(json []byte, channel string) {
	msg := parseJSON(json)
	const insertSQL = `
INSERT INTO context.messages VALUES ($1, DEFAULT, $2, NOW())`
	if _, err := db.Exec(insertSQL, channel, msg.Msg); err != nil {
		log.Printf("insert into messages failed: %s", err)
	}
}

func parseJSON(json []byte) *Message {
	msg := new(Message)
	js.Unmarshal(json, msg)
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
