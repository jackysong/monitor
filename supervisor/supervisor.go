package supervisor

import (
	"errors"
	"log"
)

var (
	errorParseHead  = errors.New("parse head fail")
	errorParseBody  = errors.New("parse body fail")
	errorBodyLength = errors.New("body length inconformity")
)

func Initialize(handler func(*Head, *Body, ...interface{}), params ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	reload()

	listen(handler, params)
}

func ListPrograms() []Program {
	return c.list()
}
