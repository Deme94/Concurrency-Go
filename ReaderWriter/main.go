package main

import (
	"fmt"
	"sync"
)

var data string
var mutex sync.RWMutex

var wg sync.WaitGroup // wait until every thread has finished

var start sync.WaitGroup // Optional

func main() {
	data = "Data every thread needs to access"

	start.Add(1)

	for i := 0; i < 100; i++ {

		if i < 7 {
			w := writer{}
			go w.write(&data)
			wg.Add(1)
		}

		r := reader{}
		go r.read(&data)
		wg.Add(1)
	}

	start.Done()

	wg.Wait()
	fmt.Println("Everyone is done")
}
