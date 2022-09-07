package main

import (
	"fmt"
	"strings"
)

type reader struct {
}

func (r *reader) read(data *string) {

	start.Wait()

	mutex.RLock()

	// Read data
	fmt.Println("reader is reading...")
	_ = strings.Contains(*data, "blo")

	mutex.RUnlock()

	wg.Done()
}
