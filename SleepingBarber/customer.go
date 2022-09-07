package main

import (
	"fmt"
	"sync"
	"time"
)

type Customer struct {
	name       string
	barbershop *Barbershop
	quit       *sync.WaitGroup
}

func (c *Customer) leaveBarbershop() {
	fmt.Printf("%s is dissapointed...\n", c.name)
	c.quit.Done()
}

func (c *Customer) wakeUpBarber() {
	c.barbershop.barber.wakeUp(c)
}

func (c *Customer) wait() bool {
	select {
	case c.barbershop.waitingRoom <- c:
		fmt.Printf("%s is waiting...\n", c.name)
		return true
	default:
		return false
	}
}

func (c *Customer) checkBarberAsleep() bool {
	return !c.barbershop.barber.awake
}

func (c *Customer) goToBarbershopAfter(seconds int) {
	c.quit.Add(1)

	time.Sleep(time.Duration(seconds) * time.Second)
	fmt.Printf("%s has entered the barbershop.\n", c.name)

	c.barbershop.mutex.Lock()
	if c.checkBarberAsleep() {
		c.wakeUpBarber()
	} else {
		if !c.wait() {
			c.leaveBarbershop()
		}
	}
	c.barbershop.mutex.Unlock()

	c.quit.Wait()

	fmt.Printf("%s leaves the barbershop.\n", c.name)
}
