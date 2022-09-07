package main

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func Test_updateMessage(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	go updateMessage("Hola Mundo", &wg)

	wg.Wait()

	if msg != "Hola Mundo" {
		t.Error("Expected 'Hola Mundo', but message was not updated")
	}
}

func Test_printMessage(t *testing.T) {
	msg = "Caracas"

	stdout := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	printMessage()

	w.Close()
	os.Stdout = stdout

	result, _ := io.ReadAll(r)
	output := string(result)

	if !strings.Contains(output, "Caracas") {
		t.Error("Expected 'Caracas', but it is not there")
	}
}

func Test_main(t *testing.T) {
	stdout := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	w.Close()
	os.Stdout = stdout

	result, _ := io.ReadAll(r)
	output := string(result)

	if !strings.Contains(output, "Hello, universe!\nHello, cosmos!\nHello, world!") {
		t.Errorf("Expected:\n'Hello, universe!\nHello, cosmos!\nHello, world!', \nbut got:\n'%s'", output)
	}
}
