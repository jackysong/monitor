package main

import (
	"sync"

	"github.com/kisekivul/messager/sample"
)

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(1)

	go sample.Server()

	wg.Wait()
}
