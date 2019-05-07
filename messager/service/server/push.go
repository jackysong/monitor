package server

import (
	"time"

	"github.com/kisekivul/messager/common"
	"github.com/kisekivul/messager/protocol"
	"github.com/kisekivul/utils"
)

var chPush = make(chan *common.Unit)

func Push(action *common.Unit) {
	chPush <- action
}

func push(limit int, requeue bool) {
	if limit <= 0 {
		execute(requeue)
	} else {
		for i := 0; i < limit; i++ {
			go execute(requeue)
		}
	}
}

func execute(requeue bool) {
	for unit := range chPush {
		if unit != nil {
			send(unit, requeue)
		} else {
			time.Sleep(time.Second)
			continue
		}
	}
}

func send(action *common.Unit, requeue bool) {
	val, exist := common.Get(action.Code)
	if !exist {
		if requeue {
			chPush <- action
		}
		return
	}
	conn := val.(*common.Connection)

	_, err := conn.Conn.Write(protocol.Enpack([]byte(action.Data)))
	if utils.ErrorCheck(err) {
		drop(conn)
	}
}
