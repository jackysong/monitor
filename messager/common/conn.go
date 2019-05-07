package common

import (
	"net"
	"sync"
)

type Connection struct {
	Code  string
	State *State
	Conn  net.Conn
}

type State struct {
	Packet
	Status bool `json:"status"`
}

type Packet struct {
	Date   int    `json:"date"`
	Data   string `json:"data"`
	Action Action `json:"action"`
}

var dicConn sync.Map

func Add(key, val interface{}) {
	dicConn.Store(key, val)
}

func Get(key interface{}) (interface{}, bool) {
	return dicConn.Load(key)
}

func Del(key interface{}) {
	dicConn.Delete(key)
}

func List() []string {
	list := make([]string, 0)
	f := func(key, value interface{}) (result bool) {
		list = append(list, key.(string))
		return true
	}
	dicConn.Range(f)
	return list
}
