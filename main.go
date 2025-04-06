package main

import (
	"sync"

	"github.com/silver-volt4/swapdoodle/nex"
)

var wg sync.WaitGroup

func main() {
	wg.Add(2)
	// TODO - Add gRPC server
	go nex.StartHppServer()
	wg.Wait()
}
