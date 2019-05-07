package sample

import (
	"sync"

	"github.com/kisekivul/messager/service/client"
	"github.com/kisekivul/utils"
)

func Client() {
	wg := new(sync.WaitGroup)
	wg.Add(1)

	client.Run("demo", "localhost:9999", "", 1, 0, false, testClient)

	wg.Wait()
}

func testClient(msg string, params ...interface{}) {
	switch msg {
	case "set":
		utils.LogData("set case")
		client.Push("finish set")
	case "get":
		utils.LogData("get case")
		client.Push("finish get")
	default:
		utils.LogData(msg)
	}
}
