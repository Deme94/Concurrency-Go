package main

import (
	"fmt"
	"sync"
)

var msg string
var wg sync.WaitGroup

func writeMsg(s string, m *sync.Mutex) {
	defer wg.Done()

	m.Lock()
	msg = s
	m.Unlock()
}

func main() {

	var mutex sync.Mutex
	//var mutexRW sync.RWMutex // (optional mutex for write/read threads)

	wg.Add(2)
	go writeMsg("Hola", &mutex)
	go writeMsg("Adios", &mutex)
	wg.Wait()

	fmt.Println(msg)
}
