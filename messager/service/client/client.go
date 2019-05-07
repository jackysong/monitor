package client

import (
	"encoding/json"
	"net"
	"os"
	"sync"
	"time"

	"github.com/kisekivul/messager/common"
	"github.com/kisekivul/messager/protocol"
	"github.com/kisekivul/utils"
)

var Processor func(string, ...interface{})

func Run(code, server, inform string, interval, limit int, requeue bool, processor func(string, ...interface{}), params ...interface{}) {
	Processor = processor
	//initialize protocol code
	protocol.Initialize(code)
	//start service
	go run(server, inform, interval, limit, requeue, params...)
}

func run(server, inform string, interval, limit int, requeue bool, params ...interface{}) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		utils.LogData(err.Error())
		os.Exit(1)
	}

	for {
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			utils.LogData(err.Error())
			time.Sleep(time.Duration((interval)) * time.Second)
			continue
		}
		defer conn.Close()

		connection := &common.Connection{
			State: &common.State{
				Status: true,
			},
			Conn: conn,
		}
		go read(connection, params...)
		go push(connection, limit, requeue)
		beating(connection, interval, inform)
	}
}

func beating(conn *common.Connection, interval int, inform string) {
	//predeclare
	send(conn, inform, common.Message, true)

	for {
		if send(conn, "", common.Registe, false) != nil {
			return
		}
		//interval
		time.Sleep(time.Duration((interval)) * time.Second)
	}
}

func read(conn *common.Connection, params ...interface{}) {
	data := make([][]byte, 0)
	code := make([]byte, 0)
	tail := make([]byte, 0)
	buffer := make([]byte, 1024)
	ch_msg := make(chan []byte)

	var once sync.Once
	for conn.State.Status {
		//receive message
		n, err := conn.Conn.Read(buffer)
		if err != nil {
			return
		}
		//unpacket message
		data, code, tail = protocol.Depack(append(tail, buffer[:n]...))
		if string(code) != "" {
			conn.Code = string(code)
			once.Do(func() {
				//handle message
				go handler(ch_msg, params...)
			})
			//transfer message
			transfer(ch_msg, data)
		} else {
			transfer(ch_msg, [][]byte{[]byte("disconn")})
		}
	}
}

func transfer(ch_msg chan []byte, data [][]byte) {
	// panic if channel is closed
	defer func() {
		recover()
	}()
	//transfer
	for _, msg := range data {
		ch_msg <- msg
	}
}

func send(conn *common.Connection, msg string, action common.Action, requeue bool) error {
	message := &common.Packet{
		Date:   common.Now(),
		Data:   msg,
		Action: action,
	}
	result, _ := json.Marshal(message)

	err := common.Send(*conn, result)
	if utils.ErrorCheck(err) {
		//requeue
		if requeue {
			ch_push <- msg
		}
		//close conn
		drop(*conn)
		//start new client
		return err
	}
	return nil
}

func handler(ch_msg chan []byte, params ...interface{}) {
	if Processor != nil {
		for msg := range ch_msg {
			switch string(msg) {
			case "disconn":
				return
			default:
				Processor(string(msg), params...)
			}
		}
	}
}

func drop(conn common.Connection) {
	//set status
	conn.State.Status = false
	//disconnect
	conn.Conn.Close()

	common.Del(conn.Code)

	utils.LogData("Close", conn.Conn.RemoteAddr().String())
}
