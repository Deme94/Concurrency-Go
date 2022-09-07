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
	numberOfPizzas int
	buffer         [3]*pizza
	mutex          *sync.Mutex
}

const urgentCallTime = 70 * time.Second
const cookingMaxTime = 5
const cookingMinTime = 3
const goingToPizzeriaTime = 2 * time.Second
const deliveryMaxTime = 12
const deliveryMinTime = 10

var PIZZERIA = pizzeria{
	0,
	[3]*pizza{},
	new(sync.Mutex),
}
var wg sync.WaitGroup

func produce(numberOfOrders int) {
	for i := 0; i < numberOfOrders; {

		if i == 8 {
			color.Red("THE CHEF HAS RECEIVED AN URGENT CALL FROM HIS WIFE!")
			time.Sleep(urgentCallTime)
		}

		PIZZERIA.mutex.Lock()
		pizzasReady := PIZZERIA.numberOfPizzas
		PIZZERIA.mutex.Unlock()
		if pizzasReady == len(PIZZERIA.buffer) {
			color.Red("Chef is waiting for delivery man...")
			for pizzasReady == len(PIZZERIA.buffer) {
				// Chef is waiting...
				PIZZERIA.mutex.Lock()
				pizzasReady = PIZZERIA.numberOfPizzas
				PIZZERIA.mutex.Unlock()
			}
		} else {
			color.Red("Pizza #%d is being prepared...", i+1)

			cookingTime := rand.Intn(cookingMaxTime-cookingMinTime) + cookingMinTime
			time.Sleep(time.Duration(cookingTime) * time.Second)

			pizza := pizza{i + 1}
			PIZZERIA.mutex.Lock()
			// Mutual exclusion ----
			for i, p := range PIZZERIA.buffer {
				if p == nil {
					PIZZERIA.buffer[i] = &pizza
					PIZZERIA.numberOfPizzas++
					color.Red("Pizza #%d is ready for delivery!", pizza.id)
					break
				}
			}
			// ---------------------
			PIZZERIA.mutex.Unlock()

			i++
		}

	}
	color.Red("Producer has finished")
}

func consume(numberOfOrders int) {
	for i := 0; i < numberOfOrders; {

		PIZZERIA.mutex.Lock()
		pizzasReady := PIZZERIA.numberOfPizzas
		PIZZERIA.mutex.Unlock()
		if pizzasReady == 0 {
			color.Yellow("Delivery man is waiting for chef...")
			for pizzasReady == 0 {
				// Delivery man is waiting...
				PIZZERIA.mutex.Lock()
				pizzasReady = PIZZERIA.numberOfPizzas
				PIZZERIA.mutex.Unlock()
			}
		} else {
			targetId := i + 1
			var targetPizza *pizza

			PIZZERIA.mutex.Lock()
			// Mutual exclusion ----
			for i, p := range PIZZERIA.buffer {

				if p != nil && p.id == targetId {
					targetPizza = PIZZERIA.buffer[i]
					PIZZERIA.buffer[i] = nil
					PIZZERIA.numberOfPizzas--
					color.Yellow("Pizza #%d is out for delivery!", targetPizza.id)
					break
				}
			}
			// ---------------------
			PIZZERIA.mutex.Unlock()

			deliveryTime := rand.Intn(deliveryMaxTime-deliveryMinTime) + deliveryMinTime
			time.Sleep(time.Duration(deliveryTime) * time.Second)
			color.Yellow("Pizza #%d delivered!", targetPizza.id)

			time.Sleep(goingToPizzeriaTime)
			color.Yellow("Delivery man is at the Pizzeria!")
			i++
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
	go consume(numberOfOrders)

	wg.Wait()

	fmt.Println("Pizzeria is CLOSED. See you tomorrow!")
}
