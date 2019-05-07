package sample

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/kisekivul/messager/common"
	"github.com/kisekivul/messager/service/server"
	"github.com/kisekivul/utils"
)

func Server() {
	wg := new(sync.WaitGroup)
	wg.Add(1)

	go server.Run("backend", "localhost:9999", 5, 0, false, testServer)
	go test_push()

	wg.Wait()
}

func test_push() {
	//demo
	for {
		unit := new(common.Unit)
		unit.Code = "demo"
		data, _ := json.Marshal(&common.Packet{
			Date:   common.Now(),
			Data:   "get",
			Action: common.Message,
		})
		unit.Data = string(data)

		server.Push(unit)
		time.Sleep(time.Second)

		data, _ = json.Marshal(&common.Packet{
			Date:   common.Now(),
			Data:   "set",
			Action: common.Message,
		})
		unit.Data = string(data)

		server.Push(unit)
		time.Sleep(time.Second)
	}
}

func testServer(action common.Action, code, msg string, params ...interface{}) {
	switch action {
	case common.Offline:
		fallthrough
	case common.Registe:
		fallthrough
	case common.Message:
		fallthrough
	default:
		utils.LogData(action, code, msg)
	}
}
