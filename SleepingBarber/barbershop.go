package main

import "sync"

type Barbershop struct {
	barber      *Barber
	mutex       *sync.Mutex    // Syncronize barber and customer
	waitingRoom chan *Customer // Shared buffer
	close       *sync.WaitGroup
}
