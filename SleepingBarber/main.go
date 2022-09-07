package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/goombaio/namegenerator"
)

func main() {

	maxCustomers := 10
	var barbershop Barbershop

	barber := Barber{true, &sync.WaitGroup{}, &barbershop, maxCustomers, nil}
	barbershop = Barbershop{
		&barber,
		&sync.Mutex{},
		make(chan *Customer, 3),
		&sync.WaitGroup{},
	}

	barbershop.close.Add(1)
	fmt.Println("Barbershop is open!")

	// Start barber thread
	go barbershop.barber.work()

	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	for i := 0; i < 35; i++ {
		name := nameGenerator.Generate()

		c := Customer{name, &barbershop, &sync.WaitGroup{}}

		randomSeconds := rand.Intn(100)
		go c.goToBarbershopAfter(randomSeconds)
	}

	barbershop.close.Wait()
}
