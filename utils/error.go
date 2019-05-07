package utils

import (
	"log"
)

func ErrorCheck(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}

func ErrorStop(err error) bool {
	if err != nil {
		log.Fatal(err)
		return true
	}
	return false
}
