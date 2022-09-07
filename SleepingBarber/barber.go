package main

import (
	"fmt"
	"sync"
	"time"
)

type Barber struct {
	awake        bool
	sleep        *sync.WaitGroup
	barbershop   *Barbershop
	maxCustomers int
	nextCustomer *Customer
}

func (b *Barber) rest() {
	fmt.Println("Barber: *zzZZzZZZzZZzzzzZz*")
	b.sleep.Add(1)
	b.awake = false
	b.barbershop.mutex.Unlock()
	b.sleep.Wait()
	b.awake = true
	fmt.Println("Barber: *'o_O*")
}
func (b *Barber) wakeUp(c *Customer) {
	b.nextCustomer = c
	b.sleep.Done()
}

func (b *Barber) cutHair() {
	fmt.Printf("Barber: *cuts %s hair...*\n", b.nextCustomer.name)
	time.Sleep(2 * time.Second)
	b.nextCustomer.quit.Done()
	b.nextCustomer = nil
}

func (b *Barber) checkWaitingRoom() *Customer {
	fmt.Println("Barber: *checks the waiting room...*")

	//time.Sleep(1 * time.Second)

	var c *Customer

	select {
	case c = <-b.barbershop.waitingRoom:
		return c
	default:
		return nil
	}
}

// main thread
func (b *Barber) work() {
	for i := 0; i < b.maxCustomers; i++ {

		b.barbershop.mutex.Lock()
		b.nextCustomer = b.checkWaitingRoom()

		if b.nextCustomer == nil {
			b.rest()
		} else {
			b.barbershop.mutex.Unlock()
		}

		b.cutHair()
	}
	fmt.Println("THE BARBER HAS FINISHED.")
	b.barbershop.close.Done()
}
