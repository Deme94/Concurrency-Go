package main

import (
	"fmt"
	"time"
)

type writer struct {
}

func (w *writer) write(data *string) {

	start.Wait()

	mutex.Lock()

	// Write data
	fmt.Println("writer is writing...")
	*data = "blobloblobloblo"
	time.Sleep(2 * time.Second)

	mutex.Unlock()

	wg.Done()
}
