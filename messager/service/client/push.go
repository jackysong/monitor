package client

import (
	"time"

	"github.com/kisekivul/messager/common"
)

var ch_push = make(chan string)

func Push(msg string) {
	ch_push <- msg
}

func push(conn *common.Connection, limit int, requeue bool) {
	if limit <= 0 {
		execute(conn, requeue)
	} else {
		for i := 0; i < limit; i++ {
			go execute(conn, requeue)
		}
	}
}

func execute(conn *common.Connection, requeue bool) {
	for msg := range ch_push {
		if !conn.State.Status {
			return
		}
		//send
		if msg != "" {
			send(conn, msg, common.Message, requeue)
		} else {
			time.Sleep(time.Second)
			continue
		}
	}
}
