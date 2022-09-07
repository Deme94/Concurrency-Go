package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/fatih/color"
)

type pizza struct {
	id int
}
type pizzeria struct {
	channel chan pizza
}

const urgentCallTime = 70 * time.Second
const cookingMaxTime = 5
const cookingMinTime = 3
const goingToPizzeriaTime = 2 * time.Second
const deliveryMaxTime = 12
const deliveryMinTime = 10

var PIZZERIA = pizzeria{
	make(chan pizza, 3),
}
var wg sync.WaitGroup

func produce(numberOfOrders int) {
	for i := 0; i < numberOfOrders; i++ {

		if i == 8 {
			color.Red("THE CHEF HAS RECEIVED AN URGENT CALL FROM HIS WIFE!")
			time.Sleep(urgentCallTime)
		}

		color.Red("Pizza #%d is being prepared...", i+1)
		cookingTime := rand.Intn(cookingMaxTime-cookingMinTime) + cookingMinTime
		time.Sleep(time.Duration(cookingTime) * time.Second)
		p := pizza{i + 1}

		if len(PIZZERIA.channel) == cap(PIZZERIA.channel) {
			color.Red("Chef is waiting for delivery man...")
		}
		PIZZERIA.channel <- p // locked until channel has free space
		color.Red("Pizza #%d ready for delivery!", p.id)
	}
	color.Red("Producer has finished")
	close(PIZZERIA.channel)
}

func consume() {
	// locked in each iteration until channel length is greater than 0
	for pizza := range PIZZERIA.channel {
		color.Yellow("Pizza #%d is out for delivery!", pizza.id)
		deliveryTime := rand.Intn(deliveryMaxTime-deliveryMinTime) + deliveryMinTime
		time.Sleep(time.Duration(deliveryTime) * time.Second)
		color.Yellow("Pizza #%d delivered!", pizza.id)
		time.Sleep(goingToPizzeriaTime)
		color.Yellow("Delivery man is at the Pizzeria!")
		if len(PIZZERIA.channel) == 0 {
			color.Yellow("Delivery man is waiting for chef...")
		}
	}

	color.Yellow("Consumer has finished")
	wg.Done()
}

func main() {

	fmt.Println("Pizzeria is OPEN. Make your order!")

	wg.Add(1)

	numberOfOrders := 10
	go produce(numberOfOrders)
	go consume()

	wg.Wait()

	fmt.Println("Pizzeria is CLOSED. See you tomorrow!")
}
