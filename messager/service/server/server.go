package server

import (
	"net"
	"os"
	"sync"
	"time"

	"github.com/kisekivul/messager/common"
	"github.com/kisekivul/messager/protocol"
	"github.com/kisekivul/utils"
)

var Processor func(common.Action, string, string, ...interface{})

func Run(code, host string, timeout, limit int, requeue bool, processor func(common.Action, string, string, ...interface{}), params ...interface{}) {
	Processor = processor
	//initialize protocol code
	protocol.Initialize(code)
	//start service
	go run(host, timeout, limit, requeue, params...)
}

func run(host string, timeout, limit int, requeue bool, params ...interface{}) {
	listener, err := net.Listen("tcp", host)
	if utils.ErrorCheck(err) {
		os.Exit(1)
	}
	defer listener.Close()

	//push message
	go push(limit, requeue)
	//accept connection && read message
	for {
		conn, err := listener.Accept()
		if utils.ErrorCheck(err) {
			continue
		}
		//accept connection
		utils.LogData("Accept", conn.RemoteAddr().String())
		connection := &common.Connection{
			State: &common.State{
				Status: true,
			},
			Conn: conn,
		}
		//read message
		go read(connection, timeout, params...)
	}
}

func read(conn *common.Connection, timeout int, params ...interface{}) {
	data := make([][]byte, 0)
	code := make([]byte, 0)
	tail := make([]byte, 0)
	chMsg := make(chan []byte)

	defer func() {
		conn.Conn.Close()
		close(chMsg)
	}()

	var once sync.Once
	//set timeout
	timer := time.NewTimer(time.Second * time.Duration(timeout))

	for conn.State.Status {
		buffer := make([]byte, 1024)
		//receive message
		n, err := conn.Conn.Read(buffer)
		if err != nil {
			drop(conn)
			return
		}
		//unpacket message
		data, code, tail = protocol.Depack(append(tail, buffer[:n]...))
		if string(code) != "" {
			conn.Code = string(code)
			once.Do(func() {
				//handle message
				go handler(conn, chMsg, timeout, timer, params...)
			})
			//transfer message
			transfer(chMsg, data)
		}
	}
}

func handler(conn *common.Connection, chMsg chan []byte, timeout int, timer *time.Timer, params ...interface{}) {
	for conn.State.Status {
		select {
		case <-timer.C:
			drop(conn)
			return
		case data := <-chMsg:
			if data == nil {
				return
			}
			timer.Reset(time.Second * time.Duration(timeout))
			//set status
			conn.State.Status = true
			//analyze data
			action, _ := common.Analyze(data)
			switch action.Action {
			case common.Message:
				fallthrough
			case common.Registe:
				conn.State.Date = action.Date
				conn.State.Action = action.Action
				common.Add(conn.Code, conn)
				fallthrough
			// case common.Offline:
			// 	fallthrough
			default:
				if Processor != nil {
					Processor(action.Action, conn.Code, action.Data, params...)
				}
			}
			utils.LogData(conn.State.Action.String(), conn.Conn.RemoteAddr().String(), conn.Code)
		}
	}
}

func transfer(chMsg chan []byte, data [][]byte) {
	// panic if channel is closed
	defer func() {
		recover()
	}()
	//transfer
	for _, msg := range data {
		chMsg <- msg
	}
}

func drop(conn *common.Connection) {
	//set status
	conn.State.Status = false
	//disconnect
	conn.Conn.Close()

	common.Del(conn.Code)
	if Processor != nil {
		Processor(common.Offline, conn.Code, "")
	}
	utils.LogData(common.Offline.String(), conn.Conn.RemoteAddr().String(), conn.Code)
}

func List() []string {
	return common.List()
}
